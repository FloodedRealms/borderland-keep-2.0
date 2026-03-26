package gamesystems

import (
	"errors"
	"fmt"
	"math"
)

type ACKSII struct {
	goldValueLookupTable []float64
	henchmenShare float64
	koppenModifiers  map[string][4]WeatherModifiers
	winds []string
}

func NewACKSII() ACKSII {

var KoppenModifiers = map[string][4]WeatherModifiers{
	// [Winter, Spring, Summer, Fall]

	// ── Tropical ──────────────────────────────────────────────────────────────
	"Af": {
		{TempDayMod: +3, TempNightMod: +0, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -1, WindMod: +0},
	},
	"Am": {
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -1, WindMod: +2},
		{TempDayMod: +7, TempNightMod: +0, PrecipMod: -1, WindMod: +0},
		{TempDayMod: +6, TempNightMod: +0, PrecipMod: +4, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: -1, WindMod: +0},
	},
	"Aw": {
		{TempDayMod: +5, TempNightMod: +0, PrecipMod: -5, WindMod: +2},
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +0, PrecipMod: -1, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
	},
	"As": {
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: -1, WindMod: +2},
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: -5, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
	},

	// ── Arid ──────────────────────────────────────────────────────────────────
	"BWh": {
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -5, WindMod: +2},
		{TempDayMod: +6, TempNightMod: +0, PrecipMod: -5, WindMod: +0},
		{TempDayMod: +7, TempNightMod: +2, PrecipMod: -5, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: +0, WindMod: +0},
	},
	"BWk": {
		{TempDayMod: +0, TempNightMod: -4, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -5, WindMod: +0},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +0, TempNightMod: -1, PrecipMod: -5, WindMod: +0},
	},
	"BSh": {
		{TempDayMod: +6, TempNightMod: +0, PrecipMod: -5, WindMod: +2},
		{TempDayMod: +7, TempNightMod: +3, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +6, TempNightMod: +2, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +6, TempNightMod: +2, PrecipMod: -5, WindMod: +0},
	},
	"BSk": {
		{TempDayMod: +0, TempNightMod: -1, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +0, PrecipMod: -5, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
	},

	// ── Temperate ─────────────────────────────────────────────────────────────
	"Csa": {
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
	},
	"Csb": {
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -2, WindMod: +2},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -1, WindMod: +0},
	},
	"Csc": {
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: +2, WindMod: +2},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
	},
	"Cwa": {
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +3, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +2, PrecipMod: -1, WindMod: +2},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
	},
	"Cwb": {
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -1, WindMod: +2},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -5, WindMod: +0},
	},
	"Cwc": {
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: +0, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: +2, WindMod: +2},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
	},
	"Cfa": {
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +1, PrecipMod: -1, WindMod: +0},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
	},
	"Cfb": {
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -1, WindMod: +2},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -1, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +1, PrecipMod: -1, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
	},
	"Cfc": {
		{TempDayMod: +0, TempNightMod: -1, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
	},

	// ── Continental ───────────────────────────────────────────────────────────
	"Dsa": {
		{TempDayMod: +0, TempNightMod: -2, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +3, TempNightMod: +1, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +5, TempNightMod: +3, PrecipMod: -5, WindMod: +0},
		{TempDayMod: +2, TempNightMod: +1, PrecipMod: -4, WindMod: +0},
	},
	"Dsb": {
		{TempDayMod: -2, TempNightMod: -5, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +1, TempNightMod: -2, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +1, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +0, TempNightMod: -1, PrecipMod: -3, WindMod: +0},
	},
	"Dsc": {
		{TempDayMod: -2, TempNightMod: -3, PrecipMod: -2, WindMod: +2},
		{TempDayMod: +1, TempNightMod: -1, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +2, TempNightMod: -1, PrecipMod: -2, WindMod: +0},
		{TempDayMod: -1, TempNightMod: -1, PrecipMod: -1, WindMod: +0},
	},
	"Dwa": {
		{TempDayMod: +0, TempNightMod: -3, PrecipMod: -5, WindMod: +2},
		{TempDayMod: +1, TempNightMod: +0, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -1, WindMod: +0},
	},
	"Dwb": {
		{TempDayMod: -1, TempNightMod: -3, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +1, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
	},
	"Dwc": {
		{TempDayMod: -1, TempNightMod: -5, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +0, TempNightMod: -3, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +1, TempNightMod: -1, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +0, TempNightMod: -2, PrecipMod: -3, WindMod: +0},
	},
	"Dwd": {
		{TempDayMod: -9, TempNightMod: -11, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +0, TempNightMod: -4, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +2, TempNightMod: -2, PrecipMod: -3, WindMod: +2},
		{TempDayMod: -4, TempNightMod: -7, PrecipMod: -5, WindMod: +0},
	},
	"Dfa": {
		{TempDayMod: -2, TempNightMod: -4, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +3, TempNightMod: -1, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +4, TempNightMod: +2, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +2, TempNightMod: -2, PrecipMod: -3, WindMod: +0},
	},
	"Dfb": {
		{TempDayMod: -3, TempNightMod: -4, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +2, TempNightMod: -2, PrecipMod: -2, WindMod: +0},
		{TempDayMod: +3, TempNightMod: +1, PrecipMod: -3, WindMod: +2},
		{TempDayMod: +0, TempNightMod: -2, PrecipMod: -3, WindMod: +0},
	},
	"Dfc": {
		{TempDayMod: -3, TempNightMod: -5, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +1, TempNightMod: -2, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +2, TempNightMod: +0, PrecipMod: -3, WindMod: +0},
		{TempDayMod: -1, TempNightMod: -2, PrecipMod: -3, WindMod: +0},
	},
	"Dfd": {
		{TempDayMod: -9, TempNightMod: -11, PrecipMod: -4, WindMod: +2},
		{TempDayMod: +0, TempNightMod: -4, PrecipMod: -4, WindMod: +0},
		{TempDayMod: +2, TempNightMod: -2, PrecipMod: -4, WindMod: +0},
		{TempDayMod: -5, TempNightMod: -6, PrecipMod: -5, WindMod: +0},
	},

	// ── Polar ─────────────────────────────────────────────────────────────────
	"ET": {
		{TempDayMod: -8, TempNightMod: -11, PrecipMod: -2, WindMod: +2},
		{TempDayMod: +0, TempNightMod: -1, PrecipMod: -3, WindMod: +0},
		{TempDayMod: +0, TempNightMod: +0, PrecipMod: -2, WindMod: +0},
		{TempDayMod: -1, TempNightMod: -2, PrecipMod: -1, WindMod: +0},
	},
	"EF": {
		{TempDayMod: -11, TempNightMod: -11, PrecipMod: -5, WindMod: -2},
		{TempDayMod: -11, TempNightMod: -11, PrecipMod: -5, WindMod: +0},
		{TempDayMod: -11, TempNightMod: -11, PrecipMod: -5, WindMod: +0},
		{TempDayMod: -11, TempNightMod: -11, PrecipMod: -5, WindMod: +0},
	},
}
var Winds = []string{"Northerly", "Northeasterly", "Easterly", "Southeasterly", "Southerly", "SouthWesterly", "Westerly", "Northerwesterly"}
	return ACKSII{
		goldValueLookupTable: []float64{0.01, 0.1, 0.5, 1, 5},
		henchmenShare: 0.5,
		koppenModifiers: KoppenModifiers,
		winds: Winds,
	}
}

// GetModifiers returns the WeatherModifiers for the given Koppen climate code
// and season. Returns false if the code is not found.
func (a ACKSII) GetWeatherModifiers(koppenCode string, season Season) (WeatherModifiers, bool) {
	seasons, ok := a.koppenModifiers[koppenCode]
	if !ok {
		return WeatherModifiers{}, false
	}
	return seasons[season], true
}

/* Treasure and XP */

var goldValueLookupTable = []float64{0.01, 0.1, 0.5, 1, 5}

func (a ACKSII) copperToGold(copperamount int) int {
	return int(math.Floor(float64(copperamount) * a.goldValueLookupTable[0]))
}
func (a ACKSII) silverToGold(silver int) int {
	return int(math.Floor(float64(silver) * a.goldValueLookupTable[1]))
}
func (a ACKSII) electrumToGold(electrum int) int {
	return int(math.Floor(float64(electrum) * a.goldValueLookupTable[2]))
}
func (a ACKSII) platinumToGold(platium int) int {
	return int(math.Floor(float64(platium) * a.goldValueLookupTable[4]))
}

func (a ACKSII) CalculateNumberOfShares(players, henchmen int) float64 {
	return float64(players) + float64(henchmen) / 2.0
}


func (a ACKSII) CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum int) float64 {
	copperGoldValue := a.copperToGold(copper)
	silverGoldValue := a.silverToGold(silver)
	electrumGoldValue := a.electrumToGold(electrum)
	platinumGoldValue := a.platinumToGold(platinum)
	return float64(copperGoldValue + silverGoldValue + electrumGoldValue + gold + platinumGoldValue)
}

func tenthsRemainder(value, shares float64) float64 {
	quotient := value / shares
	return float64(int(math.Abs(quotient * 10)) % 10)
}

func hundredthsRemainder(value, shares float64) float64 {
	quotient := value / shares
	return float64(int(math.Abs(quotient * 100)) % 10)
}


func (a ACKSII) CalculateDetailedCoinage(shares float64, copper, silver, electrum, gold, platinum int) []int {
	if shares == 0.0 {
		return []int{0,0,0,0,0,0,0,0,0,0}
	}
	//copperAmount, silverAmount, electrumAmount, goldAmount, platinumAmount := 0.0,0.0,0.0,0.0,0.0
	platinumAmount := math.Floor(float64(platinum) / shares)
	goldAmount := math.Floor(float64(gold) / shares)
	electrumAmount := math.Floor(float64(electrum) / shares) + tenthsRemainder(float64(platinum), shares)
	silverAmount := math.Floor(float64(silver) / shares) + tenthsRemainder(float64(gold), shares)
	copperAmount := math.Floor(float64(copper) / shares) + tenthsRemainder(float64(silver), shares) + (5.0 * tenthsRemainder(float64(electrum), shares)) + hundredthsRemainder(float64(gold), shares) + (5.0 * hundredthsRemainder(float64(platinum), shares))

	return []int{int(copperAmount), int(copperAmount * a.henchmenShare), int(silverAmount), int(silverAmount * a.henchmenShare), int(electrumAmount), int(electrumAmount * a.henchmenShare), int(goldAmount), int(goldAmount * a.henchmenShare), int(platinumAmount), int(platinumAmount * a.henchmenShare)}
}

func (a ACKSII) CalculateTotalXPFromCombat(combatants, xpValue []int) float64 {
	total := 0.0
	for i, v := range combatants {
		total = total + float64(v * xpValue[i])
	}
	return total
}

func (a ACKSII) CalculateTotalGPFromSpecialTreasure(numberRetrieved, gpValue []int) float64 {
	total := 0.0
	for i, v := range numberRetrieved {
		total = total + float64(v * gpValue[i])
	}
	return total
}

func (a ACKSII) CalculateGPSharesFromSpecialTreasure(numberRetrieved, gpValue []int, shares float64) (int, int) {
	if shares == 0.0 {
		return 0,0
	}
	total := a.CalculateTotalGPFromSpecialTreasure(numberRetrieved, gpValue)
	return int(math.RoundToEven(total / shares)), int(math.RoundToEven((total / shares) * a.henchmenShare))
}


func (a ACKSII) CalculateTotalGPFromMagicItems(apparentValue, sellValue []int, isSold []bool) float64 {
	total := 0.0
	for i, v := range apparentValue {
		total = total + float64(v)
		if (isSold[i]) {
			total = total + float64(sellValue[i])
		}
	}
	return total
}

func (a ACKSII) CalculateGPSharesFromMagicItems(apparentValue, sellValue []int, isSold []bool, shares float64) (int, int) {
	if shares == 0.0 {
		return 0,0
	}
	total := a.CalculateTotalGPFromMagicItems(apparentValue, sellValue, isSold)
	return int(math.RoundToEven(total / shares)), int(math.RoundToEven((total / shares) * a.henchmenShare))
}

func (a ACKSII)	CalculateTotalXP(copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, combatantsDefeated, combatantXPValue, magicItemAV, magicItemSV []int, magicItemIsSold []bool) int {
	coinXP := a.CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum)
	specialTreasureXP := a.CalculateTotalGPFromSpecialTreasure(specialTreasureRetrieved, specialTreasureValue)
	combatXP := a.CalculateTotalXPFromCombat(combatantsDefeated, combatantXPValue)
	magicItemXP := a.CalculateTotalGPFromMagicItems(magicItemAV, magicItemSV, magicItemIsSold)
	return int(math.RoundToEven(coinXP + specialTreasureXP + combatXP + magicItemXP))
}

func (a ACKSII)	CalculateTotalGP(copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, magicItemAV, magicItemSV []int, magicItemIsSold []bool) int {
	coinXP := a.CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum)
	specialTreasureXP := a.CalculateTotalGPFromSpecialTreasure(specialTreasureRetrieved, specialTreasureValue)
	magicItemXP := a.CalculateTotalGPFromMagicItems(magicItemAV, magicItemSV, magicItemIsSold)
	return int(math.RoundToEven(coinXP + specialTreasureXP + magicItemXP))
}

func (a ACKSII)	CalculateXPShares(totalShares float64, copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, combatantsDefeated, combatantXPValue, magicItemAV, magicItemSV []int, magicItemIsSold []bool) (fullXPShare, henchmenShare int) {
	if totalShares == 0 {
		return 0, 0
	}
	coinXP := a.CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum)
	specialTreasureXP := a.CalculateTotalGPFromSpecialTreasure(specialTreasureRetrieved, specialTreasureValue)
	combatXP := a.CalculateTotalXPFromCombat(combatantsDefeated, combatantXPValue)
	magicItemXP := a.CalculateTotalGPFromMagicItems(magicItemAV, magicItemSV, magicItemIsSold)
	totalXP := coinXP + specialTreasureXP + combatXP + magicItemXP
	fullShare := math.RoundToEven(totalXP / totalShares)
	henchShare := fullShare * a.henchmenShare
	return int(fullShare), int(henchShare)
}

func (a ACKSII)	CalculateGPShares(totalShares float64, copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) (fullShare, henchmenShare int) {
	if totalShares == 0.0 {
		return 0,0
	}
	coinGP := a.CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum)
	specialTreasureGP := a.CalculateTotalGPFromSpecialTreasure(specialTreasureRetrieved, specialTreasureValue)
	magicItemGP := a.CalculateTotalGPFromMagicItems(magicItemAV, magicItemSV, magicItemIsSold)
	totalGP := coinGP + specialTreasureGP + magicItemGP
	fullGPShare := math.RoundToEven(totalGP / totalShares)
	henchShare := fullGPShare * a.henchmenShare
	return int(fullGPShare), int(henchShare)
}

/*
   func (a ACKSII)	CalulateDetailedPayment(totalShares, copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) (fullCopper, fullSilver, fullEletrum, fullGold, fullPlatinum int) {
	coinGP := a.CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum)
	specialTreasureGP := a.CalculateTotalGPFromSpecialTreasure(specialTreasureRetrieved, specialTreasureValue)
	magicItemGP := a.CalculateTotalGPFromMagicItems(magicItemAV, magicItemSV, magicItemIsSold)
	totalGP := coinGP + specialTreasureGP + magicItemGP
	fullShare := math.RoundToEven(totalGP / float64(totalShares))
    }
*/

func (a ACKSII) CalculateXPShareAmount(shares float64, copper, silver, electrum, gold, platinum, totalJewelryValue, totalGemValue, totalMagicItemApparentValue, totalMagicItemSoldValue,
	totalCombatXP int) (fullShare, halfShare float64) {
	copperGoldValue := a.copperToGold(copper)
	silverGoldValue := a.silverToGold(silver)
	electrumGoldValue := a.electrumToGold(electrum)
	platinumGoldValue := a.platinumToGold(platinum)
	totalXPAvailable := float64(copperGoldValue + silverGoldValue + electrumGoldValue + gold + platinumGoldValue + totalJewelryValue + totalGemValue + totalMagicItemApparentValue + totalMagicItemSoldValue + totalCombatXP)
	fullShare = math.RoundToEven(totalXPAvailable / shares)
	halfShare = math.RoundToEven(float64(fullShare) / 2.0)
	return fullShare, halfShare
}

func (a ACKSII) CalculateGoldShareAmount(shares float64, copper, silver, electrum, gold, platinum, totalJewelryValue, totalGemValue, totalMagicItemApparentValue, totalMagicItemSoldValue int) (fullShare, halfShare float64) {
	copperGoldValue := a.copperToGold(copper)
	silverGoldValue := a.silverToGold(silver)
	electrumGoldValue := a.electrumToGold(electrum)
	platinumGoldValue := a.platinumToGold(platinum)
	totalGoldAvailable := float64(copperGoldValue + silverGoldValue + electrumGoldValue + gold + platinumGoldValue + totalJewelryValue + totalGemValue + totalMagicItemApparentValue + totalMagicItemSoldValue)
	fullShare = math.RoundToEven(totalGoldAvailable / shares)
	halfShare = math.RoundToEven(float64(fullShare) / 2.0)
	return fullShare, halfShare
}


/* Weather */
type Season int

const (
	Winter Season = iota
	Spring
	Summer
	Fall
)

func (s Season) String() string {
	return [...]string{"Winter", "Spring", "Summer", "Fall"}[s]
}

type WeatherModifiers struct {
	TempDayMod   int
	TempNightMod int
	PrecipMod    int
	WindMod      int
}


func (a ACKSII) GetWeatherModifers(koppenCode string, season Season) (WeatherModifiers, bool) {
	seasons, ok := a.koppenModifiers[koppenCode]
	if !ok {
		return WeatherModifiers{}, false
	}
	return seasons[season], true
}

func (a ACKSII) GetWindResult(diceRoll int, prevailingWind string) string {
	if diceRoll >= 8 {
		return prevailingWind
	}
	return a.winds[diceRoll-1]
}

func (a ACKSII) GetWeatherRolls() []int {
	temp := Sum(RollMultipleDice(2, 6))
	percip := Sum(RollMultipleDice(2, 6))
	wind := Sum(RollMultipleDice(2, 6))
	winddir := RollDice(12)

	return []int{temp, percip, wind, winddir}

}


func (a ACKSII) DailyWeather(diceRolls, previousDiceRolls []int, koppenCode, prevailingWind string, season Season, simulateFront bool) ([]string, error) {

	// Since the rolls get the same modifiers, we can apply the simulation here
	if simulateFront {
		for i, roll := range diceRolls {
			if roll == 2 || roll == 12 || roll == previousDiceRolls[i] {
				continue
			}
			if roll > previousDiceRolls[i] {
				roll = roll - 1
			}
			if roll < previousDiceRolls[i] {
				roll = roll + 1
			}
			diceRolls[i] = roll
		}
	}
	mods, ok := a.GetWeatherModifers(koppenCode, season)
	if !ok {
		return []string{}, errors.New(fmt.Sprintf("Unknown koppen code \"%s\" or season \"%s\"", koppenCode, string(season)))
	}

		tempDayRoll := diceRolls[0] + mods.TempDayMod
		tempNightRoll := diceRolls[0] + mods.TempNightMod
		percipRoll := diceRolls[1] + mods.PrecipMod
		windRoll := diceRolls[2] + mods.WindMod
		windDirRoll := diceRolls[3]

		dayTemp, ok := a.GetTempature(tempDayRoll, mods.TempDayMod)
	if !ok {
		return []string{}, errors.New(fmt.Sprintf("Unaccounted for day temp roll %d", tempDayRoll))
	}
		nightTemp, ok := a.GetTempature(tempNightRoll, mods.TempNightMod)
	if !ok {
		return []string{}, errors.New(fmt.Sprintf("Unaccounted for night temp roll %d", tempNightRoll))
	}
		percip, ok := a.GetPercipation(percipRoll)
	if !ok {
		return []string{}, errors.New(fmt.Sprintf("Unaccounted for percip temp roll %d", tempDayRoll))
	}
		wind, ok := a.GetWind(windRoll)
	if !ok {
		return []string{}, errors.New(fmt.Sprintf("Unaccounted wind roll %d", tempNightRoll))
	}
	windDir := a.GetWindDirection(windDirRoll, prevailingWind)

	return []string{dayTemp, nightTemp, percip, wind, windDir}, nil

}

func (a ACKSII) GetTempature(roll, modifier int) (string, bool) {
	if modifier <= 0 {
		switch {
		case roll <= 0:
			return "Frigid", true
		case roll <= 4:
			return "Cold", true
		case roll <= 6:
			return "Very Chilly", true
		case roll <= 8:
			return "Chilly", true
		case roll <= 10:
			return "Brisk", true
		}
		return "Balmy", true
	}
	if modifier > 0 {
		switch {
		case roll <= 1:
			return "Very Chilly", true
		case roll <= 4 :
			return "Chilly", true
		case roll == 5:
			return "Brisk", true
		case roll <= 8:
			return "Balmy", true
		case roll <= 10:
			return "Warm", true
		case roll <= 12:
			return "Hot", true
		}
		return "Sweltering", true
	}
	return "", false
}

func (a ACKSII) GetPercipation(roll int) (string, bool) {
	switch  {
	case roll <= -2:
		return "Sunbaked", true
	case roll <=3:
		return "Clear", true
	case roll == 4:
		return "Partly Cloudy", true
	case roll == 5:
		return "Mostly Cloudy", true
	case roll == 6:
		return "Overcast", true
	case roll <= 9:
		return "Drizzly", true
	case roll >= 10:
		return "Rainy", true
	}
	return "", false
}

func (a ACKSII) GetWind(roll int) (string, bool) {
	switch {
	case roll <= 4:
		return "Still", true
	case roll <= 6:
		return "Gentle", true
	case roll <= 9:
		return "Moderate", true
	case roll <= 11:
		return "Strong", true
	case roll <= 13:
		return "Very Strong, Windy", true
	case roll >= 14:
		return "Gale, Stormy", true
	}
	return "", false
}

func (a ACKSII) GetWindDirection(roll int, prevailing string) string {
	if roll >= len(a.winds) {
		return prevailing
	}
	return a.winds[roll]
}

func (a ACKSII) ListWinds() []string {
	return a.winds
}

func (a ACKSII) ListKoppenCodes() []string {
	codes := make([]string, len(a.koppenModifiers))

	i := 0
	for c := range a.koppenModifiers {
		codes[i] = c
		i++
	}
	return codes
}

func (a ACKSII) ListModifiers(koppenCode string, season Season) map[string]int {
	if mods, ok := a.GetWeatherModifers(koppenCode, season); ok {
		return map[string]int{
			"Day Temperature": mods.TempDayMod,
			"Night Temperature": mods.TempNightMod,
			"Percipitation": mods.PrecipMod,
			"Wind": mods.WindMod,
		}
	}
	return map[string]int{}
}
