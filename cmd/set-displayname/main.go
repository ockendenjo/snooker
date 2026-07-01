package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ockendenjo/handler"
	"github.com/ockendenjo/snooker/pkg/apighandler"
	"github.com/ockendenjo/snooker/pkg/user"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {
		userTable := handler.MustGetEnv("USER_TABLE_NAME")
		ddbClient := dynamodb.NewFromConfig(awsConfig)
		userClient := user.NewClient(ddbClient, userTable)

		h := &lambdaHandler{
			userClient: userClient,
		}
		return apighandler.GetHandler(h.handle, http.StatusNoContent)
	})
}

type lambdaHandler struct {
	userClient user.Client
}

func (h *lambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (any, error) {
	email, err := apighandler.GetEmailFromClaims(event)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusUnauthorized, Err: err}
	}

	var req *SetUsernameRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Err: err, Message: "invalid JSON body"}
	}

	if valErr := req.Validate(); valErr != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Err: valErr, Message: valErr.Error()}
	}

	userObj, err := h.userClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	userObj.DisplayName = req.DisplayName

	err = h.userClient.UpdateUser(ctx, userObj)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	ctx.GetLogger().AddParam("userId", userObj.ID).Info("Display name updated")
	return nil, nil
}

type SetUsernameRequest struct {
	DisplayName string `json:"displayName"`
}

func (r *SetUsernameRequest) Validate() error {
	r.DisplayName = strings.TrimSpace(r.DisplayName)
	if r.DisplayName == "" {
		return errors.New("displayName is required")
	}
	if len(r.DisplayName) > 50 {
		return errors.New("displayName must be less than 50 characters")
	}
	return nil
}
