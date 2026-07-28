package gamesystems

import (
	"fmt"
	"math/rand"
)
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

func SeasonFromString(s string) (Season, error) {
    switch s {
    case "Winter":
        return Winter, nil
    case "Spring":
        return Spring, nil
    case "Summer":
        return Summer, nil
    case "Fall":
        return Fall, nil
    default:
        return 0, fmt.Errorf("unknown season: %q", s)
    }
}

type Roll struct {
	Name, DiceLabel, Id string
	DieFace, DieNumber, Result int
	resolved bool
}

func (r *Roll) Resolve() {
	if !r.resolved {
		r.Result = Sum(RollMultipleDice(r.DieNumber, r.DieFace))
		r.resolved = true
	}
}

type Modifier struct {
	Name string
	Value int
}

type GameSystem interface {
	CalculateNumberOfShares(playerCount, henchmenCount int) float64
	CalculateTotalGPFromCoinage(copper, silver, electrum, gold, platinum int) float64
	CalculateTotalXPFromCombat(combatants, xp []int) float64
	CalculateTotalGPFromSpecialTreasure(numberRetrieved []int, gpValue []float64) float64
	CalculateTotalGPFromMagicItems(apparentValue, sellValue []int, isSold []bool) float64
	CalculateGPSharesFromSpecialTreasure(numberRetrieved []int, gpValue []float64, shares float64) (int, int)
	CalculateGPSharesFromMagicItems(apparentValue, sellValue []int, isSold []bool,shares float64) (int, int)
	CalculateTotalXP(copper, silver, electrum, gold, platinum int, specialTreasureRetrieved []int, specialTreasureValue []float64, combatantsDefeated, combatantXPValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) int
	CalculateTotalGP(copper, silver, electrum, gold, platinum int, specialTreasureRetrieved []int, specialTreasureValue []float64, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) int
	CalculateXPShares(totalShares float64, copper, silver, electrum, gold, platinum int, specialTreasureRetrieved []int, specialTreasureValue []float64, combatantsDefeated, combatantXPValue, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) (fullXPShare, henchmenShare int)
	CalculateGPShares(totalShares float64, copper, silver, electrum, gold, platinum int, specialTreasureRetrieved []int, specialTreasureValue []float64, magicItemAV, magicItemSV[]int, magicItemIsSold []bool) (fullShare, henchmenShare int)
	CalculateDetailedCoinage(shares float64, copper, silver, electrum, gold, platinum int) []int
    DailyWeather(diceRolls, previousDiceRolls []Roll, koppenCode, prevailingWind string, season Season, simulateFront bool) ([]string, error)
	ListWinds() []string
	ListKoppenCodes() []string
	ListWeatherModifiers( koppenCode string, season Season) []Modifier
	ListWeatherRolls() []Roll
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

func ListSeasons() []string {
	return []string{"Winter", "Spring", "Summer", "Fall"}
}
