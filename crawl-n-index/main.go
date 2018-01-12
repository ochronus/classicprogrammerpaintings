package main

import (
	"log"
	"os"
	"strings"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	algoliaCliAppId, envSet := os.LookupEnv("ALGOLIA_APP_ID")
	algoliaCliApiKey, envSet := os.LookupEnv("ALGOLIA_API_KEY")

	if !envSet {
		log.Fatal("Please set the ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables")
	}

	searchClient := algoliasearch.NewClient(algoliaCliAppId, algoliaCliApiKey)
	searchIndex := searchClient.InitIndex("classicprogrammerpaintings")

	c := colly.NewCollector(
		colly.AllowedDomains("classicprogrammerpaintings.com"),
	)

	c.OnHTML("article.photo", func(e *colly.HTMLElement) {
		id := e.Attr("data-post-id")

		imgElem := e.DOM.Find("div.photo-wrapper img").First()
		imageUrl, _ := imgElem.Attr("src")
		imageAlt, _ := imgElem.Attr("alt")
		items := strings.Split(imageAlt, "\n")
		description := strings.Trim(items[0], "“”")
		algoliaObject := make(algoliasearch.Object)
		algoliaObject["objectID"] = id
		algoliaObject["ImageUrl"] = imageUrl
		algoliaObject["Description"] = description

		searchIndex.AddObject(algoliaObject)

	})

	c.OnHTML("#pagination a.next[href]", func(e *colly.HTMLElement) {
		nextPageUrl := "http://classicprogrammerpaintings.com" + e.Attr("href")
		c.Visit(nextPageUrl)
	})
	log.Println("Crawling http://classicprogrammerpaintings.com/ and saving to Algolia... ")
	c.Visit("http://classicprogrammerpaintings.com/")
	log.Println("...done")

}
