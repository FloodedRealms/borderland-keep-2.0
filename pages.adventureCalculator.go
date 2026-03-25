package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	gamesystems "github.com/FloodedRealms/borderland-keep-2.0/game-systems"
)

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
	Adventure *Adventure
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
		Adventure: &Adventure{
			SpecialTreasures: []SpecialTreasure{},
			Combat: []CombatEncounter{},
			MagicItems: []MagicItem{},
		},
		renderer: NewRenderer(),
		currentSystem: "ACKS II",
	}
}

func (p PageAdventureCalculator) RegisterRoutes(router *http.ServeMux) {
	router.Handle("/tools/calculator", p.index((*p.renderer)))
	router.Handle("PATCH /tools/calculator/shares", p.Shares(p.UpdateAdventureData(p.form(*p.renderer))) )
	router.Handle("PATCH /tools/calculator/coins", p.Coins(p.UpdateAdventureData(p.form(*p.renderer))) )

	router.Handle("POST /tools/calculator/special-treasure",p.AddSpecialTreasure(p.form(*p.renderer)))
	router.Handle("PATCH /tools/calculator/special-treasure/{id}",p.UpdateSpecialTreasure(p.UpdateAdventureData(p.form(*p.renderer))))
	router.Handle("DELETE /tools/calculator/special-treasure/{id}",p.DeleteSpecialTreasure(p.UpdateAdventureData(p.form(*p.renderer))))

    router.Handle("POST /tools/calculator/combat",p.AddCombat(p.form(*p.renderer)))
	router.Handle("PATCH /tools/calculator/combat/{id}",p.UpdateCombat(p.UpdateAdventureData(p.form(*p.renderer))))
	router.Handle("DELETE /calculator/combat/{id}",p.DeleteCombat(p.UpdateAdventureData(p.form(*p.renderer))))

    router.Handle("POST /tools/calculator/magic-item",p.AddMagicItem(p.form(*p.renderer)))
	router.Handle("PATCH /tools/calculator/magic-item/{id}",p.UpdateMagicItem(p.UpdateAdventureData(p.form(*p.renderer))))
	router.Handle("DELETE /tools/calculator/magic-item/{id}",p.DeleteMagicItem(p.UpdateAdventureData(p.form(*p.renderer))))

	router.Handle("/tools/calculator/summary", p.AdventureSummaryModal(*p.renderer))


}

func (p PageAdventureCalculator) index(pr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {
		renderedPage, err := pr.RenderPage(p.templateName, p.Adventure)
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
		renderedPage, err := pr.RenderPage(p.formTemplateName, p.Adventure)
		if err != nil {
			log.Printf("Error rendering Form of Adventure Calculator: %v\n", err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}


func (p PageAdventureCalculator) Shares(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		pcount, err := strconv.Atoi(r.Form.Get("player-count"))
		if err != nil {
			p.Adventure.Errors = []string{"Player count malformed"}
			next.ServeHTTP(w, r)
		}
		hcount, err := strconv.Atoi(r.Form.Get("hench-count"))
		if err != nil {
			p.Adventure.Errors = []string{"Henchmen count malformed"}
			next.ServeHTTP(w, r)
		}
		p.Adventure.NumberOfPlayers = pcount
		p.Adventure.NumberOfHenchmen = hcount
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) Coins(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		copper, err := strconv.Atoi(r.Form.Get("copper"))
		if err != nil {
			p.Adventure.Errors = []string{"Copper count malformed"}
			next.ServeHTTP(w, r)
		}
		silver, err := strconv.Atoi(r.Form.Get("silver"))
		if err != nil {
			p.Adventure.Errors = []string{"Silver count malformed"}
			next.ServeHTTP(w, r)
		}
		electrum, err := strconv.Atoi(r.Form.Get("electrum"))
		if err != nil {
			p.Adventure.Errors = []string{"Electrum count malformed"}
			next.ServeHTTP(w, r)
		}
		gold, err := strconv.Atoi(r.Form.Get("gold"))
		if err != nil {
			p.Adventure.Errors = []string{"Gold count malformed"}
			next.ServeHTTP(w, r)
		}
		platinum, err := strconv.Atoi(r.Form.Get("platinum"))
		if err != nil {
			p.Adventure.Errors = []string{"Platinum count malformed"}
			next.ServeHTTP(w, r)
		}
		p.Adventure.Copper = copper
		p.Adventure.Silver = silver
		p.Adventure.Electrum = electrum
		p.Adventure.Gold = gold
		p.Adventure.Platinum = platinum
		next.ServeHTTP(w, r)
	}
}

/* *** Special Treasure *** */

func (p PageAdventureCalculator) AddSpecialTreasure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.Adventure.SpecialTreasures = append(p.Adventure.SpecialTreasures, SpecialTreasure{})
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) UpdateSpecialTreasure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Broken Treasure ID: %v/n", err)
			p.Adventure.Errors = append(p.Adventure.Errors, err.Error())
			next.ServeHTTP(w, r)
			return
		}


		name := r.Form.Get(fmt.Sprintf("treasure-name-%d", id))
		numberLabel := fmt.Sprintf("treasure-collected-%d", id)
		number, err := strconv.Atoi(r.Form.Get(numberLabel))
		if err != nil {
			p.Adventure.Errors = append(p.Adventure.Errors,fmt.Sprintf("Treasure Number Broken: %v\n Label: %s", err, numberLabel))
			next.ServeHTTP(w, r)
			return
		}
		value, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("treasure-value-%d", id)))
		if err != nil {
			p.Adventure.Errors = append(p.Adventure.Errors,fmt.Sprintf("Treasure Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
			p.Adventure.SpecialTreasures[id].GPValue = value
			p.Adventure.SpecialTreasures[id].Number = number
			p.Adventure.SpecialTreasures[id].Name = name
			p.Adventure.SpecialTreasures[id].TotalValue = value * number

		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) DeleteSpecialTreasure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			next.ServeHTTP(w, r)
		}

		nt := []SpecialTreasure{}
		for i, t := range p.Adventure.SpecialTreasures {
			if i == id {
				continue
			}
			nt = append(nt, t)
		}
		p.Adventure.SpecialTreasures = nt
		next.ServeHTTP(w, r)
	}
}

/* *** Combat *** */

func (p PageAdventureCalculator) AddCombat(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
		p.Adventure.Combat = append(p.Adventure.Combat, CombatEncounter{})
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) UpdateCombat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Broken Combat ID: %v/n", err)
			p.Adventure.Errors = append(p.Adventure.Errors, err.Error())
			next.ServeHTTP(w, r)
			return
		}

		name := r.Form.Get(fmt.Sprintf("combat-name-%d", id))
		number, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("combat-defeated-%d", id)))
		if err != nil {
			p.Adventure.Errors = append(p.Adventure.Errors,fmt.Sprintf("Combat Number Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
		value, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("combat-value-%d", id)))
		if err != nil {
			p.Adventure.Errors = append(p.Adventure.Errors,fmt.Sprintf("Combat Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
			p.Adventure.Combat[id].XPValue = value
			p.Adventure.Combat[id].NumberDefeated = number
			p.Adventure.Combat[id].Name = name

		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) DeleteCombat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			next.ServeHTTP(w, r)
		}

		n := []CombatEncounter{}
		for i, t := range p.Adventure.Combat {
			if i == id {
				continue
			}
			n = append(n, t)
		}
		p.Adventure.Combat = n
		next.ServeHTTP(w, r)
	}
}

/* *** Magic Items *** */

func (p PageAdventureCalculator) AddMagicItem(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
		p.Adventure.MagicItems = append(p.Adventure.MagicItems, MagicItem{})
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) UpdateMagicItem(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Broken Magic Item ID: %v/n", err)
			p.Adventure.Errors = append(p.Adventure.Errors, err.Error())
			next.ServeHTTP(w, r)
			return
		}


		name := r.Form.Get(fmt.Sprintf("mi-name-%d", id))
		ap, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("mi-apparent-value-%d", id)))
		if err != nil {
			p.Adventure.Errors = append(p.Adventure.Errors,fmt.Sprintf("Apparent Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
		value, err := strconv.Atoi(r.Form.Get(fmt.Sprintf("mi-sell-value-%d", id)))
		if err != nil {
			p.Adventure.Errors = append(p.Adventure.Errors,fmt.Sprintf("Sell Value Broken: %v\n", err))
			next.ServeHTTP(w, r)
			return
		}
		isSold := r.Form.Get(fmt.Sprintf("mi-is-sold-%d", id)) == "true"
			p.Adventure.MagicItems[id].ApparentValue = ap
			p.Adventure.MagicItems[id].SellValue = value
			p.Adventure.MagicItems[id].Name = name
			p.Adventure.MagicItems[id].IsSold = isSold

		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) DeleteMagicItem(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			next.ServeHTTP(w, r)
		}

		n := []MagicItem{}
		for i, t := range p.Adventure.MagicItems {
			if i == id {
				continue
			}
			n = append(n, t)
		}
		p.Adventure.MagicItems = n
		next.ServeHTTP(w, r)
	}
}

/* Adventure Data and Summary */

func (p PageAdventureCalculator) UpdateAdventureData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		system, err := gamesystems.LoadGameSystem(p.currentSystem)
		if err != nil {
			log.Printf("%v", err)
		}
		p.Adventure.TotalShares = system.CalculateNumberOfShares(p.Adventure.NumberOfPlayers, p.Adventure.NumberOfHenchmen)
		totalGPFromCoins := system.CalculateTotalGPFromCoinage(p.Adventure.Copper, p.Adventure.Silver, p.Adventure.Electrum, p.Adventure.Gold, p.Adventure.Platinum)
		stn, stv := p.Adventure.SpecialTreasureArrays()
		cbn, cbv := p.Adventure.CombatArrays()
		mgav, mgsv, mgis := p.Adventure.MagicItemArrays()
		fs, hs := system.CalculateXPShares(p.Adventure.TotalShares, p.Adventure.Copper, p.Adventure.Silver, p.Adventure.Electrum, p.Adventure.Gold, p.Adventure.Platinum, stn, stv, cbn, cbv, mgav, mgsv, mgis)

		p.Adventure.TotalXPValue = system.CalculateTotalXP(p.Adventure.Copper, p.Adventure.Silver, p.Adventure.Electrum, p.Adventure.Gold, p.Adventure.Platinum, stn, stv, cbn, cbv, mgav, mgsv, mgis)
		p.Adventure.TotalCoinGold = totalGPFromCoins
		p.Adventure.TotalGPValue = system.CalculateTotalGP(p.Adventure.Copper, p.Adventure.Silver, p.Adventure.Electrum, p.Adventure.Gold, p.Adventure.Platinum, stn, stv, mgav, mgsv, mgis)
		p.Adventure.FullXPShare = fs
		p.Adventure.HalfXPShare = hs
		next.ServeHTTP(w, r)
	}
}

func (p PageAdventureCalculator) AdventureSummaryModal(pr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		system, err := gamesystems.LoadGameSystem(p.currentSystem)
		if err != nil {
			log.Printf("%v", err)
		}
		stn, stv := p.Adventure.SpecialTreasureArrays()
		cbn, cbv := p.Adventure.CombatArrays()
		mgav, mgsv, mgis := p.Adventure.MagicItemArrays()
		fs, hs := system.CalculateXPShares(p.Adventure.TotalShares, p.Adventure.Copper, p.Adventure.Silver, p.Adventure.Electrum, p.Adventure.Gold, p.Adventure.Platinum, stn, stv, cbn, cbv, mgav, mgsv, mgis)
		fgs, hgs := system.CalculateGPShares(p.Adventure.TotalShares, p.Adventure.Copper, p.Adventure.Silver, p.Adventure.Electrum, p.Adventure.Gold, p.Adventure.Platinum, stn, stv, mgav, mgsv, mgis)
		coins := system.CalculateDetailedCoinage(p.Adventure.TotalShares, p.Adventure.Copper, p.Adventure.Silver, p.Adventure.Electrum, p.Adventure.Gold, p.Adventure.Platinum)
		pCoinValue := system.CalculateTotalGPFromCoinage(coins[0], coins[2], coins[4], coins[6], coins[8])
		hCoinValue := system.CalculateTotalGPFromCoinage(coins[1], coins[3], coins[5], coins[7], coins[9])
		pSTAmount, hSTAmount := system.CalculateGPSharesFromSpecialTreasure(stn, stv, p.Adventure.TotalShares)
		pMIAmount, hMIAmount := system.CalculateGPSharesFromMagicItems(mgav, mgsv, mgis, p.Adventure.TotalShares)

		summary := AdventureSummary {
				PlayerCount: p.Adventure.NumberOfPlayers,
				HenchCount: p.Adventure.NumberOfHenchmen,
				Shares: p.Adventure.TotalShares,
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
				SpecialTreasures: p.Adventure.SpecialTreasures,
				MagicItems: p.Adventure.MagicItems,
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
