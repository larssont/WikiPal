package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//DefaultLanguage for searching
var DefaultLanguage = "en"

//Languages all language codes available to use
var Languages = map[string]string{
	"en": "english",
	"sv": "swedish",
	"de": "german",
	"fr": "french",
	"es": "spanish",
	"ru": "russian",
	"ja": "japanese",
	"nl": "dutch",
	"it": "italian",
	"pl": "polish",
	"vi": "vietnamese",
	"pt": "portuguese",
	"ar": "arabic",
	"zh": "chinese",
	"uk": "ukrainian",
	"ca": "catalan",
	"no": "norwegian",
	"fi": "finnish",
}

//Response struct
type Response struct {
	URL             string
	Totalhits       int
	AlternativeHits []string
}

//WikiQuery struct
type query struct {
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

func convertToWikiQuery(search string, langCode string) (q query) {

	url := queryPage(search, langCode)
	jsonData := getJSONData(url)

	json.Unmarshal([]byte(jsonData), &q)

	return
}

func getLanguage(a string) string {
	for k, v := range Languages {
		if a == k || a == v {
			return k
		}
	}
	return DefaultLanguage
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

//Search searches wikipedia for a given string
func Search(search string, langCode string) (res Response) {

	langCode = getLanguage(langCode)

	q := convertToWikiQuery(search, langCode)

	baseURL, wikiPath := parseWikipediaURL(langCode)

	baseURL.Path = wikiPath + q.Query.Search[0].Title

	res.URL = getFinalURL(baseURL.String())
	res.Totalhits = q.Query.Searchinfo.Totalhits

	for i := 1; i < len(q.Query.Search); i++ {
		baseURL.Path = wikiPath + q.Query.Search[i].Title
		res.AlternativeHits = append(res.AlternativeHits, getFinalURL(baseURL.String()))
	}

	return
}
