package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	armyN4URL        = getEnvVar("ARMY_N4_EN_URL")
	armyCodeoneURL   = getEnvVar("ARMY_CODEONE_EN_URL")
	factionBaseEnURL = getEnvVar("FACTION_N4_BASE_EN_URL")
)

func init() {
	createFolder("assets")
	createFolder("assets/factions")
	createFolder("assets/units")
	createFolder("wiki")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	fetchArmyData("n4", armyN4URL, &wg)
	fetchArmyData("codeone", armyCodeoneURL, &wg)
	wiki(&wg)

	wg.Wait()
	showFinalMessage()
}

func fetchArmyData(version string, endpoint string, wg *sync.WaitGroup) {
	defer wg.Done()

	c := httpClient()
	var armyData = sendRequest(c, endpoint)
	var armyObject Army
	json.Unmarshal(armyData, &armyObject)

	createFolder(version)
	createFile(version+"/army.json", armyData, true)

	for i := 0; i < len(armyObject.Factions); i++ {
		var factionID = armyObject.Factions[i].ID
		var factionSlug = armyObject.Factions[i].Slug
		var factionLogoPath = "assets/factions/" + factionSlug + ".svg"

		if _, err := os.Stat(factionLogoPath); os.IsNotExist(err) {
			factionLogoData := sendRequest(c, armyObject.Factions[i].Logo)
			createFile(factionLogoPath, factionLogoData, false)
		}

		if factionID != 901 { // skipping non-aligned armies
			var factionData = sendRequest(c, factionBaseEnURL+fmt.Sprintf("%d", factionID))
			var factionObject Faction
			json.Unmarshal(factionData, &factionObject)
			var factionFolderPath = version + "/" + factionSlug

			var fileName = factionFolderPath + "/" + factionObject.Version + ".json"
			createFolder(factionFolderPath)
			createFile(fileName, factionData, false)

			for j := 0; j < len(factionObject.Resume); j++ {
				var unitLogoURL = factionObject.Resume[j].Logo
				var unitLogoURLArray = strings.Split(unitLogoURL, "/")
				var unitLogoFileName = unitLogoURLArray[len(unitLogoURLArray)-1]
				var unitLogoPath = "assets/units/" + unitLogoFileName

				if _, err := os.Stat(unitLogoPath); os.IsNotExist(err) {
					var unitLogoData = sendRequest(c, unitLogoURL)
					createFile(unitLogoPath, unitLogoData, false)
				}
			}
		}
	}
}

type Army struct {
	Factions []struct {
		ID     int    `json:"id"`
		Parent int    `json:"parent"`
		Name   string `json:"name"`
		Slug   string `json:"slug"`
		Logo   string `json:"logo"`
	} `json:"factions"`
	Ammunitions []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Wiki string `json:"wiki,omitempty"`
	} `json:"ammunitions"`
	Weapons []struct {
		ID         int      `json:"id"`
		Type       string   `json:"type"`
		Name       string   `json:"name"`
		Ammunition int      `json:"ammunition"`
		Burst      string   `json:"burst"`
		Damage     string   `json:"damage"`
		Saving     string   `json:"saving"`
		Properties []string `json:"properties"`
		Distance   struct {
			Short struct {
				Max int    `json:"max"`
				Mod string `json:"mod"`
			} `json:"short"`
			Max struct {
				Max int    `json:"max"`
				Mod string `json:"mod"`
			} `json:"max"`
			Med struct {
				Max int    `json:"max"`
				Mod string `json:"mod"`
			} `json:"med"`
			Long struct {
				Max int    `json:"max"`
				Mod string `json:"mod"`
			} `json:"long"`
		} `json:"distance"`
		Mode    string `json:"mode,omitempty"`
		Profile string `json:"profile,omitempty"`
		Wiki    string `json:"wiki,omitempty"`
	} `json:"weapons"`
	Skills []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Wiki string `json:"wiki,omitempty"`
	} `json:"skills"`
	Equips []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Wiki string `json:"wiki"`
	} `json:"equips"`
	Hack []struct {
		Damage    string   `json:"damage"`
		Devices   []int    `json:"devices,omitempty"`
		Attack    string   `json:"attack"`
		Name      string   `json:"name"`
		Burst     string   `json:"burst"`
		Opponent  string   `json:"opponent"`
		Special   string   `json:"special"`
		SkillType []string `json:"skillType"`
		Extra     int      `json:"extra"`
		Target    []string `json:"target"`
	} `json:"hack"`
	MartialArts []struct {
		Opponent string `json:"opponent"`
		Damage   string `json:"damage"`
		Attack   string `json:"attack"`
		Name     string `json:"name"`
		Burst    string `json:"burst"`
	} `json:"martialArts"`
	Metachemistry []struct {
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"metachemistry"`
	Booty []struct {
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"booty"`
}
