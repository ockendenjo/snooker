package teststubs

import (
	"context"
	"iter"
	"time"

	"github.com/ockendenjo/snooker/pkg/drinks"
)

// MockDrinksClient is a mock implementation of drinks.Client for testing
type MockDrinksClient struct {
	GetDrinkFn          func(ctx context.Context, userID string, timestamp time.Time) (*drinks.Drink, error)
	PutDrinkFn          func(ctx context.Context, drink *drinks.Drink) error
	ListDrinksForUserFn func(ctx context.Context, userID string) drinks.Query
	ListAllFn           func(ctx context.Context) ([]*drinks.Drink, error)
	ListUnknownDrinksFn func(ctx context.Context) ([]*drinks.Drink, error)
}

func (m *MockDrinksClient) ListUnknownDrinks(ctx context.Context) ([]*drinks.Drink, error) {
	return m.ListUnknownDrinksFn(ctx)
}

func (m *MockDrinksClient) ListAll(ctx context.Context) ([]*drinks.Drink, error) {
	return m.ListAllFn(ctx)
}

func (m *MockDrinksClient) GetDrink(ctx context.Context, userID string, timestamp time.Time) (*drinks.Drink, error) {
	return m.GetDrinkFn(ctx, userID, timestamp)
}

func (m *MockDrinksClient) PutDrink(ctx context.Context, drink *drinks.Drink) error {
	return m.PutDrinkFn(ctx, drink)
}

func (m *MockDrinksClient) ListDrinksForUser(ctx context.Context, userID string) drinks.Query {
	return m.ListDrinksForUserFn(ctx, userID)
}

// MockQuery is a mock implementation of drinks.Query for testing
type MockQuery struct {
	RunFn func() ([]*drinks.Drink, error)
}

func (m *MockQuery) Iterate() iter.Seq2[*drinks.Drink, error] {
	return nil
}

func (m *MockQuery) Limit(_ int32) drinks.Query {
	return m
}

func (m *MockQuery) Reverse() drinks.Query {
	return m
}

func (m *MockQuery) MaybeBefore(_ *time.Time) drinks.Query {
	return m
}

func (m *MockQuery) After(_ time.Time) drinks.Query {
	return m
}

func (m *MockQuery) Before(_ time.Time) drinks.Query {
	return m
}

func (m *MockQuery) Run() ([]*drinks.Drink, error) {
	return m.RunFn()
}
