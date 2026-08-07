package drinks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ScoreDrinks(t *testing.T) {
	loc := time.UTC

	tests := []struct {
		name    string
		drinks  []*Drink
		checkFn func(t *testing.T, got []*ScoreRecord)
	}{
		{
			name:   "empty",
			drinks: []*Drink{},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				assert.Empty(t, got)
			},
		},
		{
			name: "one drink",
			drinks: []*Drink{
				{ABV: 38, Timestamp: tm(7, 7, 16)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "one drink with foul",
			drinks: []*Drink{
				{ABV: 38, Timestamp: tm(7, 7, 16), IsFoul: true},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 16), IsFoul: true}), TotalPoints: -4},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "second drink with green",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 43, Timestamp: tm(7, 7, 17)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 7, 17)}), TotalPoints: 3},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "second drink red instead of colour",
			drinks: []*Drink{
				{ABV: 38, Timestamp: tm(7, 7, 16)},
				{ABV: 38, Timestamp: tm(7, 7, 17)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 17)}), TotalPoints: -4},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "second drink is colour on different day",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 43, Timestamp: tm(7, 8, 16)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 8, 16)}), TotalPoints: -4},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "second drink is red on different day",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 5, Timestamp: tm(7, 8, 16)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 5, Timestamp: tm(7, 8, 16)}), TotalPoints: 1},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "third drink is black",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 43, Timestamp: tm(7, 7, 17)},
				{ABV: 58, Timestamp: tm(7, 7, 18)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 7, 17)}), TotalPoints: 3},
					{Drink: new(Drink{ABV: 58, Timestamp: tm(7, 7, 18)}), TotalPoints: -7},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "third drink is pink",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 43, Timestamp: tm(7, 7, 17)},
				{ABV: 52, Timestamp: tm(7, 7, 18)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 7, 17)}), TotalPoints: 3},
					{Drink: new(Drink{ABV: 52, Timestamp: tm(7, 7, 18)}), TotalPoints: -6},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "third drink is blue",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 43, Timestamp: tm(7, 7, 17)},
				{ABV: 48, Timestamp: tm(7, 7, 18)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 7, 17)}), TotalPoints: 3},
					{Drink: new(Drink{ABV: 48, Timestamp: tm(7, 7, 18)}), TotalPoints: -5},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "third drink is other colour",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 43, Timestamp: tm(7, 7, 17)},
				{ABV: 44, Timestamp: tm(7, 7, 18)},
			},
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
					{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 7, 17)}), TotalPoints: 3},
					{Drink: new(Drink{ABV: 44, Timestamp: tm(7, 7, 18)}), TotalPoints: -4},
				}
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "yellow in final break",
			drinks: appendFrames([]*Drink{
				{ABV: 40, Timestamp: tm(7, 16, 16)},
			}),
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 40, Timestamp: tm(7, 16, 16)}), TotalPoints: 2},
				}
				assert.Equal(t, exp, got[30:])
			},
		},
		{
			name: "yellow and green in final break",
			drinks: appendFrames([]*Drink{
				{ABV: 40, Timestamp: tm(7, 16, 16)},
				{ABV: 43, Timestamp: tm(7, 16, 18)},
			}),
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 40, Timestamp: tm(7, 16, 16)}), TotalPoints: 2},
					{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 16, 18)}), TotalPoints: 3},
				}
				assert.Equal(t, exp, got[30:])
			},
		},
		{
			name: "yellow and brown in final break",
			drinks: appendFrames([]*Drink{
				{ABV: 40, Timestamp: tm(7, 16, 16)},
				{ABV: 46, Timestamp: tm(7, 16, 18)},
			}),
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 40, Timestamp: tm(7, 16, 16)}), TotalPoints: 2},
					{Drink: new(Drink{ABV: 46, Timestamp: tm(7, 16, 18)}), TotalPoints: -4},
				}
				assert.Equal(t, exp, got[30:])
			},
		},
		{
			name: "yellow and black in final break",
			drinks: appendFrames([]*Drink{
				{ABV: 40, Timestamp: tm(7, 16, 16)},
				{ABV: 59, Timestamp: tm(7, 16, 18)},
			}),
			checkFn: func(t *testing.T, got []*ScoreRecord) {
				exp := []*ScoreRecord{
					{Drink: new(Drink{ABV: 40, Timestamp: tm(7, 16, 16)}), TotalPoints: 2},
					{Drink: new(Drink{ABV: 59, Timestamp: tm(7, 16, 18)}), TotalPoints: -7},
				}
				assert.Equal(t, exp, got[30:])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScoreDrinks(tt.drinks, loc)
			tt.checkFn(t, got)
		})
	}
}

func tm(m, d, h int) *time.Time {
	return new(time.Date(2026, time.Month(m), d, h, 0, 0, 0, time.UTC))
}

func appendFrames(drinks []*Drink) []*Drink {
	a := make([]*Drink, 0, 30+len(drinks))

	for i := range 15 {
		//Add red
		a = append(a, new(Drink{ABV: 37, Timestamp: tm(7, i+1, 16)}))
		//Add brown
		a = append(a, new(Drink{ABV: 47, Timestamp: tm(7, i+1, 18)}))
	}

	a = append(a, drinks...)
	return a
}
