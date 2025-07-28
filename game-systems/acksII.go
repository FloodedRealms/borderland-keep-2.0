package gamesystems

import "math"

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

func GetInputForm() struct {
	return form
}
