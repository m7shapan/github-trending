package main

import (
	"log"
	"os"
)

func main() {

	items, err := githubTrending("go", "daily")
	if err != nil {
		log.Fatal(err)
	}

	creds := Credentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	client, err := authTwitter(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Fatal(err)
	}

	err = tweet(client, items)
	if err != nil {
		log.Println("Error tweeting")
		log.Fatal(err)
	}
}
