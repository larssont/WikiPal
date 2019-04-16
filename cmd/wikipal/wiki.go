package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

//WikiPage struct
type WikiPage struct {
	Title   string
	URL     string
	Image   Thumbnail
	Snippet string
}

//Thumbnail struct
type Thumbnail struct {
	ThumbnailFile     *os.File
	ThumbnailFileName string
}

// WikiPageSearch Struct for JSON
type WikiPageSearch struct {
	Query struct {
		Search []struct {
			Title  string `json:"title"`
			Pageid int    `json:"pageid"`
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
	Title     string `json:"title"`
	Thumbnail struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
}

//WikiPageExtract struct
type WikiPageExtract struct {
	Query struct {
		Pages map[string]ExtractPage `json:"pages"`
	} `json:"query"`
}

//ExtractPage struct
type ExtractPage struct {
	Pageid  int    `json:"pageid"`
	Title   string `json:"title"`
	Fullurl string `json:"fullurl"`
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
	q.Add("prop", "info")
	q.Add("srsearch", search)
	q.Add("srlimit", "3")
	q.Add("srprop", "")
	q.Add("utf8", "")
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func queryPageExtract(id int) string {
	req := getWiki()
	q := req.URL.Query()

	q.Add("action", "query")
	q.Add("prop", "info|extracts")
	q.Add("inprop", "url")
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
	absPath, _ := filepath.Abs("../../tmp/" + fileName)

	out, err := os.Create(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(absPath)
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

func test() {
	fmt.Println("test func")
}

//Searches wiki
func SearchWiki(input string) WikiPage {

	chExtract := make(chan string)
	chURL := make(chan string)
	chThumbnail := make(chan Thumbnail)

	pageSearchURL := queryPageSearch(input)
	pageSearch := convertToWikiPageSearch(pageSearchURL)

	pageid := pageSearch.Query.Search[0].Pageid
	title := pageSearch.Query.Search[0].Title

	if len(pageSearch.Query.Search) < 1 {
		return WikiPage{"", "", Thumbnail{nil, ""}, ""}
	}

	go func() {
		pageExtractURL := queryPageExtract(pageid)
		pageExtract := convertToWikiPageExtract(pageExtractURL)

		chExtract <- pageExtract.Query.Pages[strconv.Itoa(pageid)].Extract
		chURL <- pageExtract.Query.Pages[strconv.Itoa(pageid)].Fullurl

	}()
	go func() {
		pageThumbnailURL := queryPageThumbnail(pageid)
		pageThumbnail := convertToWikiPageImages(pageThumbnailURL)

		if pageThumbnail.Query.Pages[strconv.Itoa(pageid)].Thumbnail.Source != "" {
			file, fileName := downloadImage(pageThumbnail.Query.Pages[strconv.Itoa(pageid)].Thumbnail.Source)
			chThumbnail <- Thumbnail{file, fileName}
		}
	}()

	extract := <-chExtract
	URL := <-chURL
	image := <-chThumbnail

	return WikiPage{title, URL, image, extract}
}
