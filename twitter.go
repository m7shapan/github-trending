package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
)

// tweet max length is 280 characters.
const maxTweetLength = 280

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func authTwitter(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	_, _, err := client.Accounts.VerifyCredentials(verifyParams) // nolint: bodyclose
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify credentials")
	}

	return client, nil
}

// write a tweet.
func tweet(client *twitter.Client, items []Item) error {
	mainTweetBody := fmt.Sprintf("%s\n%s\n\n%s\n",
		"Today's #Golang Trending",
		time.Now().UTC().Format("2006-01-02"),
		"#go #golang #trending")

	tweet, resp, err := client.Statuses.Update(mainTweetBody, nil)
	if err != nil {
		log.Println(err)

		return errors.Wrap(err, "failed to tweet")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("%+v\n", resp)
	}

	for i, item := range items {
		// tweet max length is 280 characters
		// - (5 for newlines + 1 for hashtag and tweet number + repository name length + link length)
		item.Description = truncateText(item.Description, maxTweetLength-(5+4+len(item.Repo)+len(item.Link)))

		tweetBody := fmt.Sprintf("#%d\n\n#%s\n%s\n%s\n",
			i+1,
			item.Repo,
			item.Description, item.Link)

		tweet, resp, err = client.Statuses.Update(tweetBody, &twitter.StatusUpdateParams{
			InReplyToStatusID: tweet.ID,
		})
		defer resp.Body.Close()

		if err != nil {
			return errors.Wrap(err, "failed to tweet")
		}

		if resp.StatusCode != http.StatusOK {
			return errors.Wrap(ErrUnexpectedStatusCode, fmt.Sprintf("%+v", resp))
		}
	}

	return nil
}

func truncateText(text string, maxLength int) string {
	if len(text) > maxLength {
		return text[:maxLength-3] + "..."
	}

	return text
}
