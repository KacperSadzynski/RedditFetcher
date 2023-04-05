package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// A Response represents the structure of the data fetched from reddit
type Response struct {
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditFetcher interface {
	Fetch(context.Context) error
	Save(io.Writer) error
}

type Fetcher struct {
	c      *http.Client
	host   string
	output Response
}

// NewFetcher creates new instace of Fetcher object.
func NewFetcher(host string, t time.Duration) *Fetcher {
	return &Fetcher{
		host: host,
		c: &http.Client{
			Timeout: t,
			Transport: &http.Transport{
				TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
			},
		},
	}
}

type key int

const keyPrincipalID key = iota

// Fetch fetches the data from given subreddit host and returns Response struct with data.
func (e *Fetcher) Fetch(ctx context.Context) error {

	ctx = context.WithValue(ctx, keyPrincipalID, time.Now().Unix())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.host, http.NoBody)

	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("User-Agent", "Custom Agent")
	resp, err := e.c.Do(req)

	if err != nil {
		return fmt.Errorf("cannot get data: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return fmt.Errorf("cannot unmarshal data: %w", err)
	}

	e.output = data
	return nil
}

// Save writes the data to a file.
func (e *Fetcher) Save(w io.Writer) error {

	for _, child := range e.output.Data.Children {

		d := fmt.Sprintf("%s\n%s\n\n", child.Data.Title, child.Data.URL)
		_, err := w.Write([]byte(d))

		if err != nil {
			return err
		}
	}
	return nil
}
