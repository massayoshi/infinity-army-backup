package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	armyN4URL        string = getEnvVar("ARMY_N4_EN_URL")
	armyCodeoneURL   string = getEnvVar("ARMY_CODEONE_EN_URL")
	factionBaseEnURL string = getEnvVar("FACTION_N4_BASE_EN_URL")
)

func main() {
	fetchArmyData("n4", armyN4URL)
	fetchArmyData("codeone", armyCodeoneURL)
	wiki()
}

func fetchArmyData(version string, endpoint string) {
	var armyData = getHTTPResponse(endpoint)
	var armyObject Army
	json.Unmarshal(armyData, &armyObject)

	createFolder(version)
	createFile(version+"/army.json", armyData, true)

	for i := 0; i < len(armyObject.Factions); i++ {
		var factionID = armyObject.Factions[i].ID
		var factionSlug string = armyObject.Factions[i].Slug
		var factionLogoURL = getHTTPResponse(armyObject.Factions[i].Logo)

		createFile("assets/factions/"+factionSlug+".svg", factionLogoURL, false)

		if factionID != 901 { // skipping non-aligned armies
			var factionData = getHTTPResponse(factionBaseEnURL + fmt.Sprintf("%d", factionID))
			var factionObject Faction
			json.Unmarshal(factionData, &factionObject)
			var factionFolderPath = version + "/" + factionSlug

			var fileName string = factionFolderPath + "/" + factionObject.Version + ".json"
			createFolder(factionFolderPath)
			createFile(fileName, factionData, false)

			for j := 0; j < len(factionObject.Resume); j++ {
				var unitLogoURL = factionObject.Resume[j].Logo
				var unitLogoURLArray = strings.Split(unitLogoURL, "/")
				var unitLogoFileName = unitLogoURLArray[len(unitLogoURLArray)-1]
				var unitLogoPath = "assets/units/" + unitLogoFileName

				if _, err := os.Stat(unitLogoPath); os.IsNotExist(err) {
					var unitLogoData = getHTTPResponse(unitLogoURL)
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
