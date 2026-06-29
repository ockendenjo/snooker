package drinks

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDrink_ValidateTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(now time.Time) (*Drink, *time.Time, *time.Time)
		checkErr func(*testing.T, error)
	}{
		{
			name: "no start or end",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				return &Drink{Timestamp: new(now)}, nil, nil
			},
			checkErr: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			name: "timestamp before start",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				past := now.Add(-time.Hour)
				return &Drink{Timestamp: new(past)}, new(now), nil
			},
			checkErr: func(t *testing.T, err error) { assert.ErrorAs(t, err, &ErrTimestampBeforeStart{}) },
		},
		{
			name: "timestamp equal to start",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				return &Drink{Timestamp: new(now)}, new(now), nil
			},
			checkErr: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			name: "timestamp after end",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				future := now.Add(time.Hour)
				return &Drink{Timestamp: new(future)}, nil, new(now)
			},
			checkErr: func(t *testing.T, err error) { assert.ErrorAs(t, err, &ErrTimestampAfterEnd{}) },
		},
		{
			name: "timestamp equal to end",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				return &Drink{Timestamp: new(now)}, nil, new(now)
			},
			checkErr: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			name: "timestamp within range",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				past, future := now.Add(-time.Hour), now.Add(time.Hour)
				return &Drink{Timestamp: new(now)}, new(past), new(future)
			},
			checkErr: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			name: "timestamp before range",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				past, future := now.Add(-time.Hour), now.Add(time.Hour)
				return &Drink{Timestamp: new(past)}, new(now), new(future)
			},
			checkErr: func(t *testing.T, err error) { assert.ErrorAs(t, err, &ErrTimestampBeforeStart{}) },
		},
		{
			name: "timestamp after range",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				past, future := now.Add(-time.Hour), now.Add(time.Hour)
				return &Drink{Timestamp: new(future)}, new(past), new(now)
			},
			checkErr: func(t *testing.T, err error) { assert.ErrorAs(t, err, &ErrTimestampAfterEnd{}) },
		},
		{
			name: "timestamp in future",
			setup: func(now time.Time) (*Drink, *time.Time, *time.Time) {
				past, future := now.Add(-time.Hour), now.Add(time.Hour)
				return &Drink{Timestamp: new(now.Add(30 * time.Minute))}, new(past), new(future)
			},
			checkErr: func(t *testing.T, err error) {
				assert.ErrorAs(t, err, &ErrTimestampInFuture{})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				d, start, end := tc.setup(time.Now())
				tc.checkErr(t, d.ValidateTimestamp(start, end, time.Now()))
			})
		})
	}
}
