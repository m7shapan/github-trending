package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

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

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// write a tweet
func tweet(client *twitter.Client, items []Item) error {
	mainTweetBody := fmt.Sprintf("%s\n%s\n\n%s\n",
		"Today's #Golang Trending",
		time.Now().UTC().Format("2006-01-02"),
		"#go #golang #trending")

	tweet, resp, err := client.Statuses.Update(mainTweetBody, nil)

	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != 200 {
		log.Println(fmt.Errorf("%+v", resp))
	}

	for i, item := range items {
		// tweet max length is 280 characters - (5 for newlines + 1 for hashtag and tweet number + repository name length + link length)
		item.Description = truncateText(item.Description, 280-(5+4+len(item.Repo)+len(item.Link)))

		tweetBody := fmt.Sprintf("#%d\n\n#%s\n%s\n%s\n",
			i+1,
			item.Repo,
			item.Description, item.Link)

		tweet, resp, err = client.Statuses.Update(tweetBody, &twitter.StatusUpdateParams{
			InReplyToStatusID: tweet.ID,
		})

		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("%+v", resp)
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
