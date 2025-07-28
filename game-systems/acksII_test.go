package gamesystems

import (
	"testing"
)

func TestGoldConversions(t *testing.T) {
	var tests = []struct{
		name string
		coin string
		input int
		want int
	}{
		{"90 Copper should be 0 Gold", "copper", 90, 0},
		{"100 Copper should 1 Gold", "copper", 100, 1},
		{"1010 Copper should be 10 Gold", "copper", 1010, 10},
		{"9 Silver should be 0 Gold", "silver", 9, 0},
		{"100 Silver should 10 Gold", "silver", 100, 10},
		{"110 Silver should be 11 Gold", "silver", 110, 11},
		{"1 Electrum should be 0 Gold", "electrum", 1, 0},
		{"2 Electrum should 1 Gold", "electrum", 2, 1},
		{"11 Electrum should be 5 Gold", "electrum", 11, 5},
		{"1 Platinum should be 5 Gold", "platinum", 1, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			ans := -1
			switch (tt.coin) {
			case "copper":
				ans = copperToGold(tt.input)
			case "silver":
				ans = silverToGold(tt.input)
			case "electrum":
				ans = electrumToGold(tt.input)
			case "platinum":
				ans =platinumToGold(tt.input)
			}
			if ans != tt.want {
				t.Errorf("Got: %d\tWanted: %d\n", ans, tt.want)
			}
		})
	}
}

func TestXPCalculation(t *testing.T) {
	var treasure = []int{ // Total 1000
		10000, // 100 Gold in Copper
		1000,   // 100 Gold in Silver
		200,   // 100 Gold in Electrum
		100,   // 100 Gold
		20,    // 100 Gold in Platinum
		100,   // Jewelery
		100,   // Gems
		100,   // Magic Items
		100,   // Magic Items Sold
		100,   // Combat XP
	}
	var tests = []struct {
		name string
		players int
		henchmen int
		expShareNumber float64
		expFullShare float64
		expHalfShare float64
	}{
		{"5 Players", 5, 0, 5.0,200, 100},
		{"3 Players, 2 Henchmen", 3, 2, 4.0, 250, 125},
		{"3 Players, 1 Henchmen", 3, 1, 3.5, 286, 143},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			numberShares := calculateNumberOfShares(tt.players, tt.henchmen)

			if (numberShares != tt.expShareNumber) {
				t.Errorf("Wrong number of Shares!\nGot:%f\tWanted:%f", numberShares, tt.expShareNumber)
				return //no need to check calculation.
			}
			fShare, hShare := calculateXPShareAmount(numberShares, treasure[0], treasure[1], treasure[2], treasure[3], treasure[4], treasure[5], treasure[6], treasure[7], treasure[8], treasure[9])
			if fShare != tt.expFullShare {
				t.Errorf("Wrong Full Share Got: %f\tWanted: %f\n", fShare, tt.expFullShare)
			}
			if hShare != tt.expHalfShare {
				t.Errorf("Wrong Half Share Got: %f\tWanted: %f\n", hShare, tt.expHalfShare)
			}
		})
	}

}

func TestGoldCalculation(t *testing.T) {
	var treasure = []int{ // Total 1000
		10000, // 100 Gold in Copper
		1000,   // 100 Gold in Silver
		200,   // 100 Gold in Electrum
		100,   // 100 Gold
		20,    // 100 Gold in Platinum
		100,   // Jewelery
		100,   // Gems
		100,   // Magic Items
		100,   // Magic Items Sold
	}
	var tests = []struct {
		name string
		players int
		henchmen int
		expShareNumber float64
		expFullShare float64
		expHalfShare float64
	}{
		{"5 Players", 5, 0, 5.0, 180, 90},
		{"3 Players, 2 Henchmen", 3, 2, 4.0, 225, 112},
		{"3 Players, 1 Henchmen", 3, 1, 3.5, 257, 128},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			numberShares := calculateNumberOfShares(tt.players, tt.henchmen)

			if (numberShares != tt.expShareNumber) {
				t.Errorf("Wrong number of Shares!\nGot:%f\tWanted:%f", numberShares, tt.expShareNumber)
				return //no need to check calculation.
			}
			fShare, hShare := calculateGoldShareAmount(numberShares, treasure[0], treasure[1], treasure[2], treasure[3], treasure[4], treasure[5], treasure[6], treasure[7], treasure[8])
			if fShare != tt.expFullShare {
				t.Errorf("Wrong Full Share Got: %f\tWanted: %f\n", fShare, tt.expFullShare)
			}
			if hShare != tt.expHalfShare {
				t.Errorf("Wrong Half Share Got: %f\tWanted: %f\n", hShare, tt.expHalfShare)
			}
		})
	}

}
