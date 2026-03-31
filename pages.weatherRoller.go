package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	gamesystems "github.com/FloodedRealms/borderland-keep-2.0/game-systems"
)

const weatherKey contextKey = "weather"

type WeatherForm struct {
	// Form Data
	currentSystem string
	CurrentMode string
	ModeOptions []*SelectOption
	KoppenCodeOptions []*SelectOption
	WindOptions []*SelectOption
	CurrentWind string
	SimulateFront bool
	CurrentKoppenCode string
	SeasonOptions []*SelectOption
	CurrentSeason gamesystems.Season
	WeatherResults []*Weather
}

type Weather struct {
	// Output Data
	UnmodifiedRolls []gamesystems.Roll
	Modifiers []gamesystems.Modifier
	TemperatureDay string
	TemperatureNight string
	Precipitation string
	Wind string
	WindDirection string
}

func NewWeatherForm() *WeatherForm {
	acks, _ := gamesystems.LoadGameSystem("ACKS II")
	koppenCodes := acks.ListKoppenCodes()
	winds := acks.ListWinds()
	currentWind := "Westerly"
	currentKoppen := "Cfa"
	allowedModes := []*SelectOption{
		&SelectOption{
			Value: string(manual),
			Text: string(manual),
			Selected: false,
		},
		&SelectOption{
			Value: string(automatic),
			Text: string(automatic),
			Selected: true,
		},
	}
	seasons := gamesystems.ListSeasons()
	cseason := gamesystems.Spring
	seasonOptions := buildSelectListFromStrings(seasons, cseason.String())
	windOptions := buildSelectListFromStrings(winds, currentWind)
	koppenOptions := buildSelectListFromStrings(koppenCodes, currentKoppen)
	simulateFront := false
	wr := make([]*Weather, 31)
	for i, _ := range wr {
		pr := []gamesystems.Roll{}
		if i > 0 {
			pr = wr[i-1].UnmodifiedRolls
		}
		rolls := acks.ListWeatherRolls()
		for _, r := range rolls {
			r.Resolve()
		}
		results, _ := acks.DailyWeather(rolls, pr, currentKoppen, currentWind, cseason, simulateFront)
		modifiers := acks.ListWeatherModifiers(currentKoppen, cseason)
		weather := &Weather{
			UnmodifiedRolls: rolls,
			Modifiers: modifiers,
			TemperatureDay: results[0],
			TemperatureNight: results[1],
			Precipitation: results[2],
			Wind: results[3],
			WindDirection: results[4],
		}
		wr[i] = weather
	}
	return &WeatherForm{
		currentSystem: "ACKS II",
		CurrentMode: string(automatic),
		ModeOptions: allowedModes,
		KoppenCodeOptions: koppenOptions,
		WindOptions: windOptions,
		CurrentWind: currentWind,
		SimulateFront: simulateFront,
		CurrentKoppenCode: currentKoppen,
		SeasonOptions: seasonOptions,
		CurrentSeason: cseason,
		WeatherResults: wr,
	}
}


/* This recalculates the weather results, keeping the old rolls intact. */
func (w WeatherForm) ApplyWeatherSettings() {
	system, _ := gamesystems.LoadGameSystem(w.currentSystem)
	for i, wr := range w.WeatherResults {
		// Setting the results to any empty array will tell the system that the first day has no previous rolls.
		// This means the first day will not be adjusted even if simulate front is selected.
		pr := []gamesystems.Roll{}
		if i > 0 {
			pr = w.WeatherResults[i-1].UnmodifiedRolls
		}

		results, _ := system.DailyWeather(wr.UnmodifiedRolls, pr, w.CurrentKoppenCode, w.CurrentWind, w.CurrentSeason, w.SimulateFront)
		mods := system.ListWeatherModifiers(w.CurrentKoppenCode, w.CurrentSeason)

		wr.Modifiers = mods
		wr.TemperatureDay = results[0]
		wr.TemperatureNight = results[1]
		wr.Precipitation = results[2]
		wr.Wind =results[3]
		wr.WindDirection = results[4]
	}
}

func (w WeatherForm) RerollWeatherDice() {
	system, _ := gamesystems.LoadGameSystem(w.currentSystem)
	for _, wr := range w.WeatherResults {
		// Setting the results to any empty array will tell the system that the first day has no previous rolls.
		// This means the first day will not be adjusted even if simulate front is selected.
		rolls := system.ListWeatherRolls()
		for _, r := range rolls {
			r.Resolve()
		}
		wr.UnmodifiedRolls = rolls
	}
}

func (w WeatherForm) CSVSummary() string {
	output := fmt.Sprint("Day Temperature,Night Temperature,Precipitation,Wind Strength,Wind Direction\n")
	for _, wt := range w.WeatherResults {
		output += fmt.Sprintf("%s,%s,%s,%s,%s\n", wt.TemperatureDay, wt.TemperatureNight, wt.Precipitation, wt.Wind, wt.WindDirection)
	}
	return output
}

func (w WeatherForm) MarkDownSummary() string {
	output := fmt.Sprintf("|%-25s|%-25s|%-25s|%-25s|%-25s|\n","Day Temperature","Night Temperature","Precipitation","Wind Strength","Wind Direction")
	output += fmt.Sprintf("|%s|%s|%s|%s|%s|\n", strings.Repeat("-", 25), strings.Repeat("-", 25), strings.Repeat("-", 25), strings.Repeat("-", 25), strings.Repeat("-", 25), )
	for _, wt := range w.WeatherResults {
		output += fmt.Sprintf("|%-25s|%-25s|%-25s|%-25s|%-25s|\n", wt.TemperatureDay, wt.TemperatureNight, wt.Precipitation, wt.Wind, wt.WindDirection)
	}

	return output
}

func buildSelectListFromStrings(optionStrings []string, selectedOption string) []*SelectOption {
	options := make([]*SelectOption, len(optionStrings))
	i := 0
	for _, s := range optionStrings {
		o := &SelectOption{
			Value: s,
			Text: s,
			Selected: false,
		}
		if selectedOption != "" && s == selectedOption {
			o.Selected = true
		}
		options[i] = o
		i++
	}
	return options
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
	sessionData map[string]*WeatherForm
	mu *sync.RWMutex
}

type SelectOption struct {
	Value, Text string
	Selected bool
}

func NewPageWeather() *PageWeather {
	return &PageWeather{
		templateName: "weatherGenerator.html",
		formTemplateName: "weatherForm.html",
		renderer: NewRenderer(),
		sessionData: map[string]*WeatherForm{},
		mu: &sync.RWMutex{},
	}
}


func (p PageWeather) addWeatherToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {
		session := getSessionFromCtx(r)
		if p.sessionData[session.Id] == nil {
			p.sessionData[session.Id] = NewWeatherForm()
		}
		a := p.sessionData[session.Id]
		ctx := context.WithValue(r.Context(), weatherKey, a)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (p PageWeather) getWeatherFromContext (r *http.Request) *WeatherForm {
	p.mu.RLock()
	a, _ := r.Context().Value(weatherKey).(*WeatherForm)
	p.mu.RUnlock()
	return a
}

func (p *PageWeather) CleanupSessionData(sessionId string) {
	p.mu.Lock()
	delete(p.sessionData, sessionId)
	p.mu.Unlock()
}


func (p PageWeather) RegisterRoutes(router *http.ServeMux) {
	router.Handle("/tools/weather", p.addWeatherToContext(p.index()))
	router.Handle("PUT /tools/weather", p.addWeatherToContext(p.reRollWeather(p.form())))

	router.Handle("PATCH /tools/weather/mode", p.addWeatherToContext(p.updateGenerationMode(p.form())))
	router.Handle("PATCH /tools/weather/code", p.addWeatherToContext(p.updateKoppenCode(p.form())))
	router.Handle("PATCH /tools/weather/season", p.addWeatherToContext(p.updateSeason(p.form())))
	router.Handle("PATCH /tools/weather/prevailing-wind", p.addWeatherToContext(p.updatePrevailingWind(p.form())))
	router.Handle("PATCH /tools/weather/front", p.addWeatherToContext(p.updateFrontSimulation(p.form())))

	router.Handle("GET /tools/printable-weather", p.printerpage())

	router.Handle("GET /tools/weather/summary", p.addWeatherToContext(p.weatherSummary()))
}


func (p PageWeather) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		wt := p.getWeatherFromContext(r)
		params := r.URL.Query()
		mode := params.Get("format")
		code := params.Get("code")
		season := params.Get("season")
		wind := params.Get("wind")
		front := params.Get("front")
		applySettings := false
		if ValidateSelectOption(wt.ModeOptions, mode) {
			wt.CurrentMode = mode
			SetOptionSelected(wt.ModeOptions, mode)
			applySettings = true

		}
		if ValidateSelectOption(wt.KoppenCodeOptions, code) {
			wt.CurrentKoppenCode = code
			SetOptionSelected(wt.KoppenCodeOptions, code)
			applySettings = true

		}
		if ValidateSelectOption(wt.SeasonOptions, season) {
			wt.CurrentSeason, _ = gamesystems.SeasonFromString(season)
			SetOptionSelected(wt.SeasonOptions, season)
			applySettings = true

		}
		if ValidateSelectOption(wt.WindOptions, wind) {
			wt.CurrentWind = wind
			SetOptionSelected(wt.WindOptions, wind)
			applySettings = true

		}
		if front == "true" {
			wt.SimulateFront = true
			applySettings = true
		}
		if applySettings {
			wt.ApplyWeatherSettings()
		}
		renderedPage, err := p.renderer.RenderPage(p.templateName, wt)
		if err != nil {
			log.Printf("Error rendering Index of the Weather Page: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}

func (p PageWeather) form() http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		weather := p.getWeatherFromContext(r)
		weather.ApplyWeatherSettings()
		renderedPage, err := p.renderer.RenderPage(p.formTemplateName, weather)
		if err != nil {
			log.Printf("Error rendering Index of the Weather Page: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}


func (p PageWeather) weatherSummary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pdata struct {
			Text string
		} = struct {
			Text string
		} {
			Text: "",
		}
		wt := p.getWeatherFromContext(r)
	params := r.URL.Query()
	format := params.Get("format")
	if format == "" || format == "csv" {
		pdata.Text = wt.CSVSummary()
	}
		if format == "markdown" {
			pdata.Text = wt.MarkDownSummary()
		}
		renderedPage, err := p.renderer.RenderPage("weatherSummaryModal.html", pdata)
		if err != nil {
			log.Printf("Error rendering Summary of the Weather Page: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}

func (p PageWeather) printerpage() http.HandlerFunc {

	return func(w http.ResponseWriter, r* http.Request) {


		var printData struct {
			Weather []*Weather
			Valid bool
			Koppen, Season, PrevailingWind, Error string
		} = struct{
			Weather []*Weather
			Valid bool
			Koppen, Season, PrevailingWind, Error string
		} {
			Weather: []*Weather{},
			Valid: true,
		}

		session := getSessionFromCtx(r)
		var weather *WeatherForm
		if p.sessionData[session.Id] != nil {
			weather = p.sessionData[session.Id]
		}


		if weather == nil {
			printData.Error = "You have no existing weather data. Did you click the link on the generator page?"
			printData.Valid = false
			weather = NewWeatherForm()
		}

		printData.Koppen = weather.CurrentKoppenCode
		printData.PrevailingWind = weather.CurrentWind
		printData.Season = weather.CurrentSeason.String()
		printData.Weather = weather.WeatherResults

		renderedPage, err := p.renderer.RenderPage("printable-weather.html", printData)
		if err != nil {
			log.Printf("Error rendering Index of the Weather Page: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}

func SetOptionSelected(opts []*SelectOption, option string) bool {
		optionSet := false
		for _, o := range opts {
			o.Selected = false
			if o.Value == option {
				o.Selected = true
				optionSet = true
			}
		}
	return optionSet
}

func ValidateSelectOption(opts []*SelectOption, option string) bool {
		for _, o := range opts {
			if o.Value == option {
				return true
			}
		}
	return false
}

func (p PageWeather) updateGenerationMode(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wt := p.getWeatherFromContext(r);
		r.ParseForm()
		newMode := r.Form.Get("generation-mode")
		wt.CurrentMode = newMode
		_ = SetOptionSelected(wt.ModeOptions, newMode)
		next.ServeHTTP(w, r)
	}
}

func (p PageWeather) updateKoppenCode(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wt := p.getWeatherFromContext(r);
		r.ParseForm()
		o := r.Form.Get("koppen-code")
		wt.CurrentKoppenCode = o
		_ = SetOptionSelected(wt.KoppenCodeOptions, o)
		next.ServeHTTP(w, r)
	}
}

func (p PageWeather) updateSeason(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wt := p.getWeatherFromContext(r);
		r.ParseForm()
		o := r.Form.Get("season")
		wt.CurrentSeason, _ = gamesystems.SeasonFromString(o)
		_ = SetOptionSelected(wt.SeasonOptions, o)
		next.ServeHTTP(w, r)
	}
}

func (p PageWeather) updatePrevailingWind(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wt := p.getWeatherFromContext(r);
		r.ParseForm()
		o := r.Form.Get("prevailing-wind")
		wt.CurrentWind = o
		_ = SetOptionSelected(wt.WindOptions, o)

		next.ServeHTTP(w, r)
	}
}

func (p PageWeather) updateFrontSimulation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wt := p.getWeatherFromContext(r);
		r.ParseForm()
		o := r.Form.Get("front")
		wt.SimulateFront = o == "true"
		next.ServeHTTP(w, r)
	}
}


func (p PageWeather) reRollWeather(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wt := p.getWeatherFromContext(r);
		wt.RerollWeatherDice()
		next.ServeHTTP(w, r)
	}
}
