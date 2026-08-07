package drinks

import (
	"iter"
	"time"
)

func ScoreDrinks(drinkList []*Drink, loc *time.Location) []*ScoreRecord {
	records := make([]*ScoreRecord, 0, len(drinkList))
	numReds := 0

	for drinksForDay := range iterateDrinksByDay(drinkList) {
		expRed := numReds < 15

		for _, d := range drinksForDay {
			points := d.getScoreFromABV()
			if points == 1 {
				numReds++
			}

			if d.IsFoul {
				points = -4
				r := &ScoreRecord{Drink: d, TotalPoints: points}
				records = append(records, r)
				continue
			}

			if expRed && points != 1 {
				// Then it's a foul
				points = min(-4, -1*points)
			} else if !expRed && points == 1 {
				// Then it's a foul
				points = min(-4, -1*points)
			}

			r := &ScoreRecord{Drink: d, TotalPoints: points}
			records = append(records, r)

			if numReds < 15 {
				expRed = !expRed
			} else {
				expRed = false
			}
		}
	}

	return records
}

func iterateDrinksByDay(drinkList []*Drink) iter.Seq[[]*Drink] {
	return func(yield func([]*Drink) bool) {
		chunk := make([]*Drink, 0)
		var lastTS string

		for _, d := range drinkList {
			ts := d.Timestamp.Add(-3 * time.Hour).Format(time.DateOnly)

			if lastTS != ts && len(chunk) > 0 {
				if !yield(chunk) {
					return
				}
				chunk = make([]*Drink, 0)
			}

			lastTS = ts
			chunk = append(chunk, d)
		}

		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

// ScoreRecord extends Drink by adding some extra fields
type ScoreRecord struct {
	*Drink

	// TotalPoints indicates the total number of points this drink scored, including any penalty points
	TotalPoints int `json:"points"`
}

func (d *Drink) getScoreFromABV() int {
	switch {
	case d.ABV >= 54:
		return 7
	case d.ABV >= 51:
		return 6
	case d.ABV >= 48:
		return 5
	case d.ABV >= 45:
		return 4
	case d.ABV >= 42:
		return 3
	case d.ABV >= 39:
		return 2
	default:
		return 1
	}
}
