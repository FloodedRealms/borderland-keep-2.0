package gamesystems

import (
	"html/template"
	"math"
)

var goldValueLookupTable = []float64{0.01, 0.1, 0.5, 1, 5}

func copperToGold(copperamount int) int {
	return int(math.Floor(float64(copperamount) * goldValueLookupTable[0]))
}
func silverToGold(silver int) int {
	return int(math.Floor(float64(silver) * goldValueLookupTable[1]))
}
func electrumToGold(electrum int) int {
	return int(math.Floor(float64(electrum) * goldValueLookupTable[2]))
}
func platinumToGold(platium int) int {
	return int(math.Floor(float64(platium) * goldValueLookupTable[4]))
}

func calculateNumberOfShares(players, henchmen int) float64 {
	return float64(players) + float64(henchmen) / 2.0
}

func calculateXPShareAmount(shares float64, copper, silver, electrum, gold, platinum, totalJewelryValue, totalGemValue, totalMagicItemApparentValue, totalMagicItemSoldValue,
	totalCombatXP int) (fullShare, halfShare float64) {
	copperGoldValue := copperToGold(copper)
	silverGoldValue := silverToGold(silver)
	electrumGoldValue := electrumToGold(electrum)
	platinumGoldValue := platinumToGold(platinum)
	totalXPAvailable := float64(copperGoldValue + silverGoldValue + electrumGoldValue + gold + platinumGoldValue + totalJewelryValue + totalGemValue + totalMagicItemApparentValue + totalMagicItemSoldValue + totalCombatXP)
	fullShare = math.RoundToEven(totalXPAvailable / shares)
	halfShare = math.RoundToEven(float64(fullShare) / 2.0)
	return fullShare, halfShare
}

func calculateGoldShareAmount(shares float64, copper, silver, electrum, gold, platinum, totalJewelryValue, totalGemValue, totalMagicItemApparentValue, totalMagicItemSoldValue int) (fullShare, halfShare float64) {
	copperGoldValue := copperToGold(copper)
	silverGoldValue := silverToGold(silver)
	electrumGoldValue := electrumToGold(electrum)
	platinumGoldValue := platinumToGold(platinum)
	totalGoldAvailable := float64(copperGoldValue + silverGoldValue + electrumGoldValue + gold + platinumGoldValue + totalJewelryValue + totalGemValue + totalMagicItemApparentValue + totalMagicItemSoldValue)
	fullShare = math.RoundToEven(totalGoldAvailable / shares)
	halfShare = math.RoundToEven(float64(fullShare) / 2.0)
	return fullShare, halfShare
}

type FormSection struct {
	RenderedForm template.HTML
}
type ACKSIIForm struct {
	FormSections []FormSection
}
func GetACKSIIInputForm() (*ACKSIIForm, error) {
	type coinSectionData struct {
		Copper int
		Silver int
		Electrum int
		Gold int
		Platinum int
		CalculationLink string
	}
	coinData := coinSectionData{
			Copper: 0,
			Silver: 0,
			Electrum: 0,
			Gold: 0,
			Platinum: 0,
		    CalculationLink: "/ACKSII/Coins",
		}

	renderedCoinSection, err := renderTemplateWithData("/game-systems/templates/ACKSII", "coinFormSection", coinData)
	if err != nil {
		return nil, err
	}
	CoinSection := FormSection{
		RenderedForm: template.HTML(renderedCoinSection),
	}

	Form := ACKSIIForm {
		FormSections: []FormSection{
			CoinSection,
		},
		}

	return &Form, nil
}
