package main

import (
	"fmt"
	"github.com/rentziass/territory_wars_report/guild"
	"github.com/tealeg/xlsx"
)

type ComparisonType int
const (
	IsOwn ComparisonType = iota
	IsOpponent
)

type ZetaCount struct {
	Own int
	Opponent int
}

type ToonCount struct {
	Total int
	Gear10 int
	Gear11 int
	Gear12 int
	SixStars int
	SevenStars int
}

type ToonComparison struct {
	Own *ToonCount
	Opponent *ToonCount
}

type ShipCount struct {
	Total int
	SixStars int
	SevenStars int
}

type ShipComparison struct {
	Own *ShipCount
	Opponent *ShipCount
}

func main() {
	fmt.Println("Welcome to Territory Wars Report")
	fmt.Println("--------------------------------")
	fmt.Println("Please report any issue on https://github.com/rentziass/territory_wars_report")
	fmt.Println("Come say hi on https://discord.gg/YGUa2Fy")
	fmt.Println("")
	fmt.Printf("Enter your guild swgoh.gg's URL: ")

	var ownGuildURL, opponentGuildURL string
	fmt.Scan(&ownGuildURL)
	fmt.Printf("Enter your opponent swgoh.gg's URL: ")
	fmt.Scan(&opponentGuildURL)

	zetas := map[string]*ZetaCount{}
	toons := map[string]*ToonComparison{}
	ships := map[string]*ShipComparison{}

	// Own guild
	ownGuild, err := guild.GetGuild(ownGuildURL)
	if err != nil {
		fmt.Println("There was a problem getting your own guild")
		fmt.Println(err)
	}

	// Zetas
	for _, m := range ownGuild.Members {
		for _, z := range m.Zetas {
			if zetas[z] == nil {
				zetas[z] = &ZetaCount{Own: 1}
				continue
			}

			zetas[z].Own++
		}
	}

	// Toons
	total := len(ownGuild.Members)
	for i, m := range ownGuild.Members {
		err := m.GetToonRoster()
		if err != nil {
			panic(err)
		}
		err = m.GetShipRoster()
		if err != nil {
			panic(err)
		}

		for _, t := range m.Toons {
			addToonToComparison(toons, t, IsOwn)
		}

		for _, s := range m.Ships {
			addShipToComparison(ships, s, IsOwn)
		}
		fmt.Printf("\rProcessed %v/%v toon rosters for %v...", i+1, total, ownGuild.Name)
	}
	fmt.Print("\n")

	// Opponent guild
	opponentGuild, err := guild.GetGuild(opponentGuildURL)
	if err != nil {
		fmt.Println("There was a problem getting your own guild")
		fmt.Println(err)
	}

	// Zetas
	for _, m := range opponentGuild.Members {
		for _, z := range m.Zetas {
			if zetas[z] == nil {
				zetas[z] = &ZetaCount{Opponent: 1}
				continue
			}

			zetas[z].Opponent++
		}
	}

	// Toons
	total = len(opponentGuild.Members)
	for i, m := range opponentGuild.Members {
		err := m.GetToonRoster()
		if err != nil {
			panic(err)
		}

		for _, t := range m.Toons {
			addToonToComparison(toons, t, IsOpponent)
		}

		for _, s := range m.Ships {
			addShipToComparison(ships, s, IsOpponent)
		}
		fmt.Printf("\rProcessed %v/%v toon rosters for %v...", i+1, total, opponentGuild.Name)
	}
	fmt.Print("\n")

	writeToXLSX(zetas, toons, ships, ownGuild, opponentGuild)

	fmt.Println()
	fmt.Println("Done!")
}

func writeToXLSX(zetas map[string]*ZetaCount, toons map[string]*ToonComparison, ships map[string]*ShipComparison, ownGuild, opponentGuild *guild.Guild) error {
	excelFileName := "./" + ownGuild.Name + " vs " + opponentGuild.Name +  ".xlsx"
	xlFile := xlsx.NewFile()

	// Zetas
	zetaSheet, err := xlFile.AddSheet("Zetas")
	if err != nil {
		panic(err)
	}

	header := zetaSheet.AddRow()
	header.AddCell().SetString("Zeta name")
	header.AddCell().SetString(ownGuild.Name)
	header.AddCell().SetString(opponentGuild.Name)

	for z, c := range zetas {
		r := zetaSheet.AddRow()
		r.AddCell().SetString(z)
		r.AddCell().SetInt(c.Own)
		r.AddCell().SetInt(c.Opponent)
	}

	// Toons
	// Own
	ownToonsSheet, err := xlFile.AddSheet("Toons - " + ownGuild.Name)
	if err != nil {
		panic(err)
	}

	header = ownToonsSheet.AddRow()
	header.AddCell().SetString("Toon")
	header.AddCell().SetString("Total")
	header.AddCell().SetString("6 Stars")
	header.AddCell().SetString("7 Stars")
	header.AddCell().SetString("Gear 10")
	header.AddCell().SetString("Gear 11")
	header.AddCell().SetString("Gear 12")

	for t, c := range toons {
		r := ownToonsSheet.AddRow()
		r.AddCell().SetString(t)
		r.AddCell().SetInt(c.Own.Total)
		r.AddCell().SetInt(c.Own.SixStars)
		r.AddCell().SetInt(c.Own.SevenStars)
		r.AddCell().SetInt(c.Own.Gear10)
		r.AddCell().SetInt(c.Own.Gear11)
		r.AddCell().SetInt(c.Own.Gear12)
	}

	// Opponent
	opponentToonsSheet, err := xlFile.AddSheet("Toons - " + opponentGuild.Name)
	if err != nil {
		panic(err)
	}

	header = opponentToonsSheet.AddRow()
	header.AddCell().SetString("Toon")
	header.AddCell().SetString("Total")
	header.AddCell().SetString("6 Stars")
	header.AddCell().SetString("7 Stars")
	header.AddCell().SetString("Gear 10")
	header.AddCell().SetString("Gear 11")
	header.AddCell().SetString("Gear 12")

	for t, c := range toons {
		r := opponentToonsSheet.AddRow()
		r.AddCell().SetString(t)
		r.AddCell().SetInt(c.Opponent.Total)
		r.AddCell().SetInt(c.Opponent.SixStars)
		r.AddCell().SetInt(c.Opponent.SevenStars)
		r.AddCell().SetInt(c.Opponent.Gear10)
		r.AddCell().SetInt(c.Opponent.Gear11)
		r.AddCell().SetInt(c.Opponent.Gear12)
	}

	// Ships
	// Own
	ownShipsSheet, err := xlFile.AddSheet("Ships - " + ownGuild.Name)
	if err != nil {
		panic(err)
	}

	header = ownShipsSheet.AddRow()
	header.AddCell().SetString("Ship")
	header.AddCell().SetString("Total")
	header.AddCell().SetString("6 Stars")
	header.AddCell().SetString("7 Stars")

	for t, c := range ships {
		r := ownShipsSheet.AddRow()
		r.AddCell().SetString(t)
		r.AddCell().SetInt(c.Own.Total)
		r.AddCell().SetInt(c.Own.SixStars)
		r.AddCell().SetInt(c.Own.SevenStars)
	}

	// Opponent
	opponentShipsSheet, err := xlFile.AddSheet("Ships - " + opponentGuild.Name)
	if err != nil {
		panic(err)
	}

	header = opponentShipsSheet.AddRow()
	header.AddCell().SetString("Ship")
	header.AddCell().SetString("Total")
	header.AddCell().SetString("6 Stars")
	header.AddCell().SetString("7 Stars")

	for t, c := range ships {
		r := opponentShipsSheet.AddRow()
		r.AddCell().SetString(t)
		r.AddCell().SetInt(c.Opponent.Total)
		r.AddCell().SetInt(c.Opponent.SixStars)
		r.AddCell().SetInt(c.Opponent.SevenStars)
	}

	xlFile.Save(excelFileName)
	return nil
}

