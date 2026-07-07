package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ockendenjo/handler"
	"github.com/ockendenjo/snooker/pkg/apighandler"
	"github.com/ockendenjo/snooker/pkg/drinks"
	"github.com/ockendenjo/snooker/pkg/user"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {
		drinksTable := handler.MustGetEnv("DRINKS_TABLE_NAME")
		userTable := handler.MustGetEnv("USER_TABLE_NAME")
		ddbClient := dynamodb.NewFromConfig(awsConfig)

		startDate := optTimeVar("START_DATE")
		endDate := optTimeVar("END_DATE")

		h := &lambdaHandler{
			drinksClient: drinks.NewClient(ddbClient, drinksTable),
			userClient:   user.NewClient(ddbClient, userTable),
			startDate:    startDate,
			endDate:      endDate,
		}
		return apighandler.GetHandler(h.handle, http.StatusCreated)
	})
}

func optTimeVar(key string) *time.Time {
	s := os.Getenv(key)
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

type lambdaHandler struct {
	drinksClient drinks.Client
	userClient   user.Client
	startDate    *time.Time
	endDate      *time.Time
}

func (h *lambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (*LogDrinkResponse, error) {
	email, err := apighandler.GetEmailFromClaims(event)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusUnauthorized, Err: err}
	}

	userObj, err := h.userClient.GetUserByEmail(ctx, email)
	if err != nil || userObj == nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	var drink *drinks.Drink
	if err := json.Unmarshal([]byte(event.Body), &drink); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Err: err, Message: "invalid JSON body"}
	}

	drink.UserID = userObj.ID
	if drink.Timestamp == nil {
		drink.Timestamp = new(time.Now())
	}

	if err := drink.Validate(); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Err: err, Message: err.Error()}
	}
	if err := drink.ValidateTimestamp(h.startDate, h.endDate, time.Now()); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Err: err, Message: err.Error()}
	}

	drink.FixFields()
	if err := h.drinksClient.PutDrink(ctx, drink); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	cid, err := getCompositeID(userObj.ID, drink.Timestamp)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return &LogDrinkResponse{CompositeID: cid}, nil
}

type LogDrinkResponse struct {
	CompositeID string `json:"cid"`
}

func getCompositeID(userID string, timeStamp *time.Time) (string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(userID)); err != nil {
		return "", err
	}
	if _, err := h.Write([]byte(timeStamp.Format(time.RFC3339))); err != nil {
		return "", err
	}

	s := hex.EncodeToString(h.Sum(nil))
	return s[0:32], nil
}
