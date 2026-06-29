package drinks

import "time"

const (
	MinLength       = 3
	PointsNotInWord = 2
)

func ScoreDrinks(d []*Drink, loc *time.Location) []*ScoreRecord {
	records := make([]*ScoreRecord, 0, len(d))
	pointsForWord := 0
	length := 0
	var wordDrinks []*Drink
	var wordRecords []*ScoreRecord

	for _, drink := range d {
		r := &ScoreRecord{Drink: drink}

		if drink.NotInWord {
			r.ActualPoints = PointsNotInWord
			records = append(records, r)
			continue
		}

		r.ActualPoints = drink.Points
		pointsForWord += drink.Points
		length++
		wordDrinks = append(wordDrinks, drink)
		wordRecords = append(wordRecords, r)

		if drink.EndOfWord {
			if length < MinLength {
				r.TooShort = true
				r.WordPoints = length * PointsNotInWord * 10
				for _, wr := range wordRecords {
					wr.ActualPoints = PointsNotInWord
				}
			} else {
				multiplier := wordMultiplier(wordDrinks, loc)
				r.Multiplier = multiplier
				r.SumLetters = pointsForWord
				r.WordLength = length
				r.WordPoints = length * pointsForWord * multiplier
			}

			pointsForWord = 0
			length = 0
			wordDrinks = nil
			wordRecords = nil
		}

		records = append(records, r)
	}

	for _, wr := range wordRecords {
		wr.ActualPoints = PointsNotInWord
	}

	return records
}

func wordMultiplier(wordDrinks []*Drink, loc *time.Location) int {
	var minTime *time.Time
	var maxTime *time.Time

	for _, d := range wordDrinks {
		if d.Timestamp == nil {
			continue
		}
		t := d.Timestamp.In(loc).Add(-3 * time.Hour)

		if minTime == nil {
			minTime = &t
			maxTime = &t
			continue
		}

		if t.Before(*minTime) {
			minTime = &t
		}
		if t.After(*maxTime) {
			maxTime = &t
		}
	}

	count := 1
	if (minTime != nil) && (maxTime != nil) {
		for minTime.Day() != maxTime.Day() || minTime.Month() != maxTime.Month() || minTime.Year() != maxTime.Year() {
			count++
			minTime = new(minTime.AddDate(0, 0, 1))
		}
	}

	switch count {
	case 1, 2:
		return 10
	case 3, 4:
		return 9
	case 5, 6:
		return 8
	default:
		return 7
	}
}

type ScoreRecord struct {
	*Drink

	// ActualPoints indicates the number of points a letter scores. For a completed word it is equal to the number of points in a Drink
	ActualPoints int `json:"actualPoints"`

	// TooShort indicates whether a word is too short. It is only ever true for the last letter of a completed word
	TooShort bool `json:"tooShort,omitzero"`

	// Multiplier represents the duration multiplier for a completed word. It is only non-zero for the last letter of a completed word.
	Multiplier int `json:"multiplier,omitzero"`

	// SumLetters represents the sum of the letter points for a completed word. It is only non-zero for the last letter of a completed word
	SumLetters int `json:"sumLetters,omitzero"`

	// WordLength is the length of a completed word. It is only non-zero for the last letter of a completed word
	WordLength int `json:"wordLength,omitzero"`

	// WordPoints is the total number of points for a completed word. It is only non-zero for the last letter of a completed word.
	WordPoints int `json:"wordPoints,omitzero"`
}
