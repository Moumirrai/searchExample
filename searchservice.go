package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Snippet  string `json:"snippet"`
	Position int    `json:"position"`
}

func Search(query string) []SearchResult {
	//covert the query to a URL

	rawQuery := url.QueryEscape(query)
	url := "https://www.google.com/search?q=" + rawQuery
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.3")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	rawResults := doc.Find("div.g")

	searchResults := make([]SearchResult, rawResults.Length())

	// iterate over results and populate the searchResults slice
	c := 0
	doc.Find("div.g").Each(func(i int, result *goquery.Selection) {
		title := result.Find("h3").First().Text()
		link, _ := result.Find("a").First().Attr("href")
		snippet := result.Find(".VwiC3b").First().Text()

		searchResults[c] = SearchResult{
			Title:    title,
			Link:     link,
			Snippet:  snippet,
			Position: c + 1,
		}

		c++
	})

	return searchResults

}
