package drinks

import (
	"context"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Query builds and executes a drinks query for a single user.
type Query interface {
	After(t time.Time) Query
	Before(t time.Time) Query
	MaybeBefore(t *time.Time) Query
	Run() ([]*Drink, error)
	Iterate() iter.Seq2[*Drink, error]
	Limit(limit int32) Query
	Reverse() Query
}

type query struct {
	c                  *client
	ctx                context.Context
	userID             string
	after              *time.Time
	before             *time.Time
	scanIndexBackwards bool
	limit              *int32
}

func (q *query) Limit(limit int32) Query {
	q.limit = &limit
	return q
}

func (q *query) MaybeBefore(t *time.Time) Query {
	if t == nil {
		return q
	}
	return q.Before(*t)
}

func (q *query) Reverse() Query {
	q.scanIndexBackwards = true
	return q
}

// After restricts results to drinks with a timestamp >= t.
func (q *query) After(t time.Time) Query {
	q.after = &t
	return q
}

// Before restricts results to drinks with a timestamp < t.
func (q *query) Before(t time.Time) Query {
	q.before = &t
	return q
}

// Run executes the query and collects all drinks into a slice.
func (q *query) Run() ([]*Drink, error) {
	var drinks []*Drink
	for d, err := range q.Iterate() {
		if err != nil {
			return nil, err
		}
		drinks = append(drinks, d)
	}
	return drinks, nil
}

// Iterate executes the query and yields each drink in turn. If an error occurs
// it is yielded as the second value and iteration stops.
func (q *query) Iterate() iter.Seq2[*Drink, error] {
	return func(yield func(*Drink, error) bool) {
		keyExpr, exprValues := q.buildKeyCondition()

		input := dynamodb.QueryInput{
			TableName:                 &q.c.tableName,
			ConsistentRead:            new(true),
			KeyConditionExpression:    &keyExpr,
			ScanIndexForward:          new(!q.scanIndexBackwards),
			ExpressionAttributeValues: exprValues,
		}
		if q.limit != nil {
			input.Limit = q.limit
		}

		paginator := dynamodb.NewQueryPaginator(q.c.ddbClient, &input)
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(q.ctx)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, item := range page.Items {
				var drink Drink
				if err := attributevalue.UnmarshalMap(item, &drink); err != nil {
					yield(nil, err)
					return
				}
				if !yield(&drink, nil) {
					return
				}
			}
		}
	}
}

func (q *query) buildKeyCondition() (string, map[string]dynamoTypes.AttributeValue) {
	parts := []string{fmt.Sprintf("%s = :userID", pk)}
	values := map[string]dynamoTypes.AttributeValue{
		":userID": &dynamoTypes.AttributeValueMemberS{Value: q.userID},
	}

	if q.after != nil && q.before != nil {
		parts = append(parts, fmt.Sprintf("%s BETWEEN :after AND :before", sk))
		values[":after"] = &dynamoTypes.AttributeValueMemberS{Value: q.after.Format(time.RFC3339)}
		values[":before"] = &dynamoTypes.AttributeValueMemberS{Value: q.before.Format(time.RFC3339)}
	} else if q.after != nil {
		parts = append(parts, fmt.Sprintf("%s >= :after", sk))
		values[":after"] = &dynamoTypes.AttributeValueMemberS{Value: q.after.Format(time.RFC3339)}
	} else if q.before != nil {
		parts = append(parts, fmt.Sprintf("%s < :before", sk))
		values[":before"] = &dynamoTypes.AttributeValueMemberS{Value: q.before.Format(time.RFC3339)}
	}

	return strings.Join(parts, " AND "), values
}
