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
	Fetch(context.Context) (Response, error)
	Save(io.Writer) error
}

type Fetcher struct {
	c      *http.Client
	host   string
	output Response
}

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

func (e *Fetcher) Fetch(ctx context.Context) (Response, error) {
	ctx = context.WithValue(ctx, "requestID", time.Now().Unix())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.host, http.NoBody)
	if err != nil {
		return Response{}, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("User-Agent", "Custom Agent")
	resp, err := e.c.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("cannot get data: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Response{}, fmt.Errorf("cannot unmarshal data: %w", err)
	}
	e.output = data
	return data, nil
}
func (e *Fetcher) Save(w io.Writer) error {
	//CREATE A TXT FILE
	for _, child := range e.output.Data.Children {
		d := fmt.Sprintf("%s\n%s\n\n", child.Data.Title, child.Data.URL)
		_, err := w.Write([]byte(d))
		if err != nil {
			return err
		}
	}
	return nil
}
