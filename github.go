package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
)

const itemsHTMLPath = ".application-main main article.Box-row"

var ErrUnexpectedStatusCode = errors.New("unexpected status code")

type Item struct {
	User        string
	Repo        string
	Link        string
	Description string
}

func githubTrending(language, timeRange string) ([]Item, error) {
	res, err := http.Get(fmt.Sprintf("https://github.com/trending/%s?since=%s", language, timeRange))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trending")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Wrap(ErrUnexpectedStatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse response body")
	}

	// Find the review items
	items := getItems(doc)

	return items, nil
}

func getItems(doc *goquery.Document) []Item {
	var items []Item

	doc.Find(itemsHTMLPath).Each(func(i int, s *goquery.Selection) {
		title, found := s.Find("h2 a").Attr("href")
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

	return items
}
