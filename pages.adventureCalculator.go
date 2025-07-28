package main

import (
	"net/http"

	gamesystems "github.com/FloodedRealms/borderland-keep-2.0/game-systems"
)

type PageAdventureCalculator struct {}

/* This handles the requests for the page itself, and the full calculation*/
func adventureCalculator(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := gamesystems.GetACKSIIInputForm()
if err != nil {
			panic(err)
			return
		}
		output, err := renderTemplateWithData("adventureCalculator", data)
		if err != nil {
			panic(err)
			return
		}
		w.Write([]byte(output))
	}
}


func (p PageAdventureCalculator) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/calculator", adventureCalculator)
}
