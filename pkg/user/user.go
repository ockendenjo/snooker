package user

import (
	"time"
)

type User struct {
	Email        string     `dynamodbav:"email"`
	ID           string     `dynamodbav:"id"`
	Version      int        `dynamodbav:"version"`
	Code         string     `dynamodbav:"code"`
	Expiry       *time.Time `dynamodbav:"expiry"`
	CodeAttempts int        `dynamodbav:"attempts"`
	DisplayName  string     `dynamodbav:"display_name"`
	OverrideCode *string    `dynamodbav:"override_code,omitempty"`
	TotalPoints  int        `dynamodbav:"total_points,omitempty"`
}
