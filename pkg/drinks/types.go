package drinks

import (
	"errors"
	"time"
)

type Drink struct {
	UserID      string     `dynamodbav:"user_id"                json:"userID"`
	Timestamp   *time.Time `dynamodbav:"tstamp"                 json:"timestamp"`
	PubID       int        `dynamodbav:"pub_id"                 json:"pubID"`
	Name        string     `dynamodbav:"name"                   json:"name"`
	Brewery     string     `dynamodbav:"brewery"                json:"brewery"`
	UntappdID   *int       `dynamodbav:"untappd_id,omitempty"   json:"untappdID"`
	Points      int        `dynamodbav:"points"                 json:"points"`
	SelectedPub bool       `dynamodbav:"selected_pub"           json:"selectedPub"`
	EndOfWord   bool       `dynamodbav:"end_of_word"            json:"endOfWord"`
	NotInWord   bool       `dynamodbav:"not_in_word"            json:"notInWord"`
	With        string     `dynamodbav:"with"                   json:"with"`
	Version     int        `dynamodbav:"version"                json:"version"`
	Letter      string     `dynamodbav:"letter"                 json:"letter"`
	Notes       *string    `dynamodbav:"notes,omitempty"        json:"notes,omitempty"`
	UnknownBeer *int       `dynamodbav:"unknown_beer,omitempty" json:"-"`
}

func (d *Drink) Validate() error {
	if d.UserID == "" {
		return errors.New("userID is required")
	}
	if d.Timestamp == nil {
		return errors.New("timestamp is required")
	}
	if d.PubID == 0 {
		return errors.New("pubID is required")
	}
	if d.Name == "" {
		return errors.New("name is required")
	}
	if len(d.Name) > 50 {
		return errors.New("name must not exceed 50 characters")
	}
	if d.Brewery == "" {
		return errors.New("brewery is required")
	}
	if len(d.Brewery) > 50 {
		return errors.New("brewery must not exceed 50 characters")
	}
	if d.With == "" {
		return errors.New("with is required")
	}
	if len(d.With) > 100 {
		return errors.New("with must not exceed 100 characters")
	}
	if d.Notes != nil && len(*d.Notes) > 200 {
		return errors.New("notes must not exceed 200 characters")
	}
	return nil
}

type ErrTimestampBeforeStart struct{}

func (e ErrTimestampBeforeStart) Error() string {
	return "drink timestamp is before start of allowed range"
}

type ErrTimestampAfterEnd struct{}

func (e ErrTimestampAfterEnd) Error() string {
	return "drink timestamp is after end of allowed range"
}

type ErrTimestampInFuture struct{}

func (e ErrTimestampInFuture) Error() string {
	return "drink timestamp is in future"
}

func (d *Drink) ValidateTimestamp(start *time.Time, end *time.Time, now time.Time) error {
	if start != nil {
		if d.Timestamp.Before(*start) {
			return ErrTimestampBeforeStart{}
		}
	}
	if end != nil {
		if d.Timestamp.After(*end) {
			return ErrTimestampAfterEnd{}
		}
	}
	if d.Timestamp.After(now) {
		return ErrTimestampInFuture{}
	}
	return nil
}
