package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
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

	var folders = []string{
		"assets",
		"assets/factions",
		"assets/units",
		"wiki",
	}

	for _, folder := range folders {
		createFolder(folder)
	}
}

func main() {
	force := flag.Bool("force", false, "recreate all files regardless of whether they already exist")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(2)

	go fetchArmyData("n5", armyN4URL, *force, &wg)
	go wiki(&wg)

	wg.Wait()
	showFinalMessage()
}

func fetchArmyData(version string, endpoint string, force bool, wg *sync.WaitGroup) {
	defer wg.Done()

	c := httpClient()
	var armyData = sendRequest(c, endpoint)
	var armyObject Army
	json.Unmarshal(armyData, &armyObject)

	createFolder(version)
	createFile(version+"/army.json", []byte(prettyPrint(armyData)), true)

	// Process each faction concurrently. The global request semaphore in
	// sendRequest caps total in-flight requests, so spawning one goroutine per
	// faction (and per unit logo below) never exceeds the concurrency budget.
	var factionWg sync.WaitGroup
	for i := 0; i < len(armyObject.Factions); i++ {
		var factionID = armyObject.Factions[i].ID
		var factionSlug = armyObject.Factions[i].Slug
		var factionLogo = armyObject.Factions[i].Logo
		if factionSlug == "" {
			continue
		}
		factionWg.Add(1)
		go func(factionID int, factionSlug, factionLogo string) {
			defer factionWg.Done()
			processFaction(c, version, force, factionID, factionSlug, factionLogo)
		}(factionID, factionSlug, factionLogo)
	}
	factionWg.Wait()
}

func processFaction(c *http.Client, version string, force bool, factionID int, factionSlug, factionLogo string) {
	var factionLogoPath = "assets/factions/" + factionSlug + ".svg"

	if _, err := os.Stat(factionLogoPath); os.IsNotExist(err) {
		factionLogoData := sendRequest(c, factionLogo)
		if factionLogoData != nil {
			createFile(factionLogoPath, factionLogoData, false)
		}
	}

	if factionID == 901 { // skipping non-aligned armies
		return
	}

	var factionData = sendRequest(c, factionBaseEnURL+fmt.Sprintf("%d", factionID))
	if factionData == nil {
		return
	}

	var factionObject Faction
	json.Unmarshal(factionData, &factionObject)
	var factionFolderPath = version + "/" + factionSlug

	if factionObject.Version == "" || factionSlug == "" {
		return
	}

	var fileName = factionFolderPath + "/" + factionObject.Version + ".json"
	var fileNamePretty = factionFolderPath + "/" + factionSlug + ".json"

	createFolder(factionFolderPath)
	createFolder(factionFolderPath + "/units")

	if _, err := os.Stat(fileName); !force && !os.IsNotExist(err) {
		return
	}

	createFile(fileName, factionData, false)
	createFile(fileNamePretty, []byte(prettyPrint(factionData)), true)

	// Unit logos are independent downloads — fetch them concurrently too.
	var logoWg sync.WaitGroup
	for j := 0; j < len(factionObject.Resume); j++ {
		var unitLogoURL = factionObject.Resume[j].Logo
		var unitLogoURLArray = strings.Split(unitLogoURL, "/")
		var unitLogoFileName = unitLogoURLArray[len(unitLogoURLArray)-1]
		var unitLogoPath = "assets/units/" + unitLogoFileName

		if _, err := os.Stat(unitLogoPath); os.IsNotExist(err) {
			logoWg.Add(1)
			go func(unitLogoURL, unitLogoPath string) {
				defer logoWg.Done()
				var unitLogoData = sendRequest(c, unitLogoURL)
				if unitLogoData != nil {
					createFile(unitLogoPath, unitLogoData, false)
				}
			}(unitLogoURL, unitLogoPath)
		}
	}
	logoWg.Wait()

	for j := 0; j < len(factionObject.Units); j++ {
		var unitData, _ = json.Marshal(factionObject.Units[j])
		if unitData != nil {
			var unitSlug = factionObject.Units[j].Slug
			if unitSlug != "" {
				var fileNameUnit = factionFolderPath + "/units/" + unitSlug + ".json"

				createFile(fileNameUnit, []byte(prettyPrint(unitData)), true)
			}
		}
	}
}

type Army struct {
	Factions []struct {
		ID           int    `json:"id"`
		Parent       int    `json:"parent"`
		Name         string `json:"name"`
		Slug         string `json:"slug"`
		Logo         string `json:"logo"`
		Discontinued bool   `json:"discontinued"`
	} `json:"factions"`
	Ammunitions []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Wiki string `json:"wiki,omitempty"`
	} `json:"ammunitions"`
	Weapons []struct {
		ID         int         `json:"id"`
		Type       string      `json:"type"`
		Name       string      `json:"name"`
		Ammunition interface{} `json:"ammunition"`
		Burst      string      `json:"burst"`
		Damage     string      `json:"damage"`
		Saving     string      `json:"saving"`
		SavingNum  string      `json:"savingNum"`
		Properties []string    `json:"properties"`
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
