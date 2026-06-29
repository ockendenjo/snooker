package drinks

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ErrVersionMismatch struct{}

func (e *ErrVersionMismatch) Error() string {
	return "drink version mismatch"
}

type Client interface {
	GetDrink(ctx context.Context, userID string, timestamp time.Time) (*Drink, error)
	PutDrink(ctx context.Context, drink *Drink) error
	ListDrinksForUser(ctx context.Context, userID string) Query
	ListAll(ctx context.Context) ([]*Drink, error)
	ListUnknownDrinks(ctx context.Context) ([]*Drink, error)
}

const pk = "user_id"
const sk = "tstamp"

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

func (c *client) ListAll(ctx context.Context) ([]*Drink, error) {
	pages := dynamodb.NewScanPaginator(c.ddbClient, &dynamodb.ScanInput{
		TableName: &c.tableName,
	})

	var drinks []*Drink

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		var pageDrinks []*Drink
		err = attributevalue.UnmarshalListOfMaps(page.Items, &pageDrinks)
		if err != nil {
			return nil, err
		}
		drinks = append(drinks, pageDrinks...)
	}
	return drinks, nil
}

func (c *client) GetDrink(ctx context.Context, userID string, timestamp time.Time) (*Drink, error) {
	result, err := c.ddbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      &c.tableName,
		ConsistentRead: new(true),
		Key: map[string]dynamoTypes.AttributeValue{
			pk: &dynamoTypes.AttributeValueMemberS{Value: userID},
			sk: &dynamoTypes.AttributeValueMemberS{Value: timestamp.UTC().Truncate(time.Second).Format(time.RFC3339)},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var drink Drink
	if err := attributevalue.UnmarshalMap(result.Item, &drink); err != nil {
		return nil, err
	}
	return &drink, nil
}

func (c *client) PutDrink(ctx context.Context, drink *Drink) error {
	if drink.Timestamp != nil {
		drink.Timestamp = new(drink.Timestamp.UTC().Truncate(time.Second))
	}

	oldVersion := drink.Version
	drink.Version = oldVersion + 1

	m, err := attributevalue.MarshalMap(drink)
	if err != nil {
		return err
	}

	_, err = c.ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &c.tableName,
		Item:                m,
		ConditionExpression: new("attribute_not_exists(user_id) OR version = :version"),
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

func (c *client) ListDrinksForUser(ctx context.Context, userID string) Query {
	return &query{c: c, ctx: ctx, userID: userID}
}

func (c *client) ListUnknownDrinks(ctx context.Context) ([]*Drink, error) {

	pages := dynamodb.NewScanPaginator(c.ddbClient, &dynamodb.ScanInput{
		TableName: &c.tableName,
		IndexName: new("unknown-beer"),
	})

	allUnknown := make([]*Drink, 0, 10)

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		chunks := slices.Chunk(page.Items, 25)
		for chunk := range chunks {
			d, err := c.handleChunk(ctx, chunk)
			if err != nil {
				return nil, err
			}
			allUnknown = append(allUnknown, d...)
		}
	}

	return allUnknown, nil
}

func (c *client) handleChunk(ctx context.Context, chunk []map[string]dynamoTypes.AttributeValue) ([]*Drink, error) {

	for _, item := range chunk {
		delete(item, "unknown_beer")
	}

	res, err := c.ddbClient.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]dynamoTypes.KeysAndAttributes{
			c.tableName: {
				Keys: chunk,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var drinks []*Drink
	err = attributevalue.UnmarshalListOfMaps(res.Responses[c.tableName], &drinks)
	if err != nil {
		return nil, err
	}
	return drinks, nil
}
