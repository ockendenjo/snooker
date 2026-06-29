package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/ockendenjo/handler"
	"github.com/ockendenjo/snooker/pkg/apighandler"
	"github.com/ockendenjo/snooker/pkg/session"
	"github.com/ockendenjo/snooker/pkg/user"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) handler.Handler[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse] {
		userTable := handler.MustGetEnv("USER_TABLE_NAME")
		ddbClient := dynamodb.NewFromConfig(awsConfig)

		h := &lambdaHandler{
			userClient: user.NewClient(ddbClient, userTable),
		}

		return apighandler.GetHandler(h.handle, http.StatusOK)
	})
}

type lambdaHandler struct {
	userClient user.Client
}

func (h *lambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (*session.Data, error) {
	email, err := apighandler.GetEmailFromClaims(event)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusUnauthorized, Err: err}
	}

	userObj, err := h.userClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError}
	}

	if userObj == nil {
		newUser := user.User{Email: email, ID: uuid.NewString(), Version: 0}
		if err := h.userClient.InsertUser(ctx, newUser); err != nil {
			return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError}
		}
		return &session.Data{ID: newUser.ID, Email: email}, nil
	}

	return &session.Data{
		ID:          userObj.ID,
		Email:       email,
		DisplayName: userObj.DisplayName,
		Points:      userObj.TotalPoints,
	}, nil
}
