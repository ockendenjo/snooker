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
			name: "not in word",
			drinks: []*Drink{
				{
					Points: 5, NotInWord: true,
				},
			},
			want: []*ScoreRecord{
				{
					Drink: &Drink{Points: 5, NotInWord: true}, ActualPoints: 2,
				},
			},
		},
		{
			name: "in progress word",
			drinks: []*Drink{
				{Points: 5},
				{Points: 5},
			},
			want: []*ScoreRecord{
				{
					Drink: &Drink{Points: 5}, ActualPoints: 2,
				},
				{
					Drink: &Drink{Points: 5}, ActualPoints: 2,
				},
			},
		},
		{
			name: "completed word",
			drinks: []*Drink{
				{Points: 2},
				{Points: 3},
				{Points: 4, EndOfWord: true},
			},
			want: []*ScoreRecord{
				{Drink: &Drink{Points: 2}, ActualPoints: 2},
				{Drink: &Drink{Points: 3}, ActualPoints: 3},
				{Drink: &Drink{Points: 4, EndOfWord: true}, ActualPoints: 4, WordLength: 3, WordPoints: 270, Multiplier: 10, SumLetters: 9},
			},
		},
		{
			name: "completed word over 2 days",
			drinks: []*Drink{
				{Points: 2, Timestamp: tm(3, 24, 12)},
				{Points: 2, Timestamp: tm(3, 24, 13)},
				{Points: 2, Timestamp: tm(3, 25, 14), EndOfWord: true},
			},
			want: []*ScoreRecord{
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 24, 12)}, ActualPoints: 2},
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 24, 13)}, ActualPoints: 2},
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 25, 14), EndOfWord: true}, ActualPoints: 2, WordLength: 3, WordPoints: 180, Multiplier: 10, SumLetters: 6},
			},
		},
		{
			name: "completed word over 3 days",
			drinks: []*Drink{
				{Points: 2, Timestamp: tm(3, 24, 12)},
				{Points: 2, Timestamp: tm(3, 25, 13)},
				{Points: 2, Timestamp: tm(3, 26, 14), EndOfWord: true},
			},
			want: []*ScoreRecord{
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 24, 12)}, ActualPoints: 2},
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 25, 13)}, ActualPoints: 2},
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 26, 14), EndOfWord: true}, ActualPoints: 2, WordLength: 3, WordPoints: 162, Multiplier: 9, SumLetters: 6},
			},
		},
		{
			name: "completed word over 3 days with gap",
			drinks: []*Drink{
				{Points: 2, Timestamp: tm(3, 24, 12)},
				{Points: 2, Timestamp: tm(3, 24, 13)},
				{Points: 2, Timestamp: tm(3, 26, 14), EndOfWord: true},
			},
			want: []*ScoreRecord{
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 24, 12)}, ActualPoints: 2},
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 24, 13)}, ActualPoints: 2},
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 26, 14), EndOfWord: true}, ActualPoints: 2, WordLength: 3, WordPoints: 162, Multiplier: 9, SumLetters: 6},
			},
		},
		{
			name: "incomplete word after completed word",
			drinks: []*Drink{
				{Points: 2, Timestamp: tm(3, 24, 12)},
				{Points: 3, Timestamp: tm(3, 24, 13)},
				{Points: 4, Timestamp: tm(3, 24, 14), EndOfWord: true},
				{Points: 3, Timestamp: tm(4, 1, 15)},
			},
			want: []*ScoreRecord{
				{Drink: &Drink{Points: 2, Timestamp: tm(3, 24, 12)}, ActualPoints: 2},
				{Drink: &Drink{Points: 3, Timestamp: tm(3, 24, 13)}, ActualPoints: 3},
				{Drink: &Drink{Points: 4, Timestamp: tm(3, 24, 14), EndOfWord: true}, ActualPoints: 4, WordLength: 3, WordPoints: 270, Multiplier: 10, SumLetters: 9},
				{Drink: &Drink{Points: 3, Timestamp: tm(4, 1, 15)}, ActualPoints: 2},
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
