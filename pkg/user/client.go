package user

import (
	"context"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ErrVersionMismatch struct{}

func (e *ErrVersionMismatch) Error() string {
	return "user version mismatch"
}

type Client interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	InsertUser(ctx context.Context, user User) error
	DeleteByEmail(ctx context.Context, email string) error
	ScanUsers(ctx context.Context) ([]*User, error)
}

const partitionKey = "email"
const gsiIndexKey = "id"

func NewClient(ddbClient *dynamodb.Client, tableName string) Client {
	return &client{
		ddbClient: ddbClient,
		tableName: tableName,
	}
}

type client struct {
	ddbClient *dynamodb.Client
	tableName string
}

func (c *client) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	result, err := c.ddbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: new(c.tableName),
		Key: map[string]dynamoTypes.AttributeValue{
			"email": &dynamoTypes.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var user User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *client) UpdateUser(ctx context.Context, user *User) error {
	oldVersion := user.Version
	user.Version = oldVersion + 1

	m, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = c.ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &c.tableName,
		Item:                m,
		ConditionExpression: new("version = :version"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":version": &dynamoTypes.AttributeValueMemberN{Value: strconv.Itoa(oldVersion)},
		},
	})
	if err != nil {
		if _, ok := errors.AsType[*dynamoTypes.ConditionalCheckFailedException](err); ok {
			return &ErrVersionMismatch{}
		}
		return err
	}
	return nil
}

func (c *client) InsertUser(ctx context.Context, user User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = c.ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           new(c.tableName),
		Item:                av,
		ConditionExpression: new("attribute_not_exists(email)"),
	})
	return err
}

func (c *client) GetUserByID(ctx context.Context, userID string) (*User, error) {
	result, err := c.ddbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              &c.tableName,
		IndexName:              aws.String("id-index"),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":id": &dynamoTypes.AttributeValueMemberS{Value: userID},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	var user User
	if err := attributevalue.UnmarshalMap(result.Items[0], &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *client) DeleteByEmail(ctx context.Context, email string) error {
	_, err := c.ddbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &c.tableName,
		Key: map[string]dynamoTypes.AttributeValue{
			"email": &dynamoTypes.AttributeValueMemberS{Value: email},
		},
	})
	return err
}

func (c *client) ScanUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	var lastKey map[string]dynamoTypes.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:        &c.tableName,
			FilterExpression: aws.String("total_points > :zero"),
			ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
				":zero": &dynamoTypes.AttributeValueMemberN{Value: "0"},
			},
		}
		if lastKey != nil {
			input.ExclusiveStartKey = lastKey
		}

		result, err := c.ddbClient.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			var user User
			if err := attributevalue.UnmarshalMap(item, &user); err != nil {
				return nil, err
			}
			users = append(users, &user)
		}

		if result.LastEvaluatedKey == nil {
			break
		}
		lastKey = result.LastEvaluatedKey
	}

	return users, nil
}
