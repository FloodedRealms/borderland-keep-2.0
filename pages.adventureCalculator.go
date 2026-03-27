package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	gamesystems "github.com/FloodedRealms/borderland-keep-2.0/game-systems"
)


const adventureKey contextKey = "adventure"

type Adventure struct {
	NumberOfPlayers int
	NumberOfHenchmen int
	Gold int
	Silver int
	Copper int
	Electrum int
	Platinum int
	TotalCoinGold float64
	SpecialTreasures []SpecialTreasure
	MagicItems []MagicItem
	Combat []CombatEncounter
	TotalShares float64

	TotalGPValue int
	TotalXPValue int
	FullXPShare int
	FullGPShare int
	HalfXPShare int
	HalfGPShare int
	Errors []string
}


func NewAdventure() *Adventure {
	return &Adventure{
		SpecialTreasures: []SpecialTreasure{},
		MagicItems: []MagicItem{},
		Combat: []CombatEncounter{},
	}
}

func ParseAndValidateAdventure(r *http.Request) (*Adventure, bool) {
	r.ParseForm()
	a := &Adventure{
		SpecialTreasures: []SpecialTreasure{},
		MagicItems: []MagicItem{},
		Combat: []CombatEncounter{},
		Errors: []string{},
	}
	validated := true

	errors := []string{}
	r.ParseForm()
	pcount, _, errors := ValidatePositiveIntegerField(r.Form.Get("player-count"), "Player Count", errors)
	hcount, _, errors := ValidatePositiveIntegerField(r.Form.Get("hench-count"), "Henchmen Count", errors)
	copper, _, errors := ValidatePositiveIntegerField(r.Form.Get("copper"), "Copper", errors)
	silver, _, errors := ValidatePositiveIntegerField(r.Form.Get("silver"), "Silver", errors)
	electrum, _, errors := ValidatePositiveIntegerField(r.Form.Get("electrum"), "Electrum", errors)
	gold, _, errors := ValidatePositiveIntegerField(r.Form.Get("gold"), "Gold", errors)
	platinum, _, errors := ValidatePositiveIntegerField(r.Form.Get("platinum"), "Platinum", errors)
	ts, _, errors := ParseandValidateTreasure(r.Form, errors)
	cs, _, errors := ParseandValidateCombat(r.Form, errors)
	ms, _, errors := ParseandValidateMagicItems(r.Form, errors)


	validated = len(errors) == 0
	a.Errors = errors
	a.NumberOfPlayers = pcount
	a.NumberOfHenchmen = hcount
	a.Copper = copper
	a.Silver = silver
	a.Electrum = electrum
	a.Gold = gold
	a.Platinum = platinum
	a.SpecialTreasures = ts
	a.Combat = cs
	a.MagicItems = ms

	return a, validated
}

func ParseandValidateTreasure(form url.Values, e []string) ([]SpecialTreasure, bool, []string) {

	tkey := regexp.MustCompile(`^treasure-(\w+)-(\d+)`)
	tmap := map[int]*SpecialTreasure{}

	for key, vals := range form {
		matches := tkey.FindStringSubmatch(key)
		if matches == nil {
			continue
		}
		index, _ := strconv.Atoi(matches[2])
		field := matches[1]
		val := vals[0]

		if tmap[index] == nil {
			tmap[index] = &SpecialTreasure{}
		}

		switch field {
		case "name":
			tmap[index].Name = val
		case "value":
			gp, _, es := ValidatePositiveIntegerField(val, fmt.Sprintf("GP value of treasure %s", index), e)
			e = append(e, es...)
			tmap[index].GPValue = gp
		case "collected":
			gp, _, es := ValidatePositiveIntegerField(val, fmt.Sprintf("Number of treasure %s", index), e)
			e = append(e, es...)
			tmap[index].Number = gp
		}
	}

	ts := make([]SpecialTreasure, len(tmap))
	for i, t := range tmap {
		ts[i] = *t
	}
	v := len(e) == 0
	return ts, v, e
}

func ParseandValidateCombat(form url.Values, e []string) ([]CombatEncounter, bool, []string) {

	tkey := regexp.MustCompile(`^combat-(\w+)-(\d+)`)
	tmap := map[int]*CombatEncounter{}

	for key, vals := range form {
		matches := tkey.FindStringSubmatch(key)
		if matches == nil {
			continue
		}
		index, _ := strconv.Atoi(matches[2])
		field := matches[1]
		val := vals[0]

		if tmap[index] == nil {
			tmap[index] = &CombatEncounter{}
		}

		switch field {
		case "name":
			tmap[index].Name = val
		case "value":
			v, _, es := ValidatePositiveIntegerField(val, fmt.Sprintf("XP value of combat %s", index), e)
			e = append(e, es...)
			tmap[index].XPValue = v
		case "defeated":
			v, _, es := ValidatePositiveIntegerField(val, fmt.Sprintf("Number of combet %s", index), e)
			e = append(e, es...)
			tmap[index].NumberDefeated = v
		}
	}

	ts := make([]CombatEncounter, len(tmap))
	for i, v := range tmap {
		ts[i] = *v
	}
	v := len(e) == 0
	return ts, v, e
}

func ParseandValidateMagicItems(form url.Values, e []string) ([]MagicItem, bool, []string) {

	tkey := regexp.MustCompile(`^mi-(\w+)-(\d+)`)
	tmap := map[int]*MagicItem{}

	for key, vals := range form {
		matches := tkey.FindStringSubmatch(key)
		if matches == nil {
			continue
		}
		index, _ := strconv.Atoi(matches[2])
		field := matches[1]
		val := vals[0]

		if tmap[index] == nil {
			tmap[index] = &MagicItem{}
		}

		switch field {
		case "name":
			tmap[index].Name = val
		case "apparent":
			v, _, es := ValidatePositiveIntegerField(val, fmt.Sprintf("Apparent Value of Magic Item %s", index), e)
			e = append(e, es...)
			tmap[index].ApparentValue = v
		case "sell":
			v, _, es := ValidatePositiveIntegerField(val, fmt.Sprintf("Sell value of Magic Item %s", index), e)
			e = append(e, es...)
			tmap[index].SellValue = v
		case "sold":
            v := val == "true"
			tmap[index].IsSold = v

		}
	}

	ts := make([]MagicItem, len(tmap))
	for i, v := range tmap {
		ts[i] = *v
	}
	v := len(e) == 0
	return ts, v, e
}

func ValidatePositiveIntegerField(input, fieldLabel string, errors []string) (int, bool, []string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return value, false, append(errors, fmt.Sprintf("%s should be a number!", fieldLabel))
	}
	if value < 0 {
		return value, false, append(errors, fmt.Sprintf("%s should be at least 0!", fieldLabel))
	}
	return value, true, errors
}

type AdventureSummary struct {
	PlayerCount int
	HenchCount int
	Shares float64
	PlayerXP int
	HenchXP int
	PlayerGP int
	HenchGP int
	PlayerCoinValue int
	HenchCoinValue int
	PCP, PSP, PEP, PGP, PPP int
	HCP, HSP, HEP, HGP, HPP int
	PlayerSpecialTreasure int
	HenchSpecialTreasure int
	SpecialTreasures []SpecialTreasure
	MagicItems []MagicItem

}

type SpecialTreasureType string
const (
	Jewelery SpecialTreasureType = "Jewelery"
	Gemstone SpecialTreasureType = "Gemstone"
)

type SpecialTreasure struct {
	GPValue int
	SPType SpecialTreasureType
	Number int
	Name string
	TotalValue int
}

type MagicItem struct {
	ApparentValue int
	IsSold bool
	SellValue int
	Name string
}

type CombatEncounter struct {
	Name string
	NumberDefeated int
	XPValue int
}

type PageAdventureCalculator struct {
	templateName string
	formTemplateName string
	errorTemplate string
	adventures map[string]*Adventure
	currentSystem string
	renderer *Renderer
}

func (a Adventure) SpecialTreasureArrays() (numberRetrieved, value []int) {
	numberRetrieved, value = []int{}, []int{}
	for _, v := range a.SpecialTreasures {
		numberRetrieved = append(numberRetrieved, v.Number)
		value = append(value, v.GPValue)
	}
	return
}

func (a Adventure) CombatArrays() (numberDefeated, value []int) {
	numberDefeated, value = []int{}, []int{}
	for _, v := range a.Combat {
		numberDefeated = append(numberDefeated, v.NumberDefeated)
		value = append(value, v.XPValue)
	}
	return
}

func (a Adventure) MagicItemArrays() (av, sv []int, isSold []bool) {
	av, sv = []int{}, []int{}
	isSold = []bool{}

	for _, v := range a.MagicItems {
		av = append(av, v.ApparentValue)
		sv = append(sv, v.SellValue)
		isSold = append(isSold, v.IsSold)
	}
	return
}

func NewPageAdventureCalculator() *PageAdventureCalculator {
	return &PageAdventureCalculator{
		templateName: "adventureCalculator.html",
		formTemplateName: "calculatorForm.html",
		errorTemplate: "calculatorError.html",
		adventures: map[string]*Adventure{},
		renderer: NewRenderer(),
		currentSystem: "ACKS II",
	}
}

func (p PageAdventureCalculator) RegisterRoutes(router *http.ServeMux) {
	router.Handle("/tools/calculator", p.addAventureToContext(p.index((*p.renderer))))
	router.Handle("PATCH /tools/calculator/shares", p.addAventureToContext(p.Shares(p.UpdateAdventureData(p.form(*p.renderer))) ))
	router.Handle("PATCH /tools/calculator/coins", p.addAventureToContext(p.Coins(p.UpdateAdventureData(p.form(*p.renderer))) ))

	router.Handle("POST /tools/calculator/special-treasure",p.addAventureToContext(p.AddSpecialTreasure(p.form(*p.renderer))))
	router.Handle("PATCH /tools/calculator/special-treasure/{id}",p.addAventureToContext(p.UpdateSpecialTreasure(p.UpdateAdventureData(p.form(*p.renderer)))))
	router.Handle("DELETE /tools/calculator/special-treasure/{id}",p.addAventureToContext(p.DeleteSpecialTreasure(p.UpdateAdventureData(p.form(*p.renderer)))))

    router.Handle("POST /tools/calculator/combat",p.addAventureToContext(p.AddCombat(p.form(*p.renderer))))
	router.Handle("PATCH /tools/calculator/combat/{id}",p.addAventureToContext(p.UpdateCombat(p.UpdateAdventureData(p.form(*p.renderer)))))
	router.Handle("DELETE /calculator/combat/{id}",p.addAventureToContext(p.DeleteCombat(p.UpdateAdventureData(p.form(*p.renderer)))))

    router.Handle("POST /tools/calculator/magic-item",p.addAventureToContext(p.AddMagicItem(p.form(*p.renderer))))
	router.Handle("PATCH /tools/calculator/magic-item/{id}",p.addAventureToContext(p.UpdateMagicItem(p.UpdateAdventureData(p.form(*p.renderer)))))
	router.Handle("DELETE /tools/calculator/magic-item/{id}",p.addAventureToContext(p.DeleteMagicItem(p.UpdateAdventureData(p.form(*p.renderer)))))

	router.Handle("/tools/calculator/summary", p.addAventureToContext(p.AdventureSummaryModal(*p.renderer)))


}

func (p PageAdventureCalculator) addAventureToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {
		session := getSessionFromCtx(r)
		if p.adventures[session.Id] == nil {
			p.adventures[session.Id] = &Adventure{}
		}
		a := p.adventures[session.Id]
		ctx := context.WithValue(r.Context(), adventureKey, a)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (p PageAdventureCalculator) getAdventureFromContext (r *http.Request) *Adventure {
	a, _ := r.Context().Value(adventureKey).(*Adventure)
	return a
}

func (p PageAdventureCalculator) index(pr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		a := p.getAdventureFromContext(r)
		renderedPage, err := pr.RenderPage(p.templateName, a)
		if err != nil {
			log.Printf("Error rendering Index of Adventure Calculator: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}

func (p PageAdventureCalculator) form(pr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		a := p.getAdventureFromContext(r)

		renderedPage, err := pr.RenderPage(p.formTemplateName, a)
		if err != nil {
			log.Printf("Error rendering Form of Adventure Calculator: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}

func (p PageAdventureCalculator) renderForm(a Adventure, valid bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		if !valid {
		renderedSnippet, err := p.renderer.RenderPage(p.errorTemplate, a)
			if err != nil {
				log.Printf("Error rendering Form of Adventure Calculator: %v\n", err)
				w.Header().Add("hx-redirect", "/error")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write([]byte(renderedSnippet))
			return

		}
		renderedPage, err := p.renderer.RenderPage(p.formTemplateName, a)
		if err != nil {
			log.Printf("Error rendering Form of Adventure Calculator: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(renderedPage))
	}
}


func (p PageAdventureCalculator) Shares(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		a := p.getAdventureFromContext(r)

		pcount, err := strconv.Atoi(r.Form.Get("player-count"))
		if err != nil {
			a.Errors = []string{"Player count malformed"}
			next.ServeHTTP(w, r)
		}
		hcount, err := strconv.Atoi(r.Form.Get("hench-count"))
		if err != nil {
			a.Errors = []string{"Henchmen count malformed"}
			next.ServeHTTP(w, r)
		}
		a.NumberOfPlayers = pcount
		a.NumberOfHenchmen = hcount
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) Coins(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		a := p.getAdventureFromContext(r)

		copper, err := strconv.Atoi(r.Form.Get("copper"))
		if err != nil {
			a.Errors = []string{"Copper count malformed"}
			next.ServeHTTP(w, r)
		}
		silver, err := strconv.Atoi(r.Form.Get("silver"))
		if err != nil {
			a.Errors = []string{"Silver count malformed"}
			next.ServeHTTP(w, r)
		}
		electrum, err := strconv.Atoi(r.Form.Get("electrum"))
		if err != nil {
			a.Errors = []string{"Electrum count malformed"}
			next.ServeHTTP(w, r)
		}
		gold, err := strconv.Atoi(r.Form.Get("gold"))
		if err != nil {
			a.Errors = []string{"Gold count malformed"}
			next.ServeHTTP(w, r)
		}
		platinum, err := strconv.Atoi(r.Form.Get("platinum"))
		if err != nil {
			a.Errors = []string{"Platinum count malformed"}
			next.ServeHTTP(w, r)
		}
		a.Copper = copper
		a.Silver = silver
		a.Electrum = electrum
		a.Gold = gold
		a.Platinum = platinum
		next.ServeHTTP(w, r)
	}
}

/* *** Special Treasure *** */

func (p PageAdventureCalculator) AddSpecialTreasure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := p.getAdventureFromContext(r)

		a.SpecialTreasures = append(a.SpecialTreasures, SpecialTreasure{})
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) UpdateSpecialTreasure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		a := p.getAdventureFromContext(r)

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Broken Treasure ID: %v/n", err)
			a.Errors = append(a.Errors, err.Error())
			next.ServeHTTP(w, r)
			return
		}


		name := r.Form.Get(fmt.Sprintf("treasure-name-%d", id))
		numberLabel := fmt.Sprintf("treasure-collected-%d", id)
		number, err := strconv.Atoi(r.Form.Get(numberLabel))
		if err != nil {
			a.Errors = append(a.Errors,fmt.Sprintf("Treasure Number Broken: %v\n Label: %s", err, numberLabel))
			next.ServeHTTP(w, r)
			return
		}
		value, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("treasure-value-%d", id)))
		if err != nil {
			a.Errors = append(a.Errors,fmt.Sprintf("Treasure Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
			a.SpecialTreasures[id].GPValue = value
			a.SpecialTreasures[id].Number = number
			a.SpecialTreasures[id].Name = name
			a.SpecialTreasures[id].TotalValue = value * number

		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) DeleteSpecialTreasure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := p.getAdventureFromContext(r)

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			next.ServeHTTP(w, r)
		}

		nt := []SpecialTreasure{}
		for i, t := range a.SpecialTreasures {
			if i == id {
				continue
			}
			nt = append(nt, t)
		}
		a.SpecialTreasures = nt
		next.ServeHTTP(w, r)
	}
}

/* *** Combat *** */

func (p PageAdventureCalculator) AddCombat(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
		a := p.getAdventureFromContext(r)

		a.Combat = append(a.Combat, CombatEncounter{})
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) UpdateCombat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		a := p.getAdventureFromContext(r)

		if err != nil {

			log.Printf("Broken Combat ID: %v/n", err)
			a.Errors = append(a.Errors, err.Error())
			next.ServeHTTP(w, r)
			return
		}

		name := r.Form.Get(fmt.Sprintf("combat-name-%d", id))
		number, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("combat-defeated-%d", id)))
		if err != nil {
			a.Errors = append(a.Errors,fmt.Sprintf("Combat Number Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
		value, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("combat-value-%d", id)))
		if err != nil {
			a.Errors = append(a.Errors,fmt.Sprintf("Combat Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
			a.Combat[id].XPValue = value
			a.Combat[id].NumberDefeated = number
			a.Combat[id].Name = name

		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) DeleteCombat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		a := p.getAdventureFromContext(r)

		if err != nil {
			next.ServeHTTP(w, r)
		}
		n := []CombatEncounter{}
		for i, t := range a.Combat {
			if i == id {
				continue
			}
			n = append(n, t)
		}
		a.Combat = n
		next.ServeHTTP(w, r)
	}
}

/* *** Magic Items *** */

func (p PageAdventureCalculator) AddMagicItem(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
		a := p.getAdventureFromContext(r)

		a.MagicItems = append(a.MagicItems, MagicItem{})
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) UpdateMagicItem(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		a := p.getAdventureFromContext(r)

		if err != nil {
			log.Printf("Broken Magic Item ID: %v/n", err)
			a.Errors = append(a.Errors, err.Error())
			next.ServeHTTP(w, r)
			return
		}


		name := r.Form.Get(fmt.Sprintf("mi-name-%d", id))
		ap, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("mi-apparent-value-%d", id)))
		if err != nil {
			a.Errors = append(a.Errors,fmt.Sprintf("Apparent Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
		value, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("mi-sell-value-%d", id)))
		if err != nil {
			a.Errors = append(a.Errors,fmt.Sprintf("Sell Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
		isSold := r.Form.Get(fmt.Sprintf("mi-is-sold-%d", id)) == "true"
			a.MagicItems[id].ApparentValue = ap
			a.MagicItems[id].SellValue = value
			a.MagicItems[id].Name = name
			a.MagicItems[id].IsSold = isSold

		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) DeleteMagicItem(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		a := p.getAdventureFromContext(r)

		if err != nil {
			next.ServeHTTP(w, r)
		}

		n := []MagicItem{}
		for i, t := range a.MagicItems {
			if i == id {
				continue
			}
			n = append(n, t)
		}
		a.MagicItems = n
		next.ServeHTTP(w, r)
	}
}

/* Adventure Data and Summary */

func (p PageAdventureCalculator) UpdateAdventureData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		system, err := gamesystems.LoadGameSystem(p.currentSystem)
		a := p.getAdventureFromContext(r)

		if err != nil {
			log.Printf("%v", err)
		}
		a.TotalShares = system.CalculateNumberOfShares(a.NumberOfPlayers, a.NumberOfHenchmen)
		totalGPFromCoins := system.CalculateTotalGPFromCoinage(a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum)
		stn, stv := a.SpecialTreasureArrays()
		cbn, cbv := a.CombatArrays()
		mgav, mgsv, mgis := a.MagicItemArrays()
		fs, hs := system.CalculateXPShares(a.TotalShares, a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum, stn, stv, cbn, cbv, mgav, mgsv, mgis)

		a.TotalXPValue = system.CalculateTotalXP(a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum, stn, stv, cbn, cbv, mgav, mgsv, mgis)
		a.TotalCoinGold = totalGPFromCoins
		a.TotalGPValue = system.CalculateTotalGP(a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum, stn, stv, mgav, mgsv, mgis)
		a.FullXPShare = fs
		a.HalfXPShare = hs
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) AdventureSummaryModal(pr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := p.getAdventureFromContext(r)

		system, err := gamesystems.LoadGameSystem(p.currentSystem)
		if err != nil {
			log.Printf("%v", err)
		}
		stn, stv := a.SpecialTreasureArrays()
		cbn, cbv := a.CombatArrays()
		mgav, mgsv, mgis := a.MagicItemArrays()
		fs, hs := system.CalculateXPShares(a.TotalShares, a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum, stn, stv, cbn, cbv, mgav, mgsv, mgis)
		fgs, hgs := system.CalculateGPShares(a.TotalShares, a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum, stn, stv, mgav, mgsv, mgis)
		coins := system.CalculateDetailedCoinage(a.TotalShares, a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum)
		pCoinValue := system.CalculateTotalGPFromCoinage(coins[0], coins[2], coins[4], coins[6], coins[8])
		hCoinValue := system.CalculateTotalGPFromCoinage(coins[1], coins[3], coins[5], coins[7], coins[9])
		pSTAmount, hSTAmount := system.CalculateGPSharesFromSpecialTreasure(stn, stv, a.TotalShares)
		pMIAmount, hMIAmount := system.CalculateGPSharesFromMagicItems(mgav, mgsv, mgis, a.TotalShares)

		summary := AdventureSummary {
				PlayerCount: a.NumberOfPlayers,
				HenchCount: a.NumberOfHenchmen,
				Shares: a.TotalShares,
				PlayerXP: fs,
				HenchXP: hs,
				PlayerGP: fgs,
				HenchGP: hgs,
				PCP: coins[0],
				PSP: coins[2],
				PEP: coins[4],
				PGP: coins[6],
				PPP: coins[8],
				HCP: coins[1],
				HSP: coins[3],
				HEP: coins[5],
				HGP: coins[7],
				HPP: coins[9],
				PlayerCoinValue: int(pCoinValue),
				HenchCoinValue: int(hCoinValue),
				PlayerSpecialTreasure: pSTAmount + pMIAmount,
				HenchSpecialTreasure: hSTAmount + hMIAmount,
				SpecialTreasures: a.SpecialTreasures,
				MagicItems: a.MagicItems,
		}
		renderedPage, err := pr.RenderPage("adventureSummaryModal.html", summary)
		if err != nil {
			log.Printf("Error rendering Adventure Summary Modal: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(renderedPage))
	}
}
