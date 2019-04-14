package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Page Struct for JSON
type Page struct {
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

//Response Struct for JSON
type Response struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}

func getImageJSON(search string) string {
	search = strings.Title(search)

	req, err := http.NewRequest("GET", "http://en.wikipedia.org/w/api.php", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("action", "query")
	q.Add("prop", "pageimages")
	q.Add("pithumbsize", "500")
	q.Add("titles", search)
	q.Add("format", "json")

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func convertImageJSON(url string) Response {
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

	var res = Response{}

	err = json.Unmarshal([]byte(jsonDataFromHTTP), &res) // here!

	if err != nil {
		panic(err)
	}

	return res

}

func findImage(searchParameter string) (bool, string) {
	jsonURL := getImageJSON(searchParameter)
	res := convertImageJSON(jsonURL)

	var pages []Page

	for _, v := range res.Query.Pages {
		pages = append(pages, v)
	}

	source := pages[0].Thumbnail.Source
	found := false

	if source != "" {
		found = true
	}

	return found, source
}

func downloadImage(URL string) string {

	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	fileName := "../../tmp/pic.jpg"

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return fileName
}
