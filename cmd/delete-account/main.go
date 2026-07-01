package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ockendenjo/handler"
	"github.com/ockendenjo/snooker/pkg/apighandler"
	"github.com/ockendenjo/snooker/pkg/user"
	"github.com/ockendenjo/snooker/pkg/userpool"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {
		dynamoClient := dynamodb.NewFromConfig(awsConfig)
		cognitoClient := cognitoidentityprovider.NewFromConfig(awsConfig)
		userTable := handler.MustGetEnv("USER_TABLE_NAME")
		userPoolId := handler.MustGetEnv("USER_POOL_ID")

		h := &lambdaHandler{
			userClient: user.NewClient(dynamoClient, userTable),
			poolClient: userpool.NewClient(cognitoClient, userPoolId),
		}
		return apighandler.GetHandler(h.handle, http.StatusNoContent)
	})
}

type lambdaHandler struct {
	userClient user.Client
	poolClient userpool.Client
}

func (h *lambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (any, error) {
	email, err := apighandler.GetEmailFromClaims(event)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusUnauthorized, Err: err}
	}

	if err := h.userClient.DeleteByEmail(ctx, email); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, Err: err, Message: "failed to delete user"}
	}

	if err := h.poolClient.DeleteUser(ctx, email); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, Err: err, Message: "failed to delete user from pool"}
	}

	return nil, nil
}
