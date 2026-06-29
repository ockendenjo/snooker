package apighandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ockendenjo/handler"
)

type HttpHandlerFunc[T any] func(ctx *handler.Context, event events.APIGatewayProxyRequest) (T, error)

type H = handler.Handler[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]

func GetHandler[T any](handleFn HttpHandlerFunc[T], statusCode int) handler.Handler[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse] {
	return func(ctx *handler.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		logger := ctx.GetLogger()

		res, err := handleFn(ctx, event)
		if err != nil {
			if httpError, ok := errors.AsType[HttpError](err); ok {
				response := events.APIGatewayProxyResponse{StatusCode: httpError.StatusCode, Body: httpError.Message, Headers: map[string]string{
					"Content-Type": "text/plain",
				}}
				logger.AddParam("log", httpError.Log).With("response", response).Error(err.Error())
				return response, nil
			}
			logger.Error(err.Error())
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}

		b, err := json.Marshal(res)
		if err != nil {
			logger.Error(err.Error())
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: statusCode, Body: string(b), Headers: map[string]string{
			"Content-Type": "application/json",
		}}, nil
	}
}

type HttpError struct {
	StatusCode int
	Err        error
	Message    string
	Log        string
}

func (h HttpError) Error() string {
	if h.Err != nil {
		return fmt.Sprintf("Returning HTTP error %d with message %s because %s", h.StatusCode, h.Message, h.Err.Error())
	}
	return fmt.Sprintf("Returning HTTP error %d with message %s", h.StatusCode, h.Message)
}

func (h HttpError) Unwrap() error {
	return h.Err
}

func GetEmailFromClaims(event events.APIGatewayProxyRequest) (string, error) {
	claims, ok := event.RequestContext.Authorizer["claims"]
	if !ok {
		return "", errors.New("no claims found in request context")
	}
	claimsMap, ok := claims.(map[string]any)
	if !ok {
		return "", errors.New("claims are not a map")
	}
	email, ok := claimsMap["email"].(string)
	if !ok {
		return "", errors.New("claims do not contain email")
	}
	return email, nil
}
