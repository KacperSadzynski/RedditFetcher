package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	subreddits := []string{"golang", "docker", "kubernetes", "aws", "googlecloud"}
	filePrefix := "reddit_output_"

	for _, subreddit := range subreddits {
		wg.Add(1)
		go run("https://www.reddit.com/r/"+subreddit+".json", filePrefix+subreddit+".txt")
	}
	wg.Wait()
}

// run function creates object Fetcher, fetches data and writes it to output file.
func run(hostUrl, fileName string) {

	defer wg.Done()

	fetcher := NewFetcher(hostUrl, time.Second*3)

	err := fetcher.Fetch(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	err = fetcher.Save(file)
	if err != nil {
		log.Fatal(err)
	}
}
