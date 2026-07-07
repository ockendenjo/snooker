package drinks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ScoreDrinks(t *testing.T) {
	loc := time.UTC

	tests := []struct {
		name   string
		drinks []*Drink
		want   []*ScoreRecord
	}{
		{
			name:   "empty",
			drinks: []*Drink{},
			want:   []*ScoreRecord{},
		},
		{
			name: "one drink",
			drinks: []*Drink{
				{ABV: 38, Timestamp: tm(7, 7, 16)},
			},
			want: []*ScoreRecord{
				{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
			},
		},
		{
			name: "one drink with foul",
			drinks: []*Drink{
				{ABV: 38, Timestamp: tm(7, 7, 16), IsFoul: true},
			},
			want: []*ScoreRecord{
				{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 16), IsFoul: true}), TotalPoints: -4},
			},
		},
		{
			name: "second drink with green",
			drinks: []*Drink{
				{ABV: 37, Timestamp: tm(7, 7, 16)},
				{ABV: 43, Timestamp: tm(7, 7, 17)},
			},
			want: []*ScoreRecord{
				{Drink: new(Drink{ABV: 37, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
				{Drink: new(Drink{ABV: 43, Timestamp: tm(7, 7, 17)}), TotalPoints: 3},
			},
		},
		{
			name: "second drink red instead of colour",
			drinks: []*Drink{
				{ABV: 38, Timestamp: tm(7, 7, 16)},
				{ABV: 38, Timestamp: tm(7, 7, 17)},
			},
			want: []*ScoreRecord{
				{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 16)}), TotalPoints: 1},
				{Drink: new(Drink{ABV: 38, Timestamp: tm(7, 7, 17)}), TotalPoints: -4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScoreDrinks(tt.drinks, loc)
			assert.Equal(t, tt.want, got)
		})
	}
}

func tm(m, d, h int) *time.Time {
	return new(time.Date(2026, time.Month(m), d, h, 0, 0, 0, time.UTC))
}
