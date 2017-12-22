package main

import "github.com/rentziass/territory_wars_report/guild"

func addToonToComparison(comparison map[string]*ToonComparison, toon *guild.Toon, comparisonType ComparisonType) {
	var toonCount *ToonCount
	if comparison[toon.Name] == nil {
		tc := &ToonComparison{}
		tc.Own = &ToonCount{}
		tc.Opponent = &ToonCount{}

		comparison[toon.Name] = tc
	}

	if comparisonType == IsOwn {
		toonCount = comparison[toon.Name].Own
	} else {
		toonCount = comparison[toon.Name].Opponent
	}

	// Total count
	toonCount.Total++

	// Stars
	if toon.Stars == 7 {
		toonCount.SevenStars++
	}
	if toon.Stars == 6 {
		toonCount.SixStars++
	}

	// Gear level
	if toon.GearLevel == 10 {
		toonCount.Gear10++
	}
	if toon.GearLevel == 11 {
		toonCount.Gear11++
	}
	if toon.GearLevel == 12 {
		toonCount.Gear12++
	}
}

func addShipToComparison(comparison map[string]*ShipComparison, ship *guild.Ship, comparisonType ComparisonType) {
	var shipCount *ShipCount
	if comparison[ship.Name] == nil {
		tc := &ShipComparison{}
		tc.Own = &ShipCount{}
		tc.Opponent = &ShipCount{}

		comparison[ship.Name] = tc
	}

	if comparisonType == IsOwn {
		shipCount = comparison[ship.Name].Own
	} else {
		shipCount = comparison[ship.Name].Opponent
	}

	// Total count
	shipCount.Total++

	// Stars
	if ship.Stars == 7 {
		shipCount.SevenStars++
	}
	if ship.Stars == 6 {
		shipCount.SixStars++
	}
}
