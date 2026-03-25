package gamesystems

import "fmt"


type GameSystem interface {
	CalculateNumberOfShares(playerCount, henchmenCount int) float64
	CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum int) float64
	CalculateTotalXPFromCombat(combatants, xp []int) float64
	CalculateTotalGPFromSpecialTreasure(numberRetrieved, gpValue []int) float64
	CalculateTotalGPFromMagicItems(apparentValue, sellValue []int, isSold []bool) float64
	CalculateGPSharesFromSpecialTreasure(numberRetrieved, gpValue []int, shares float64) (int, int)
	CalculateGPSharesFromMagicItems(apparentValue, sellValue []int, isSold []bool,shares float64) (int, int)
	CalculateTotalXP(copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, combatantsDefeated, combatantXPValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) int
	CalculateTotalGP(copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) int
	CalculateXPShares(totalShares float64, copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, combatantsDefeated, combatantXPValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) (fullXPShare, henchmenShare int)
	CalculateGPShares(totalShares float64, copper, silver, electrum, gold, platinum int, specialTreasureRetrieved, specialTreasureValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) (fullShare, henchmenShare int)
	CalculateDetailedCoinage(shares float64, copper, silver, electrum, gold, platinum int) []int
}

type SystemNotFoundError struct {
	Msg string
}

func (e SystemNotFoundError) Error() string {
	return e.Msg
}

func LoadGameSystem(systemName string) (GameSystem, error) {
	if (systemName == "ACKS II") {
		return NewACKSII(), nil
	}
	return nil,  SystemNotFoundError{ Msg: fmt.Sprintf("System %s not found!", systemName)}
}
