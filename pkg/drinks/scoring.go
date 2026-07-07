package drinks

import "time"

func ScoreDrinks(d []*Drink, loc *time.Location) []*ScoreRecord {
	records := make([]*ScoreRecord, 0, len(d))

	for _, drink := range d {
		points := drink.getScoreFromABV()
		if drink.IsFoul {
			points = -4
		}

		r := &ScoreRecord{Drink: drink, TotalPoints: points}

		records = append(records, r)
	}

	return records
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
