package gamesystems

import (
	"math"
)

type ACKSII struct {
	goldValueLookupTable []float64
	henchmenShare float64
}

func NewACKSII() ACKSII {
	return ACKSII{
		goldValueLookupTable: []float64{0.01, 0.1, 0.5, 1, 5},
		henchmenShare: 0.5,
	}
}

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
