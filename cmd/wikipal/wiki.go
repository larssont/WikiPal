package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//WikiPage struct
type WikiPage struct {
	ThumbnailFile     *os.File
	ThumbnailFileName string
	Snippet           string
}

// WikiPageSearch Struct for JSON
type WikiPageSearch struct {
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
			Ns        int       `json:"ns"`
			Title     string    `json:"title"`
			Pageid    int       `json:"pageid"`
			Size      int       `json:"size"`
			Wordcount int       `json:"wordcount"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

//WikiPageImages struct
type WikiPageImages struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]ImagePage `json:"pages"`
	} `json:"query"`
}

//ImagePage struct
type ImagePage struct {
	Pageid    int    `json:"pageid"`
	Ns        int    `json:"ns"`
	Title     string `json:"title"`
	Thumbnail struct {
		Source string `json:"source"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail"`
	Pageimage string `json:"pageimage"`
}

//WikiPageExtract struct
type WikiPageExtract struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]ExtractPage `json:"pages"`
	} `json:"query"`
}

//ExtractPage struct
type ExtractPage struct {
	Pageid  int    `json:"pageid"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

func getWiki() *http.Request {
	req, err := http.NewRequest("GET", "http://en.wikipedia.org/w/api.php", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	return req
}

func queryPageThumbnail(id int) string {

	req := getWiki()
	q := req.URL.Query()

	q.Add("action", "query")
	q.Add("prop", "pageimages")
	q.Add("pithumbsize", "500")
	q.Add("pilicense", "any")
	q.Add("pageids", strconv.Itoa(id))
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func queryPageSearch(search string) string {

	req := getWiki()
	q := req.URL.Query()

	q.Add("action", "query")
	q.Add("list", "search")
	q.Add("srsearch", search)
	q.Add("srlimit", "3")
	q.Add("utf8", "")
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func queryPageExtract(id int) string {
	req := getWiki()
	q := req.URL.Query()

	q.Add("action", "query")
	q.Add("prop", "extracts")
	q.Add("exintro", "true")
	q.Add("exchars", "350")
	q.Add("pageids", strconv.Itoa(id))
	q.Add("explaintext", "true")
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func getJSONData(url string) []byte {

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return jsonDataFromHTTP

}

func convertToWikiPageSearch(url string) WikiPageSearch {

	jsonData := getJSONData(url)

	var search = WikiPageSearch{}

	json.Unmarshal([]byte(jsonData), &search)
	return search

}

func convertToWikiPageExtract(url string) WikiPageExtract {

	jsonData := getJSONData(url)

	var pageExtract = WikiPageExtract{}

	json.Unmarshal([]byte(jsonData), &pageExtract)

	return pageExtract
}

func convertToWikiPageImages(url string) WikiPageImages {

	jsonData := getJSONData(url)

	var pageImages = WikiPageImages{}

	json.Unmarshal([]byte(jsonData), &pageImages)

	return pageImages
}

func downloadImage(URL string) (*os.File, string) {

	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	fileName := "pic.jpg"
	filePath := "../../" + fileName

	out, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return file, fileName

}

/*
* SearchWiki earches wikipedia for a given searchterm and returns
* the image for given page and a short snippet of text.
*
 */
func SearchWiki(input string) WikiPage {

	pageSearchURL := queryPageSearch(input)
	pageSearch := convertToWikiPageSearch(pageSearchURL)

	if len(pageSearch.Query.Search) < 1 {
		return WikiPage{nil, "", ""}
	}

	pageid := pageSearch.Query.Search[0].Pageid

	pageExtractURL := queryPageExtract(pageid)
	pageExtract := convertToWikiPageExtract(pageExtractURL)

	extract := pageExtract.Query.Pages[strconv.Itoa(pageid)].Extract

	pageThumbnailURL := queryPageThumbnail(pageid)
	pageThumbnail := convertToWikiPageImages(pageThumbnailURL)

	if pageThumbnail.Query.Pages[strconv.Itoa(pageid)].Thumbnail.Source == "" {
		return WikiPage{nil, "", extract}
	}

	file, fileName := downloadImage(pageThumbnail.Query.Pages[strconv.Itoa(pageid)].Thumbnail.Source)

	return WikiPage{file, fileName, extract}
}
