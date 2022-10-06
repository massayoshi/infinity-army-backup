package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	armyN4URL        string = getEnvVar("ARMY_N4_EN_URL")
	armyCodeoneURL   string = getEnvVar("ARMY_CODEONE_EN_URL")
	factionBaseEnURL string = getEnvVar("FACTION_N4_BASE_EN_URL")
)

func main() {
	// fetchArmyData("n4", armyN4URL, "")
	// fetchArmyData("codeone", armyCodeoneURL, "")
	wiki()
}

func fetchArmyData(version string, endpoint string, wiki string) {
	var armyData = getHTTPResponse(endpoint)
	var armyObject Army
	json.Unmarshal(armyData, &armyObject)

	createFolder(version)
	createFolder("wiki/skills")
	createFolder("wiki/equips")
	createFile(version+"/army.json", armyData, true)

	// if wiki != "" {
	// 	for i := 0; i < len(armyObject.Skills); i++ {
	// 		var skillURL = armyObject.Skills[i].Wiki
	// 		var skillURLArray = strings.Split(skillURL, "/")
	// 		var skillSlug = strings.Replace(skillURLArray[len(skillURLArray)-1], "?version=n4", "", -1)
	// 		fmt.Println(skillSlug)
	// 		if skillSlug != "" {
	// 			var skillData = getHTTPResponse(wikiBaseURL + skillSlug)
	// 			var skillObject WikiPage
	// 			json.Unmarshal(skillData, &skillObject)
	// 			var skillPageContent = skillObject.Query.Pages[0].Revisions[0].Slots.Main.Content
	// 			var skillCategory = skillObject.Query.Pages[0].Categories[0].Title
	// 			// make string url safe and lower case
	// 			skillCategory = strings.Replace(skillCategory, "Category:", "", -1)
	// 			skillCategory = strings.ReplaceAll(skillCategory, " ", "_")
	// 			skillCategory = strings.ToLower(skillCategory)

	// 			if strings.Contains(skillPageContent, "#redirect") {
	// 				var skillPageContentArray = strings.Split(skillPageContent, "#redirect [[")
	// 				var skillPageContentArray2 = strings.Split(skillPageContentArray[1], "]]")
	// 				skillSlug = skillPageContentArray2[0]
	// 				skillData = getHTTPResponse(wikiBaseURL + skillSlug)
	// 				json.Unmarshal(skillData, &skillObject)
	// 				skillPageContent = skillObject.Query.Pages[0].Revisions[0].Slots.Main.Content
	// 			}

	// 			createFile("wiki/skills/"+skillCategory+"-"+skillSlug+".json", []byte(skillPageContent), true)
	// 		}
	// 	}

	// for i := 0; i < len(armyObject.Equips); i++ {
	// 	var equipURL = armyObject.Equips[i].Wiki
	// 	var equipURLArray = strings.Split(equipURL, "/")
	// 	var equipSlug = strings.Replace(equipURLArray[len(equipURLArray)-1], "?version=n4", "", -1)
	// 	if equipSlug != "" {
	// 		var equipData = getHTTPResponse(wikiBaseURL + equipSlug)
	// 		var equipObject WikiPage
	// 		json.Unmarshal(equipData, &equipObject)
	// 		var equipPageContent = equipObject.Query.Pages[0].Revisions[0].Slots.Main.Content

	// 		if strings.Contains(equipPageContent, "#redirect") {
	// 			var equipPageContentArray = strings.Split(equipPageContent, "#redirect [[")
	// 			var equipPageContentArray2 = strings.Split(equipPageContentArray[1], "]]")
	// 			equipSlug = equipPageContentArray2[0]
	// 			equipData = getHTTPResponse(wikiBaseURL + equipSlug)
	// 			json.Unmarshal(equipData, &equipObject)
	// 			equipPageContent = equipObject.Query.Pages[0].Revisions[0].Slots.Main.Content
	// 		}

	// 		createFile("wiki/equips/"+equipSlug+".json", []byte(equipPageContent), true)
	// 	}
	// }
	// }

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

func getEnvVar(value string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	return os.Getenv(value)
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

type Faction struct {
	Version string `json:"version"`
	Units   []struct {
		ID            int         `json:"id"`
		IDArmy        int         `json:"idArmy"`
		Canonical     int         `json:"canonical"`
		Isc           string      `json:"isc"`
		IscAbbr       interface{} `json:"iscAbbr"`
		Name          string      `json:"name"`
		ProfileGroups []struct {
			ID       int         `json:"id"`
			Category int         `json:"category"`
			Isc      string      `json:"isc"`
			Notes    interface{} `json:"notes"`
			Profiles []struct {
				ID       int           `json:"id"`
				Arm      int           `json:"arm"`
				Ava      int           `json:"ava"`
				Bs       int           `json:"bs"`
				Bts      int           `json:"bts"`
				Cc       int           `json:"cc"`
				Chars    []int         `json:"chars"`
				Equip    []interface{} `json:"equip"`
				Logo     string        `json:"logo"`
				Weapons  []interface{} `json:"weapons"`
				Includes []interface{} `json:"includes"`
				Move     []int         `json:"move"`
				Ph       int           `json:"ph"`
				S        int           `json:"s"`
				Str      bool          `json:"str"`
				Type     int           `json:"type"`
				W        int           `json:"w"`
				Wip      int           `json:"wip"`
				Name     string        `json:"name"`
				Notes    interface{}   `json:"notes"`
				Skills   []struct {
					ID    int   `json:"id"`
					Order int   `json:"order"`
					Extra []int `json:"extra,omitempty"`
					Q     int   `json:"q,omitempty"`
				} `json:"skills"`
				Peripheral []interface{} `json:"peripheral"`
			} `json:"profiles"`
			Options []struct {
				ID       int           `json:"id"`
				Chars    []interface{} `json:"chars"`
				Disabled bool          `json:"disabled"`
				Equip    []interface{} `json:"equip"`
				Minis    int           `json:"minis"`
				Orders   []struct {
					Type  string `json:"type"`
					List  int    `json:"list"`
					Total int    `json:"total"`
				} `json:"orders"`
				Includes []interface{} `json:"includes"`
				Points   int           `json:"points"`
				Swc      string        `json:"swc"`
				Weapons  []struct {
					ID    int `json:"id"`
					Order int `json:"order"`
				} `json:"weapons"`
				Name       string        `json:"name"`
				Skills     []interface{} `json:"skills"`
				Peripheral []interface{} `json:"peripheral"`
			} `json:"options"`
		} `json:"profileGroups"`
		Options []interface{} `json:"options"`
		Slug    string        `json:"slug"`
		Filters struct {
			Categories []int         `json:"categories"`
			Skills     []int         `json:"skills"`
			Equip      []interface{} `json:"equip"`
			Chars      []int         `json:"chars"`
			Types      []int         `json:"types"`
			Weapons    []int         `json:"weapons"`
			Ammunition []int         `json:"ammunition"`
		} `json:"filters"`
		Factions []int  `json:"factions"`
		Notes    string `json:"notes,omitempty"`
	} `json:"units"`
	Filters struct {
		Peripheral []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Mercs bool   `json:"mercs,omitempty"`
		} `json:"peripheral"`
		Attrs []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Abbr string `json:"abbr"`
			Min  int    `json:"min"`
			Max  int    `json:"max"`
		} `json:"attrs"`
		Points []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"points"`
		Swc      []interface{} `json:"swc"`
		Category []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"category"`
		Ammunition []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Wiki  string `json:"wiki,omitempty"`
			Mercs bool   `json:"mercs,omitempty"`
		} `json:"ammunition"`
		Chars []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Wiki string `json:"wiki,omitempty"`
		} `json:"chars"`
		Type []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"type"`
		Equip []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
			Wiki  string `json:"wiki"`
			Mercs bool   `json:"mercs,omitempty"`
		} `json:"equip"`
		Skills []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Wiki    string `json:"wiki,omitempty"`
			Mercs   bool   `json:"mercs,omitempty"`
			Specops bool   `json:"specops,omitempty"`
		} `json:"skills"`
		Weapons []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
			Mercs bool   `json:"mercs,omitempty"`
			Wiki  string `json:"wiki,omitempty"`
		} `json:"weapons"`
		Extras []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
			Mercs bool   `json:"mercs,omitempty"`
		} `json:"extras"`
	} `json:"filters"`
	Resume []struct {
		ID       int    `json:"id"`
		Isc      string `json:"isc"`
		IDArmy   int    `json:"idArmy"`
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Logo     string `json:"logo"`
		Type     int    `json:"type"`
		Category int    `json:"category"`
	} `json:"resume"`
	Fireteams []interface{} `json:"fireteams"`
	Relations []struct {
		Units []struct {
			Unit    int `json:"unit"`
			Profile int `json:"profile"`
		} `json:"units"`
		Min   int  `json:"min"`
		Max   int  `json:"max"`
		Group bool `json:"group"`
	} `json:"relations"`
	Specops struct {
		Equip []struct {
			Exp   int `json:"exp"`
			ID    int `json:"id"`
			Equip []struct {
				ID int `json:"id"`
			} `json:"equip,omitempty"`
		} `json:"equip"`
		Skills []struct {
			Exp     int   `json:"exp"`
			ID      int   `json:"id"`
			Extras  []int `json:"extras,omitempty"`
			Weapons []struct {
				ID int `json:"id"`
			} `json:"weapons,omitempty"`
			Equip []struct {
				ID int `json:"id"`
			} `json:"equip,omitempty"`
		} `json:"skills"`
		Weapons []struct {
			Exp int `json:"exp"`
			ID  int `json:"id"`
		} `json:"weapons"`
		Units []struct {
			ID            int         `json:"id"`
			IDArmy        int         `json:"idArmy"`
			Canonical     int         `json:"canonical"`
			Isc           string      `json:"isc"`
			IscAbbr       interface{} `json:"iscAbbr"`
			Notes         interface{} `json:"notes"`
			Name          string      `json:"name"`
			ProfileGroups []struct {
				Notes    interface{} `json:"notes"`
				Isc      string      `json:"isc"`
				Profiles []struct {
					Bts      int           `json:"bts"`
					Cc       int           `json:"cc"`
					Move     []int         `json:"move"`
					Notes    interface{}   `json:"notes"`
					Includes []interface{} `json:"includes"`
					Type     int           `json:"type"`
					Ava      int           `json:"ava"`
					Str      bool          `json:"str"`
					Bs       int           `json:"bs"`
					S        int           `json:"s"`
					Equip    []interface{} `json:"equip"`
					W        int           `json:"w"`
					Ph       int           `json:"ph"`
					Name     string        `json:"name"`
					Logo     string        `json:"logo"`
					ID       int           `json:"id"`
					Arm      int           `json:"arm"`
					Weapons  []interface{} `json:"weapons"`
					Chars    []int         `json:"chars"`
					Wip      int           `json:"wip"`
					Skills   []struct {
						Q     int   `json:"q"`
						Extra []int `json:"extra"`
						ID    int   `json:"id"`
						Order int   `json:"order"`
					} `json:"skills"`
					Peripheral []interface{} `json:"peripheral"`
				} `json:"profiles"`
				Options []struct {
					Includes []interface{} `json:"includes"`
					Minis    int           `json:"minis"`
					Points   int           `json:"points"`
					Equip    []interface{} `json:"equip"`
					Name     string        `json:"name"`
					Disabled bool          `json:"disabled"`
					Orders   []struct {
						Type  string `json:"type"`
						List  int    `json:"list"`
						Total int    `json:"total"`
					} `json:"orders"`
					ID      int `json:"id"`
					Weapons []struct {
						ID    int `json:"id"`
						Order int `json:"order"`
					} `json:"weapons"`
					Chars  []interface{} `json:"chars"`
					Swc    string        `json:"swc"`
					Skills []struct {
						Q     int `json:"q"`
						ID    int `json:"id"`
						Order int `json:"order"`
					} `json:"skills"`
					Peripheral []interface{} `json:"peripheral"`
				} `json:"options"`
				ID       int `json:"id"`
				Category int `json:"category"`
			} `json:"profileGroups"`
			Options []interface{} `json:"options"`
			Slug    string        `json:"slug"`
			Filters struct {
				Categories []int         `json:"categories"`
				Skills     []int         `json:"skills"`
				Equip      []interface{} `json:"equip"`
				Chars      []int         `json:"chars"`
				Types      []int         `json:"types"`
				Weapons    []int         `json:"weapons"`
				Ammunition []int         `json:"ammunition"`
			} `json:"filters"`
			Factions []interface{} `json:"factions"`
		} `json:"units"`
	} `json:"specops"`
	FireteamChart struct {
		Spec struct {
			CORE  int `json:"CORE"`
			HARIS int `json:"HARIS"`
			DUO   int `json:"DUO"`
		} `json:"spec"`
		Desc  string `json:"desc"`
		Teams []struct {
			Name  string   `json:"name"`
			Obs   string   `json:"obs"`
			Type  []string `json:"type"`
			Units []struct {
				Min      int    `json:"min"`
				Max      int    `json:"max"`
				Name     string `json:"name"`
				Comment  string `json:"comment"`
				Required bool   `json:"required"`
				Slug     string `json:"slug"`
			} `json:"units"`
		} `json:"teams"`
	} `json:"fireteamChart"`
}

// type WikiPage struct {
// 	Batchcomplete bool `json:"batchcomplete"`
// 	Query         struct {
// 		Pages []struct {
// 			Pageid    int    `json:"pageid"`
// 			Ns        int    `json:"ns"`
// 			Title     string `json:"title"`
// 			Revisions []struct {
// 				Slots struct {
// 					Main struct {
// 						Contentmodel  string `json:"contentmodel"`
// 						Contentformat string `json:"contentformat"`
// 						Content       string `json:"content"`
// 					} `json:"main"`
// 				} `json:"slots"`
// 			} `json:"revisions"`
// 		} `json:"pages"`
// 	} `json:"query"`
// }

type WikiCategory struct {
	Batchcomplete bool `json:"batchcomplete"`
	Continue      struct {
		Accontinue string `json:"accontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Allcategories []struct {
			Category string `json:"category"`
			Size     int    `json:"size"`
			Pages    int    `json:"pages"`
			Files    int    `json:"files"`
			Subcats  int    `json:"subcats"`
		} `json:"allcategories"`
	} `json:"query"`
}

type WikiCategoryPages struct {
	Batchcomplete string `json:"batchcomplete"`
	Limits        struct {
		Categorymembers int `json:"categorymembers"`
	} `json:"limits"`
	Query struct {
		Pages []struct {
			Pageid int    `json:"pageid"`
			Ns     int    `json:"ns"`
			Title  string `json:"title"`
		} `json:"pages"`
	} `json:"query"`
}
