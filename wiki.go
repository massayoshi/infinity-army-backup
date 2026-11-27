package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	wikiPageURL     = getEnvVar("WIKI_PAGE_EN_URL")
	wikiPageListURL = getEnvVar("WIKI_ALL_EN_URL")
	urlList         []string

	// redirects maps a redirect page title to its target ("Page#Section"),
	// written to wiki/redirects.json so short names like "ARO" stay resolvable.
	redirects   = make(map[string]string)
	redirectsMu sync.Mutex

	// writtenPages holds the paths of the .wiki files written this run.
	writtenPages = make(map[string]bool)
	writtenMu    sync.Mutex

	redirectTarget = regexp.MustCompile(`\[\[\s*([^\]|]+?)\s*(?:\|[^\]]*)?\]\]`)
)

func wiki(wg *sync.WaitGroup) {
	defer wg.Done()
	c := httpClient()

	// Main-namespace pages hold the rules text, but much of what the live wiki
	// renders on them (e.g. {{modifiers-explained}} on Rolls) lives in templates,
	// so back up the Template namespace too or that content is lost.
	for _, ns := range []int{namespaceMain, namespaceTemplate} {
		var listURL = namespaceListURL(wikiPageListURL, ns)
		getPageList(c, listURL, listURL)
	}

	// Fetch and write each page concurrently. The global request semaphore in
	// sendRequest caps how many of these are actually in flight at once.
	var pageWg sync.WaitGroup
	for _, u := range urlList {
		pageWg.Add(1)
		go func(pageURL string) {
			defer pageWg.Done()
			fetchWikiPage(c, pageURL)
		}(u)
	}
	pageWg.Wait()

	removeStaleCopies()

	redirectsData, _ := json.Marshal(redirects)
	createFile("wiki/redirects.json", prettyPrint(redirectsData), true)
}

// parseRedirect returns the target of a "#REDIRECT [[Target]]" page.
func parseRedirect(content string) (string, bool) {
	var trimmed = strings.TrimSpace(content)
	if !strings.HasPrefix(strings.ToLower(trimmed), "#redirect") {
		return "", false
	}
	var match = redirectTarget.FindStringSubmatch(trimmed)
	if match == nil {
		return "", false
	}
	return match[1], true
}

const (
	namespaceMain     = 0
	namespaceTemplate = 10
	templatePrefix    = "Template:"
)

// namespaceListURL restricts an allpages list URL to one namespace. MediaWiki
// takes the last value of a repeated parameter, so this overrides any
// apnamespace already present in the configured URL.
func namespaceListURL(listURL string, ns int) string {
	return fmt.Sprintf("%s&apnamespace=%d", listURL, ns)
}

// releaseCategory matches errata/version tags such as "N5.2 Update",
// "N5 FAQ 0.1" or "N4 Annex". They are added to pages as rules change, so
// using them as the folder would move a page every time it is updated.
var releaseCategory = regexp.MustCompile(`(?i)^(Category:)?\s*N\d.*\b(Update|FAQ|Annex)\b`)

// pageFolder is the folder for a page's first topical category, falling back to
// its first release-tag category, or "Uncategorized".
func pageFolder(categories []string) string {
	if len(categories) == 0 {
		return "Uncategorized"
	}
	var pageCategory = categories[0]
	for _, category := range categories {
		if !releaseCategory.MatchString(category) {
			pageCategory = category
			break
		}
	}
	pageCategory = strings.Replace(pageCategory, "Category:", "", -1)
	pageCategory = strings.ReplaceAll(pageCategory, " ", "_")
	pageCategory = strings.ReplaceAll(pageCategory, "/", "")
	pageCategory = strings.ReplaceAll(pageCategory, "&", "and")
	return pageCategory
}

func fetchWikiPage(c *http.Client, pageURL string) {
	var pageData = sendRequest(c, pageURL)
	var pageObject WikiPage
	json.Unmarshal(pageData, &pageObject)

	// A 404 (nil body) or an otherwise empty/unexpected response leaves Pages
	// empty; skip rather than panic on Pages[0].
	if len(pageObject.Query.Pages) == 0 {
		return
	}

	var page = pageObject.Query.Pages[0]

	if page.Missing || len(page.Revisions) == 0 {
		return
	}

	var pageName = strings.TrimPrefix(page.Title, templatePrefix)
	pageName = strings.ReplaceAll(pageName, "/", "")
	pageName = strings.ReplaceAll(pageName, " ", "_")

	var revisions = page.Revisions
	if len(revisions) > 1 {
		fmt.Println("More than one revision for " + pageName)
	}
	var pageContent = page.Revisions[0].Slots.Main.Content

	if target, ok := parseRedirect(pageContent); ok {
		redirectsMu.Lock()
		redirects[page.Title] = target
		redirectsMu.Unlock()
		return
	}

	var categories []string
	for _, category := range page.Categories {
		categories = append(categories, category.Title)
	}
	var folderPath = "wiki/" + pageFolder(categories)
	if page.Ns == namespaceTemplate {
		folderPath = "wiki/Templates"
	}
	createFolder(folderPath)
	var filePath = folderPath + "/" + pageName + ".wiki"
	createFile(filePath, []byte(pageContent), true)

	writtenMu.Lock()
	writtenPages[filePath] = true
	writtenMu.Unlock()
}

// removeStaleCopies deletes older copies of pages written this run. A page's
// folder comes from its categories, so when those change the page is written
// to a new folder and the copy in the old one goes stale. Only files whose
// page was successfully re-fetched under another path are removed; pages that
// disappear from the wiki (or fail to fetch) are kept.
func removeStaleCopies() {
	var current = make(map[string]bool)
	for p := range writtenPages {
		current[staleKey(p)] = true
	}
	files, _ := filepath.Glob("wiki/*/*.wiki")
	for _, f := range files {
		if !writtenPages[f] && current[staleKey(f)] {
			os.Remove(f)
		}
	}
}

// staleKey identifies a page independent of its category folder. Templates
// live in their own folder and namespace, so they never match a regular page.
func staleKey(path string) string {
	var name = filepath.Base(path)
	if filepath.Base(filepath.Dir(path)) == "Templates" {
		return "Templates/" + name
	}
	return name
}

func getPageList(client *http.Client, listURL string, pageURL string) {
	var data = sendRequest(client, pageURL)
	parsePageList(client, listURL, data)
}

func parsePageList(client *http.Client, listURL string, data []byte) {
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
		var continueURL = listURL + "&apcontinue=" + url.QueryEscape(apContinue)
		getPageList(client, listURL, continueURL)
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
