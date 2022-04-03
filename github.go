package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	User        string
	Repo        string
	Link        string
	Description string
}

func githubTrending(language, timeRange string) ([]Item, error) {
	res, err := http.Get(fmt.Sprintf("https://github.com/trending/%s?since=%s", language, timeRange))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var items []Item

	// Find the review items
	doc.Find(".application-main main .Box-row").Each(func(i int, s *goquery.Selection) {
		title, found := s.Find("h1 a").Attr("href")
		description := s.Find("p").Text()

		if found {
			parts := strings.Split(strings.TrimSpace(title), "/")

			user := parts[1]
			repository := parts[2]
			link := fmt.Sprintf("https://github.com/%s/%s", user, repository)

			items = append(items, Item{
				User:        user,
				Repo:        repository,
				Link:        link,
				Description: strings.TrimSpace(description),
			})
		}
	})

	return items, nil
}
