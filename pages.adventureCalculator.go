package main

import (
	"net/http"

	gamesystems "github.com/FloodedRealms/borderland-keep-2.0/game-systems"
)

type Adventure struct {
	NumberOfPlayers int
	NumberOfHenchmen int
	GP int
	SP int
	CP int
	EP int
	PP int
	SpecialTreasures []SpecialTreasure
	MagicItems []MagicItem
	Combat []CombatEncounter
}


type SpecialTreasureType string
const (
	Jewelery SpecialTreasureType = "Jewelery"
	Gemstone SpecialTreasureType = "Gemstone"
)

type SpecialTreasure struct {
	GPValue int
	SPType SpecialTreasureType
	Name string
}

type MagicItem struct {
	ApparentValue int
	IsSold bool
	SellValue int
	Name string
}

type CombatEncounter struct {
	MonsterName string
	NumberDefeated int
	XPValue int
}

type PageAdventureCalculator struct {}

/* This handles the requests for the page itself, and the full calculation*/
func adventureCalculator(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := gamesystems.GetACKSIIInputForm()
		if err != nil {
			panic(err)
		}
		output, err := renderTemplateWithData("adventureCalculator", data)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(output))
	}
}


func (p PageAdventureCalculator) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/calculator", adventureCalculator)
}
