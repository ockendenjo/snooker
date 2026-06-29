package testdb

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TestFlag = "RUN_DYNAMODB_LOCAL_TESTS"

func SetupTable(ctx context.Context, t *testing.T, dbClient *dynamodb.Client, createTableInput dynamodb.CreateTableInput) func() {
	log.Printf("creating table %s\n", *createTableInput.TableName)

	_, err := dbClient.CreateTable(ctx, &createTableInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("created table")
	count := 0

	for {
		if count >= 30 {
			t.Fatal("timeout waiting for table")
		}

		table, err := dbClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: createTableInput.TableName})
		if err != nil {
			t.Fatal(err)
		}
		if table.Table.TableStatus == dynamoTypes.TableStatusActive {
			log.Println("table is ready")
			break
		}
		time.Sleep(time.Second)
		log.Println("waiting for table to be ready")
		count++
	}

	return func() {
		log.Println("deleting table")
		_, err := dbClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: createTableInput.TableName})
		if err != nil {
			t.Fatalf("Failed to delete table %s - %v", *createTableInput.TableName, err)
			return
		}
		log.Println("deleted table")
	}
}
