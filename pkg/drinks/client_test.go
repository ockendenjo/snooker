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
		CamraID:   new(1),
		DrinkName: "Test Drink",
		ABV:       38,
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
		CamraID:   new(2),
		DrinkName: "Another Drink",
		ABV:       45,
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
	assert.Equal(t, "Test Drink", drinks[0].DrinkName)
	assert.Equal(t, "Another Drink", drinks[1].DrinkName)
	assert.Equal(t, 38, drinks[0].ABV)
	assert.Equal(t, 45, drinks[1].ABV)
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
