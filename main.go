package main

import (
	"fmt"
	"github.com/rentziass/territory_wars_report/guild"
	"github.com/tealeg/xlsx"
)

type VsCount struct {
	Own int
	Opponent int
}

func main() {
	fmt.Println("Welcome to Territory Wars Report")
	fmt.Println("--------------------------------")
	fmt.Println("Please report any issue on GitHub repo")
	fmt.Println("Come say hi on https://discord.gg/YGUa2Fy")
	fmt.Println("")
	fmt.Printf("Enter your guild swgoh.gg's URL: ")

	var ownGuildURL, opponentGuildURL string
	fmt.Scan(&ownGuildURL)
	fmt.Printf("Enter your opponent swgoh.gg's URL: ")
	fmt.Scan(&opponentGuildURL)

	zetas := map[string]*VsCount{}
	toons := map[string]*VsCount{}

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
				zetas[z] = &VsCount{Own: 1}
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

		for _, t := range m.Toons {
			if toons[t] == nil {
				toons[t] = &VsCount{Own:1}
				continue
			}
			toons[t].Own++
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
				zetas[z] = &VsCount{Opponent: 1}
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
			if toons[t] == nil {
				toons[t] = &VsCount{Opponent:1}
				continue
			}
			toons[t].Opponent++
		}
		fmt.Printf("\rProcessed %v/%v toon rosters for %v...", i+1, total, ownGuild.Name)
	}
	fmt.Print("\n")

	writeToXLSX(zetas, toons, ownGuild, opponentGuild)

	fmt.Println()
	fmt.Println("Done!")
}

func writeToXLSX(zetas, toons map[string]*VsCount, ownGuild, opponentGuild *guild.Guild) error {
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
	toonsSheet, err := xlFile.AddSheet("Toons")
	if err != nil {
		panic(err)
	}

	header = toonsSheet.AddRow()
	header.AddCell().SetString("Toon")
	header.AddCell().SetString(ownGuild.Name)
	header.AddCell().SetString(opponentGuild.Name)

	for t, c := range toons {
		r := toonsSheet.AddRow()
		r.AddCell().SetString(t)
		r.AddCell().SetInt(c.Own)
		r.AddCell().SetInt(c.Opponent)
	}

	xlFile.Save(excelFileName)
	return nil
}

