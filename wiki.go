package main

import (
	"encoding/json"
	"net/url"
	"strings"
)

var (
	wikiPageURL     string = getEnvVar("WIKI_EN_URL")
	wikiPageListURL string = "https://infinitythewiki.com/api.php?action=query&format=json&prop=&list=allpages&indexpageids=1&continue=-||&redirects=1&converttitles=1&apnamespace=0&aplimit=max"
	urlList         []string
)

func wiki() {
	// var wikiData = getPageList(wikiPageListURL)
	// var wikiObject WikiPageList
	// json.Unmarshal(wikiData, &wikiObject)
	getPageList(wikiPageListURL)
	for _, url := range urlList {
		var pageData = getHTTPResponse(url)
		var pageObject WikiPage
		json.Unmarshal(pageData, &pageObject)
		var pageName = pageObject.Query.Pages[0].Title
		var pageContent = pageObject.Query.Pages[0].Revisions[0].Slots.Main.Content
		var pageCategory = pageObject.Query.Pages[0].Categories[0].Title
		pageCategory = strings.Replace(pageCategory, "Category:", "", -1)
		pageCategory = strings.ReplaceAll(pageCategory, " ", "_")
		pageCategory = strings.ToLower(pageCategory)

		createFolder("wiki/" + pageCategory)
		createFile(pageName+".wiki", []byte(pageContent), true)
	}
}

// getPageList gets a list of all pages from the wiki recursively
func getPageList(url string) []byte {
	var data = getHTTPResponse(url)
	var pageList = parsePageList(data)
	return pageList
}

// parsePageList parses the JSON response from the wiki API
func parsePageList(data []byte) []byte {
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
		return getPageList(continueURL)
	}

	return data
}

type WikiPage struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Pages []struct {
			Pageid     int    `json:"pageid"`
			Ns         int    `json:"ns"`
			Title      string `json:"title"`
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
