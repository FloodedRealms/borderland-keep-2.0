package gamesystems

import (
	"testing"
)


func Treasure() ([][]int, []bool) {
		var coins = []int{ // Total 1000
		10000, // 100 Gold in Copper
		1000,   // 100 Gold in Silver
		200,   // 100 Gold in Electrum
		100,   // 100 Gold
		20,    // 100 Gold in Platinum
	}
	var stn = []int{1, 2, 5}
	var stv = []int{50, 25, 20} // This represents 1 50 GP, 2 25 GP amd 5 20 GP special treasures ones
	var cbn = []int{1, 7, 3}
	var cbv = []int{15, 10, 5} // This represents 1 15 XP, 7 10 XP and 3 5 XP Combats
	var mgav = []int{50, 50}
	var mgsp = []int{1000, 100}
	var mgis = []bool{false, true}

	return [][]int{coins, stn, stv, cbn, cbv, mgav, mgsp}, mgis
	}

func TestGoldConversions(t *testing.T) {
	a := NewACKSII()
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
				ans = a.copperToGold(tt.input)
			case "silver":
				ans = a.silverToGold(tt.input)
			case "electrum":
				ans = a.electrumToGold(tt.input)
			case "platinum":
				ans =a.platinumToGold(tt.input)
			}
			if ans != tt.want {
				t.Errorf("Got: %d\tWanted: %d\n", ans, tt.want)
			}
		})
	}
}

func TestMagicItemGPCalulation(t *testing.T) {
	a := NewACKSII()

	treasure, mgis := Treasure()
	magicItemGPValue := a.CalculateTotalGPFromMagicItems(treasure[5], treasure[6], mgis)
	expectedMagicItemValue := 200.0

	if  magicItemGPValue != float64(expectedMagicItemValue) {
		t.Errorf("Wrong Magic Item Value Got: %f\tWanted: %f\n", magicItemGPValue, expectedMagicItemValue)
	}

}

func TestXPCalculation(t *testing.T) {
	a := NewACKSII()
	treasure, mgis := Treasure()
	var tests = []struct {
		name string
		players int
		henchmen int
		expShareNumber float64
		expFullShare int
		expHalfShare int
	}{
		{"5 Players", 5, 0, 5.0,200, 100},
		{"3 Players, 2 Henchmen", 3, 2, 4.0, 250, 125},
		{"3 Players, 1 Henchmen", 3, 1, 3.5, 286, 143},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			numberShares := a.CalculateNumberOfShares(tt.players, tt.henchmen)

			if (numberShares != tt.expShareNumber) {
				t.Errorf("Wrong number of Shares!\nGot:%f\tWanted:%f", numberShares, tt.expShareNumber)
				return //no need to check calculation.
			}
			fShare, hShare := a.CalculateXPShares(numberShares, treasure[0][0], treasure[0][1], treasure[0][2], treasure[0][3], treasure[0][4], treasure[1], treasure[2], treasure[3], treasure[4], treasure[5], treasure[6], mgis )
			if fShare != tt.expFullShare {
				t.Errorf("Wrong Full Share Got: %d\tWanted: %d\n", fShare, tt.expFullShare)
			}
			if hShare != tt.expHalfShare {
				t.Errorf("Wrong Half Share Got: %d\tWanted: %d\n", hShare, tt.expHalfShare)
			}
		})
	}
}

func TestGoldCalculation(t *testing.T) {
	a := NewACKSII()
	treasure, mgis := Treasure()
	var tests = []struct {
		name string
		players int
		henchmen int
		expShareNumber float64
		expFullShare int
		expHalfShare int
	}{
		{"5 Players", 5, 0, 5.0, 180, 90},
		{"3 Players, 2 Henchmen", 3, 2, 4.0, 225, 112},
		{"3 Players, 1 Henchmen", 3, 1, 3.5, 257, 128},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			numberShares := a.CalculateNumberOfShares(tt.players, tt.henchmen)

			if (numberShares != tt.expShareNumber) {
				t.Errorf("Wrong number of Shares!\nGot:%f\tWanted:%f", numberShares, tt.expShareNumber)
				return //no need to check calculation.
			}
			fShare, hShare := a.CalculateGPShares(numberShares, treasure[0][0], treasure[0][1], treasure[0][2], treasure[0][3], treasure[0][4], treasure[1], treasure[2], treasure[5], treasure[6], mgis )

			if fShare != tt.expFullShare {
				t.Errorf("Wrong Full Share Got: %d\tWanted: %d\n", fShare, tt.expFullShare)
			}
			if hShare != tt.expHalfShare {
				t.Errorf("Wrong Half Share Got: %d\tWanted: %d\n", hShare, tt.expHalfShare)
			}
		})
	}

}


func TestForLargeValueBug (t *testing.T) {
	a := NewACKSII()
	treasure, mgis := Treasure()
	treasure[0][4] = 1000 // makes total XP 5980
	var tests = []struct {
		name string
		players int
		henchmen int
		expShareNumber float64
		expFullShare int
		expHalfShare int
	}{
		{"5 Players", 5, 0, 5.0, 1160, 580},
		{"2 Players", 2, 0, 2.0, 2900, 1450},
		{"3 Players, 2 Henchmen", 3, 2, 4.0, 1450, 725},
		{"3 Players, 1 Henchmen", 3, 1, 3.5, 1657, 828},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			numberShares := a.CalculateNumberOfShares(tt.players, tt.henchmen)

			if (numberShares != tt.expShareNumber) {
				t.Errorf("Wrong number of Shares!\nGot:%f\tWanted:%f", numberShares, tt.expShareNumber)
				return //no need to check calculation.
			}
			fShare, hShare := a.CalculateGPShares(numberShares, treasure[0][0], treasure[0][1], treasure[0][2], treasure[0][3], treasure[0][4], treasure[1], treasure[2], treasure[5], treasure[6], mgis )

			if fShare != tt.expFullShare {
				t.Errorf("Wrong Full Share Got: %d\tWanted: %d\n", fShare, tt.expFullShare)
			}
			if hShare != tt.expHalfShare {
				t.Errorf("Wrong Half Share Got: %d\tWanted: %d\n", hShare, tt.expHalfShare)
			}
		})
	}

}

func TestWindDirections (t *testing.T) {
	a := NewACKSII()
	var tests = []struct {
		name string
		roll int
		expectedWind, prevailingWind string
	}{
		{"Test Northerly result", 1, "Northerly", "Southerly"},
		{"Test Northeasterly result", 2, "Northeasterly", "Southerly"},
		{"Test Easterly result", 3, "Easterly", "Southerly"},
		{"Test Southeasterly result", 4, "Southeasterly", "Southerly"},
		{"Test Southerly result", 5, "Southerly", "Northerly"},
		{"Test Southwesterly result", 6, "Southwesterly", "Southerly"},
		{"Test Westerly result", 7, "Westerly", "Southerly"},
		{"Test Northerwesterly result", 8, "Northwesterly", "Southerly"},
		{"Test Prevailing 1 result", 9, "Southerly", "Southerly"},
		{"Test Prevailing 2 result", 10, "Southerly", "Southerly"},
		{"Test Prevailing 3 result", 11, "Southerly", "Southerly"},
		{"Test Prevailing 4 result", 12, "Southerly", "Southerly"},


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			wind := a.GetWindDirection(tt.roll, tt.prevailingWind)
			if (wind != tt.expectedWind) {
				t.Errorf("Wrong wind for roll %d! Got:%s\tWanted:%s", tt.roll, wind, tt.expectedWind)
			}
		})
	}
}
