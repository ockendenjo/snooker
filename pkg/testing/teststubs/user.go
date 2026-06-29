package teststubs

import (
	"context"

	"github.com/ockendenjo/snooker/pkg/user"
)

// MockUserClient is a mock implementation of user.Client for testing
type MockUserClient struct {
	GetUserByEmailFn func(ctx context.Context, email string) (*user.User, error)
	ScanUsersFn      func(ctx context.Context) ([]*user.User, error)
}

func (m *MockUserClient) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.GetUserByEmailFn != nil {
		return m.GetUserByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *MockUserClient) GetUserByID(_ context.Context, _ string) (*user.User, error) {
	return nil, nil
}

func (m *MockUserClient) UpdateUser(_ context.Context, _ *user.User) error {
	return nil
}

func (m *MockUserClient) InsertUser(_ context.Context, _ user.User) error {
	return nil
}

func (m *MockUserClient) DeleteByEmail(_ context.Context, _ string) error {
	return nil
}

func (m *MockUserClient) ScanUsers(ctx context.Context) ([]*user.User, error) {
	return m.ScanUsersFn(ctx)
}
