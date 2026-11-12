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

	var slugCount = make(map[string]int)
	for _, f := range armyObject.Factions {
		slugCount[f.Slug]++
	}

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
		var folderName = factionFolderName(factionSlug, factionID, slugCount)
		go func(factionID int, factionSlug, factionLogo, folderName string) {
			defer factionWg.Done()
			processFaction(c, version, force, factionID, factionSlug, factionLogo, folderName)
		}(factionID, factionSlug, factionLogo, folderName)
	}
	factionWg.Wait()
}

// factionFolderName is the faction's slug, or slug-ID when several factions
// share a slug (998 and 999 are both "contracted-back-up") so neither
// overwrites the other.
func factionFolderName(slug string, id int, slugCount map[string]int) string {
	if slugCount[slug] > 1 {
		return fmt.Sprintf("%s-%d", slug, id)
	}
	return slug
}

func processFaction(c *http.Client, version string, force bool, factionID int, factionSlug, factionLogo, folderName string) {
	var factionLogoPath = "assets/factions/" + factionSlug + ".svg"

	if _, err := os.Stat(factionLogoPath); os.IsNotExist(err) {
		factionLogoData := sendRequest(c, factionLogo)
		if factionLogoData != nil {
			createFile(factionLogoPath, factionLogoData, false)
		}
	}

	var factionData = sendRequest(c, factionBaseEnURL+fmt.Sprintf("%d", factionID))
	if factionData == nil {
		return
	}

	var factionObject Faction
	json.Unmarshal(factionData, &factionObject)
	var factionFolderPath = version + "/" + folderName

	if factionObject.Version == "" || factionSlug == "" {
		return
	}

	var fileName = factionFolderPath + "/" + factionObject.Version + ".json"
	var fileNamePretty = factionFolderPath + "/" + folderName + ".json"

	createFolder(factionFolderPath)
	createFolder(factionFolderPath + "/units")

	// Download unit logos before the up-to-date check so missing logos are
	// backfilled even when the faction version hasn't changed.
	downloadUnitLogos(c, unitLogoURLs(factionObject))

	if _, err := os.Stat(fileName); !force && !os.IsNotExist(err) {
		return
	}

	createFile(fileName, factionData, false)
	createFile(fileNamePretty, []byte(prettyPrint(factionData)), true)

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

// unitLogoURLs collects every unit logo a faction references: the resume
// (one per unit) plus each profile in every profile group, including spec-ops
// units. Profile logos cover variants such as team-ops HI (-2-1) and support
// troops (-3-1).
func unitLogoURLs(f Faction) []string {
	var seen = make(map[string]bool)
	var urls []string
	var add = func(u string) {
		if u != "" && !seen[u] {
			seen[u] = true
			urls = append(urls, u)
		}
	}
	for _, r := range f.Resume {
		add(r.Logo)
	}
	for _, u := range f.Units {
		for _, g := range u.ProfileGroups {
			for _, p := range g.Profiles {
				add(p.Logo)
			}
		}
	}
	for _, u := range f.Specops.Units {
		for _, g := range u.ProfileGroups {
			for _, p := range g.Profiles {
				add(p.Logo)
			}
		}
	}
	return urls
}

// downloadUnitLogos fetches each logo into assets/units/ if not already there.
func downloadUnitLogos(c *http.Client, urls []string) {
	var logoWg sync.WaitGroup
	for _, unitLogoURL := range urls {
		var unitLogoURLArray = strings.Split(unitLogoURL, "/")
		var unitLogoPath = "assets/units/" + unitLogoURLArray[len(unitLogoURLArray)-1]

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
}
