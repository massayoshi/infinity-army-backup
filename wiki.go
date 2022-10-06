package main

import (
	"encoding/json"
	"net/url"
	"strings"
)

var (
	wikiPageURL     string = getEnvVar("WIKI_PAGE_EN_URL")
	wikiPageListURL string = getEnvVar("WIKI_ALL_EN_URL")
	urlList         []string
)

func wiki() {
	getPageList(wikiPageListURL)
	for _, url := range urlList {
		var pageData = getHTTPResponse(url)
		var pageObject WikiPage
		json.Unmarshal(pageData, &pageObject)

		if pageObject.Query.Pages[0].Missing || len(pageObject.Query.Pages[0].Categories) == 0 {
			continue
		}

		var pageName = pageObject.Query.Pages[0].Title
		pageName = strings.ReplaceAll(pageName, "/", "")
		pageName = strings.ReplaceAll(pageName, " ", "_")

		var pageContent = pageObject.Query.Pages[0].Revisions[0].Slots.Main.Content

		if strings.Contains(pageContent, "#redirect") {
			continue
		}

		var pageCategory = pageObject.Query.Pages[0].Categories[0].Title
		pageCategory = strings.Replace(pageCategory, "Category:", "", -1)
		pageCategory = strings.ReplaceAll(pageCategory, " ", "_")
		pageCategory = strings.ReplaceAll(pageCategory, "/", "")
		pageCategory = strings.ReplaceAll(pageCategory, "&", "and")
		var folderPath = "wiki/" + pageCategory
		createFolder(folderPath)
		createFile(folderPath+"/"+pageName+".wiki", []byte(pageContent), true)
	}
}

func getPageList(url string) {
	var data = getHTTPResponse(url)
	parsePageList(data)
}

func parsePageList(data []byte) {
	var pageList = make(map[string]interface{})
	json.Unmarshal(data, &pageList)
	var pages = pageList["query"].(map[string]interface{})["allpages"].([]interface{})
	for _, page := range pages {
		var title = page.(map[string]interface{})["title"].(string)
		var pageURL = wikiPageURL + url.PathEscape(title)
		urlList = append(urlList, pageURL)
	}

	if pageList["continue"] != nil {
		var apContinue = pageList["continue"].(map[string]interface{})["apcontinue"].(string)
		var continueURL = wikiPageListURL + "&apcontinue=" + apContinue
		getPageList(continueURL)
	}
}

type WikiPage struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Pages []struct {
			Pageid     int    `json:"pageid"`
			Ns         int    `json:"ns"`
			Title      string `json:"title"`
			Missing    bool   `json:"missing"`
			Categories []struct {
				Ns    int    `json:"ns"`
				Title string `json:"title"`
			} `json:"categories"`
			Revisions []struct {
				Revid    int  `json:"revid"`
				Parentid int  `json:"parentid"`
				Minor    bool `json:"minor"`
				Slots    struct {
					Main struct {
						Contentmodel  string `json:"contentmodel"`
						Contentformat string `json:"contentformat"`
						Content       string `json:"content"`
					} `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

type WikiPageList struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Apcontinue string `json:"apcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Limits struct {
		Allpages int `json:"allpages"`
	} `json:"limits"`
	Query struct {
		Allpages []struct {
			Pageid int    `json:"pageid"`
			Ns     int    `json:"ns"`
			Title  string `json:"title"`
		} `json:"allpages"`
	} `json:"query"`
}
