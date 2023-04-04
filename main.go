package main

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	//Solution to 429 error https://www.reddit.com/r/redditdev/comments/t8e8hc/getting_nothing_but_429_responses_when_using_go/

	var f RedditFetcher
	var w io.Writer

	subreddits := []string{"golang", "docker", "kubernetes", "aws", "googlecloud"}
	for _, subreddit := range subreddits {
		wg.Add(1)
		go run(f, w, "https://www.reddit.com/r/"+subreddit+".json", "reddit_output_"+subreddit+".txt")
	}
	wg.Wait()
}

var wg sync.WaitGroup

func run(f RedditFetcher, w io.Writer, hostUrl, fileName string) {
	defer wg.Done()
	f = NewFetcher(hostUrl, time.Second*3)
	_, err := f.Fetch(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w = file
	err = f.Save(w)
	if err != nil {
		log.Fatal(err)
	}
}
