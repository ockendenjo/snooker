package drinks

import "time"

func ScoreDrinks(d []*Drink, loc *time.Location) []*ScoreRecord {
	records := make([]*ScoreRecord, 0, len(d))

	for _, drink := range d {
		r := &ScoreRecord{Drink: drink}

		records = append(records, r)
	}

	return records
}

// ScoreRecord extends Drink by adding some extra fields
type ScoreRecord struct {
	*Drink

	// PenaltyPoints indicates the number of penalty points a drink incurred - e.g. by drinking out of order or because of a foul
	PenaltyPoints int `json:"penaltyPoints"`
}
