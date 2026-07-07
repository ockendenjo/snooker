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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScoreDrinks(tt.drinks, loc)
			assert.Equal(t, tt.want, got)
		})
	}
}
