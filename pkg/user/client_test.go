package user

import (
	"fmt"
	"os"
	"testing"

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

	const email = "test@example.com"

	newUser := User{
		Email: email,
		ID:    uuid.NewString(),
	}
	err := c.InsertUser(t.Context(), newUser)
	require.NoError(t, err)

	lookupUser, err := c.GetUserByEmail(t.Context(), email)
	require.NoError(t, err)
	require.NotNil(t, lookupUser)
	assert.Equal(t, newUser.ID, lookupUser.ID)

	lookupUser.DisplayName = "JohnSmith"
	err = c.UpdateUser(t.Context(), lookupUser)
	require.NoError(t, err)

	lookupUser, err = c.GetUserByID(t.Context(), newUser.ID)
	require.NoError(t, err)
	require.NotNil(t, lookupUser)
	assert.Equal(t, email, lookupUser.Email)
	assert.Equal(t, "JohnSmith", lookupUser.DisplayName)
	assert.Equal(t, 1, lookupUser.Version)

	lookupUser.Version = 0
	err = c.UpdateUser(t.Context(), lookupUser)
	var expErr *ErrVersionMismatch
	require.ErrorAs(t, err, &expErr)

	err = c.DeleteByEmail(t.Context(), email)
	require.NoError(t, err)

	delUser, err := c.GetUserByEmail(t.Context(), email)
	require.NoError(t, err)
	require.Nil(t, delUser)
}

func Test_ScanUsers(t *testing.T) {
	c, deleteFn := setupTableAndGetClient(t)
	defer deleteFn()

	// Insert user with points
	userWithPoints := User{Email: "points@example.com", ID: uuid.NewString(), TotalPoints: 50}
	err := c.InsertUser(t.Context(), userWithPoints)
	require.NoError(t, err)

	// Insert user without points — should be excluded from scan results
	userNoPoints := User{Email: "nopoints@example.com", ID: uuid.NewString(), TotalPoints: 0}
	err = c.InsertUser(t.Context(), userNoPoints)
	require.NoError(t, err)

	users, err := c.ScanUsers(t.Context())
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, userWithPoints.ID, users[0].ID)
	assert.Equal(t, 50, users[0].TotalPoints)
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

	tableName := fmt.Sprintf("test-mpan-%s", uuid.NewString())

	lbClient := NewClient(dynamoClient, tableName)

	deleteFn := testdb.SetupTable(t.Context(), t, dynamoClient, dynamodb.CreateTableInput{
		TableName: new(tableName),
		KeySchema: []dynamoTypes.KeySchemaElement{
			{
				AttributeName: new(partitionKey),
				KeyType:       dynamoTypes.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []dynamoTypes.GlobalSecondaryIndex{
			{
				IndexName: new("id-index"),
				KeySchema: []dynamoTypes.KeySchemaElement{
					{
						AttributeName: new(gsiIndexKey),
						KeyType:       dynamoTypes.KeyTypeHash,
					},
				},
				Projection: &dynamoTypes.Projection{
					ProjectionType: dynamoTypes.ProjectionTypeAll,
				},
			},
		},
		AttributeDefinitions: []dynamoTypes.AttributeDefinition{
			{
				AttributeName: new(partitionKey),
				AttributeType: dynamoTypes.ScalarAttributeTypeS,
			},
			{
				AttributeName: new(gsiIndexKey),
				AttributeType: dynamoTypes.ScalarAttributeTypeS,
			},
		},
		BillingMode: dynamoTypes.BillingModePayPerRequest,
	})

	return lbClient, deleteFn
}
