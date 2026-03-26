package main

import (
	"log"
	"net/http"
)


type Weather struct {
	TemperatureDay string
	TemperatureNight string
	Precipitation string
	Wind string
	WindDirection string
}

type GenerationMode string
const (
	manual GenerationMode = "Manual"
	automatic GenerationMode = "Automatic"
	daily GenerationMode = "Daily"
	monthly GenerationMode = "Monthly"
)

type PageWeather struct {
	renderer *Renderer
	templateName string
	formTemplateName string
	Weather []Weather
	currentSystem string
	CurrentMode string
	ModeOptions []SelectOption
	KoppenCodeOptions []SelectOption
	WindOptions []SelectOption
	SimulateFront bool
	CurrentKoppenCode string
}

type SelectOption struct {
	Value, Text string
	Selected bool
}

func NewPageWeather() *PageWeather {
	allowedModes := []SelectOption{
		SelectOption{
			Value: string(manual),
			Text: string(manual),
			Selected: false,
		},
		SelectOption{
			Value: string(daily),
			Text: string(daily),
			Selected: true,
		},
		SelectOption{
			Value: string(monthly),
			Text: string(monthly),
			Selected: false,
		},
	}
	return &PageWeather{
		templateName: "weatherGenerator.html",
		formTemplateName: "weatherForm.html",
		Weather: []Weather{},
		renderer: NewRenderer(),
		currentSystem: "ACKS II",
		CurrentMode: string(automatic),
		ModeOptions: allowedModes,
	}
}

func (p PageWeather) RegisterRoutes(router *http.ServeMux) {
	router.Handle("/tools/calculator", p.index((*p.renderer)))
}

func (p PageWeather) index(pr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		renderedPage, err := pr.RenderPage(p.templateName, p.Weather)
		if err != nil {
			log.Printf("Error rendering Index of Adventure Calculator: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}
