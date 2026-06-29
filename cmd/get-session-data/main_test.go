package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ockendenjo/handler"
	"github.com/ockendenjo/snooker/pkg/testing/teststubs"
	"github.com/ockendenjo/snooker/pkg/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LettersInProgress(t *testing.T) {

	uc := &teststubs.MockUserClient{
		GetUserByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
			return &user.User{
				ID:    "1",
				Email: "user@example.com",
			}, nil
		},
	}

	h := &lambdaHandler{
		userClient: uc,
	}
	event := events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]any{
				"claims": map[string]any{
					"email": "user@example.com",
				},
			},
		},
	}

	sd, err := h.handle(handler.Get(t.Context()), event)
	require.NoError(t, err)

	assert.Equal(t, "user@example.com", sd.Email)
}
