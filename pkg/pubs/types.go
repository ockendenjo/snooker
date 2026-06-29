package pubs

import "strings"

type groupsFile struct {
	Groups GroupMap `json:"groups"`
}

type GroupMap map[string]*PubInGroup

type PubInGroup struct {
	Name   string  `json:"name"`
	Group  int     `json:"group"`
	Letter *string `json:"letter,omitempty"`
}

func (pg *PubInGroup) GetLetter() string {
	if pg.Letter != nil {
		return *pg.Letter
	}
	return GetPubLetter(pg.Name)
}

func GetPubLetter(name string) string {
	name = strings.ToLower(name)
	name = strings.TrimPrefix(name, "the ")
	return strings.ToUpper(name[0:1])
}

type pubFile struct {
	Pubs []*Pub `json:"pubs"`
}

type Pub struct {
	CamraID    int     `json:"camraID"`
	GoodBeerID *int    `json:"goodBeerID,omitempty"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	RealAles   int     `json:"realAles"`
	Notes      *string `json:"notes,omitempty"`
	Chain      *string `json:"chain,omitempty"`
	TempClosed bool    `json:"tempClosed,omitzero"`
}
