package gamesystems

import (
	"fmt"
	"math/rand"
)


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

func RollDice(faces int) int {
	if faces < 2 {
		return 0
	}
	return rand.Intn(faces) + 1
}

func RollMultipleDice(count, faces int) []int {
	results := make([]int, count)
	for i := range results {
		results[i] = RollDice(faces)
	}
	return results
}

// Sum adds up a slice of roll results
func Sum(rolls []int) int {
	total := 0
	for _, r := range rolls {
		total += r
	}
	return total
}
