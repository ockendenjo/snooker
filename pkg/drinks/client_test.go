package drinks

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/ockendenjo/snooker/pkg/testing/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_client(t *testing.T) {
	c, deleteFn := setupTableAndGetClient(t)
	defer deleteFn()

	const userID = "test-user-id"
	timestamp := time.Now().UTC()

	newDrink := Drink{
		UserID:    userID,
		Timestamp: &timestamp,
		PubID:     1,
		Name:      "Test Drink",
		Points:    10,
		EndOfWord: true,
		NotInWord: false,
		With:      "Friend",
		Version:   0,
	}
	err := c.PutDrink(t.Context(), &newDrink)
	require.NoError(t, err)

	// Insert another drink for the same user
	timestamp2 := time.Now().Add(time.Hour)
	newDrink2 := Drink{
		UserID:    userID,
		Timestamp: &timestamp2,
		PubID:     2,
		Name:      "Another Drink",
		Points:    15,
		EndOfWord: false,
		NotInWord: true,
		With:      "Colleague",
		Version:   0,
	}
	err = c.PutDrink(t.Context(), &newDrink2)
	require.NoError(t, err)

	// List drinks for user
	query := c.ListDrinksForUser(t.Context(), userID)
	drinks, err := query.Run()
	require.NoError(t, err)
	require.Len(t, drinks, 2)

	// Verify drinks are in ascending order by timestamp (oldest first)
	assert.Equal(t, timestamp.Truncate(time.Second).UTC(), drinks[0].Timestamp.UTC())
	assert.Equal(t, timestamp2.Truncate(time.Second).UTC(), drinks[1].Timestamp.UTC())
	assert.Equal(t, "Test Drink", drinks[0].Name)
	assert.Equal(t, "Another Drink", drinks[1].Name)
	assert.Equal(t, 10, drinks[0].Points)
	assert.Equal(t, 15, drinks[1].Points)
}

func Test_client_ListInProgressDrinks(t *testing.T) {
	c, deleteFn := setupTableAndGetClient(t)
	defer deleteFn()

	const userID = "list-in-progress-user"
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	putDrink := func(offset time.Duration, endOfWord bool) {
		t.Helper()
		err := c.PutDrink(t.Context(), &Drink{
			UserID:    userID,
			Timestamp: new(base.Add(offset)),
			PubID:     1,
			Name:      "Test",
			Brewery:   "Brewery",
			With:      "Friend",
			EndOfWord: endOfWord,
		})
		require.NoError(t, err)
	}

	t.Run("no drinks returns empty slice", func(t *testing.T) {
		result, err := c.ListInProgressDrinks(t.Context(), "no-drinks-user")
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("no EndOfWord returns all drinks in chronological order", func(t *testing.T) {
		const u = "user-no-eow"
		ts1, ts2 := base, base.Add(time.Hour)
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(ts1), PubID: 1, Name: "A", Brewery: "B", With: "F"}))
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(ts2), PubID: 1, Name: "B", Brewery: "B", With: "F"}))

		result, err := c.ListInProgressDrinks(t.Context(), u)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, ts1, result[0].Timestamp.UTC())
		assert.Equal(t, ts2, result[1].Timestamp.UTC())
	})

	t.Run("returns drinks after the most recent EndOfWord, in chronological order", func(t *testing.T) {
		putDrink(0, true)            // EndOfWord marker
		putDrink(time.Hour, false)   // in-progress drink 1
		putDrink(2*time.Hour, false) // in-progress drink 2

		result, err := c.ListInProgressDrinks(t.Context(), userID)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, base.Add(time.Hour), result[0].Timestamp.UTC())
		assert.Equal(t, base.Add(2*time.Hour), result[1].Timestamp.UTC())
	})

	t.Run("most recent drink is EndOfWord returns empty slice", func(t *testing.T) {
		const u = "user-eow-last"
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(base), PubID: 1, Name: "A", Brewery: "B", With: "F"}))
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(base.Add(time.Hour)), PubID: 1, Name: "B", Brewery: "B", With: "F", EndOfWord: true}))

		result, err := c.ListInProgressDrinks(t.Context(), u)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("stops at the most recent EndOfWord when multiple exist", func(t *testing.T) {
		const u = "user-multi-eow"
		ts4 := base.Add(3 * time.Hour)
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(base), PubID: 1, Name: "A", Brewery: "B", With: "F", EndOfWord: true}))
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(base.Add(time.Hour)), PubID: 1, Name: "B", Brewery: "B", With: "F"}))
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(base.Add(2 * time.Hour)), PubID: 1, Name: "C", Brewery: "B", With: "F", EndOfWord: true}))
		require.NoError(t, c.PutDrink(t.Context(), &Drink{UserID: u, Timestamp: new(ts4), PubID: 1, Name: "D", Brewery: "B", With: "F"}))

		result, err := c.ListInProgressDrinks(t.Context(), u)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, ts4, result[0].Timestamp.UTC())
	})
}

func setupTableAndGetClient(t *testing.T) (Client, func()) {
	//This test requires a running dynamodb-local container
	//podman run -p 8000:8000 docker.io/amazon/dynamodb-local
	if os.Getenv(testdb.TestFlag) != "1" {
		t.SkipNow()
		return nil, nil
	}

	credentialsProvider := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "secret", "session"))
	awsConfig, err := config.LoadDefaultConfig(t.Context(), credentialsProvider)
	if err != nil {
		t.Fatalf("unable to load SDK config, %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(awsConfig, func(o *dynamodb.Options) {
		o.BaseEndpoint = new("http://localhost:8000")
	})

	tableName := fmt.Sprintf("test-drinks-%s", uuid.NewString())

	drinksClient := NewClient(dynamoClient, tableName)

	deleteFn := testdb.SetupTable(t.Context(), t, dynamoClient, dynamodb.CreateTableInput{
		TableName: new(tableName),
		KeySchema: []dynamoTypes.KeySchemaElement{
			{
				AttributeName: new(pk),
				KeyType:       dynamoTypes.KeyTypeHash,
			},
			{
				AttributeName: new(sk),
				KeyType:       dynamoTypes.KeyTypeRange,
			},
		},
		AttributeDefinitions: []dynamoTypes.AttributeDefinition{
			{
				AttributeName: new(pk),
				AttributeType: dynamoTypes.ScalarAttributeTypeS,
			},
			{
				AttributeName: new(sk),
				AttributeType: dynamoTypes.ScalarAttributeTypeS,
			},
		},
		BillingMode: dynamoTypes.BillingModePayPerRequest,
	})

	return drinksClient, deleteFn
}
