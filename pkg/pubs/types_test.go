package pubs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PubGroup_GetLetter(t *testing.T) {

	testcases := []struct {
		name string
		pg   *PubInGroup
		exp  string
	}{
		{
			name: "should use defined letter",
			pg:   &PubInGroup{Letter: new("U")},
			exp:  "U",
		},
		{
			name: "should use first letter",
			pg:   &PubInGroup{Name: "Cask & Barrel"},
			exp:  "C",
		},
		{
			name: "should trim the prefix when calculating letter",
			pg:   &PubInGroup{Name: "The Abbey"},
			exp:  "A",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.pg.GetLetter()
			assert.Equal(t, tc.exp, got)
		})
	}
}
