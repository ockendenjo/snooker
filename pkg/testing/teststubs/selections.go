package teststubs

import (
	"context"
	"time"

	"github.com/ockendenjo/snooker/pkg/selections"
)

// MockSelectionsClient is a mock implementation of selections.Client for testing
type MockSelectionsClient struct {
	GetAllSelectionsFn    func(ctx context.Context) ([]selections.Selection, error)
	PutSelectionsFn       func(ctx context.Context, selections []selections.Selection) error
	GetActivePubsFn       func(ctx context.Context, tm time.Time) (*selections.Selection, error)
	PutNextSelectionFn    func(ctx context.Context, sel *selections.Selection) error
	GetNextSelectionFn    func(ctx context.Context) (*selections.Selection, error)
	DeleteNextSelectionFn func(ctx context.Context) error
}

func (m *MockSelectionsClient) DeleteNextSelection(ctx context.Context) error {
	return m.DeleteNextSelectionFn(ctx)
}

func (m *MockSelectionsClient) GetNextSelection(ctx context.Context) (*selections.Selection, error) {
	return m.GetNextSelectionFn(ctx)
}

func (m *MockSelectionsClient) GetAllSelections(ctx context.Context) ([]selections.Selection, error) {
	return m.GetAllSelectionsFn(ctx)
}

func (m *MockSelectionsClient) PutSelections(ctx context.Context, sels []selections.Selection) error {
	return m.PutSelectionsFn(ctx, sels)
}

func (m *MockSelectionsClient) GetActivePubs(ctx context.Context, tm time.Time) (*selections.Selection, error) {
	return m.GetActivePubsFn(ctx, tm)
}

func (m *MockSelectionsClient) PutNextSelection(ctx context.Context, sel *selections.Selection) error {
	return m.PutNextSelectionFn(ctx, sel)
}
