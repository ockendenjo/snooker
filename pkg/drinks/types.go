package drinks

import (
	"errors"
	"fmt"
	"time"
)

type Drink struct {
	UserID      string     `dynamodbav:"user_id"                json:"userID"`
	Timestamp   *time.Time `dynamodbav:"tstamp"                 json:"timestamp"`
	CamraID     *int       `dynamodbav:"camra_id"               json:"camraID"`
	PubName     string     `dynamodbav:"pub_name"               json:"pubName"`
	DrinkName   string     `dynamodbav:"drink_name"             json:"drinkName"`
	Brewery     string     `dynamodbav:"brewery"                json:"brewery"`
	ABV         *float32   `dynamodbav:"abv"                    json:"abv"`
	UntappdID   *int       `dynamodbav:"untappd_id,omitempty"   json:"untappdID"`
	Points      int        `dynamodbav:"points"                 json:"points"`
	With        string     `dynamodbav:"with"                   json:"with"`
	Version     int        `dynamodbav:"version"                json:"version"`
	Notes       *string    `dynamodbav:"notes,omitempty"        json:"notes,omitempty"`
	IsFoul      bool       `dynamodbav:"is_foul"                json:"isFoul"`
	UnknownBeer *int       `dynamodbav:"unknown_beer,omitempty" json:"-"`
	UnknownPub  *int       `dynamodbav:"unknown_pub,omitempty"  json:"-"`
}

const (
	MaxLengthDrinkName = 50
	MaxLengthPubName   = 80
	MaxLengthBrewery   = 50
	MaxLengthWith      = 100
	MaxLengthNotes     = 200
	MinABV             = 0.0
	MaxABV             = 20.0
)

func (d *Drink) Validate() error {
	if d.UserID == "" {
		return errors.New("userID is required")
	}
	if d.Timestamp == nil {
		return errors.New("timestamp is required")
	}
	if d.PubName == "" {
		return errors.New("pubName is required")
	}
	if len(d.PubName) > MaxLengthPubName {
		return fmt.Errorf("pubName must not exceed %d characters", MaxLengthPubName)
	}

	if d.DrinkName == "" {
		return errors.New("drinkName is required")
	}
	if len(d.DrinkName) > MaxLengthDrinkName {
		return fmt.Errorf("drinkName must not exceed %d characters", MaxLengthDrinkName)
	}
	if d.Brewery == "" {
		return errors.New("brewery is required")
	}
	if len(d.Brewery) > MaxLengthBrewery {
		return fmt.Errorf("brewery must not exceed %d characters", MaxLengthBrewery)
	}
	if d.ABV == nil {
		return errors.New("abv is required")
	}
	if *d.ABV < MinABV || *d.ABV > MaxABV {
		return errors.New("abv outside allowable range")
	}

	if d.With == "" {
		return errors.New("with is required")
	}
	if len(d.With) > MaxLengthWith {
		return fmt.Errorf("with must not exceed %d characters", MaxLengthWith)
	}
	if d.Notes != nil && len(*d.Notes) > MaxLengthNotes {
		return fmt.Errorf("notes must not exceed %d characters", MaxLengthNotes)
	}
	return nil
}

func (d *Drink) FixFields() {
	if d.UntappdID == nil {
		d.UnknownBeer = new(1)
	} else {
		d.UnknownBeer = nil
	}

	if d.CamraID == nil {
		d.UnknownPub = new(1)
	} else {
		d.UnknownPub = nil
	}
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
