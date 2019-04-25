package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var defaultLangCode = "en"

var langCodes = []string{
	"en",
	"sv",
	"de",
	"fr",
	"es",
	"ru",
	"ja",
	"nl",
	"it",
	"pl",
	"vi",
	"pt",
	"ar",
	"zh",
	"uk",
	"ca",
	"no",
	"fi",
}

//WikiResponse struct
type WikiResponse struct {
	URL             string
	Totalhits       int
	AlternativeHits []string
}

//WikiQuery struct
type WikiQuery struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Sroffset int    `json:"sroffset"`
		Continue string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Searchinfo struct {
			Totalhits int `json:"totalhits"`
		} `json:"searchinfo"`
		Search []struct {
			Ns        int    `json:"ns"`
			Title     string `json:"title"`
			Pageid    int    `json:"pageid"`
			Wordcount int    `json:"wordcount"`
		} `json:"search"`
	} `json:"query"`
}

func getWiki(langCode string) *http.Request {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s.wikipedia.org/w/api.php", langCode), nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	return req
}

func queryPage(search string, langCode string) string {
	req := getWiki(langCode)
	q := req.URL.Query()

	q.Add("action", "query")
	q.Add("list", "search")
	q.Add("srsearch", search)
	q.Add("srinfo", "totalhits")
	q.Add("srlimit", "3")
	q.Add("srprop", "wordcount")
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()

	return req.URL.String()
}

func getJSONData(url string) (jsonDataFromHTTP []byte) {

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	jsonDataFromHTTP, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return

}

func convertToWikiQuery(search string, langCode string) (query WikiQuery) {

	url := queryPage(search, langCode)
	jsonData := getJSONData(url)

	json.Unmarshal([]byte(jsonData), &query)

	return
}

func containsString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func parseWikipediaURL(langCode string) (baseURL *url.URL, path string) {

	baseURL, err := url.Parse(fmt.Sprintf("https://%s.wikipedia.org/wiki/", langCode))
	path = baseURL.Path
	if err != nil {
		panic(err)
	}

	return
}

func getFinalURL(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("http.Get => %v", err.Error())
	}
	return resp.Request.URL.String()
}

func searchWiki(search string, langCode string) (response WikiResponse) {

	if !containsString(langCode, langCodes) {
		langCode = defaultLangCode
	}

	q := convertToWikiQuery(search, langCode)

	baseURL, wikiPath := parseWikipediaURL(langCode)

	baseURL.Path = wikiPath + q.Query.Search[0].Title

	response.URL = getFinalURL(baseURL.String())
	response.Totalhits = q.Query.Searchinfo.Totalhits

	for i := 1; i < len(q.Query.Search); i++ {
		baseURL.Path = wikiPath + q.Query.Search[i].Title
		response.AlternativeHits = append(response.AlternativeHits, getFinalURL(baseURL.String()))
	}

	return
}
