package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type redditResponseJSON struct {
	Data struct {
		Children []struct {
			Data Item `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Item struct {
	ID          string `json:"id"`
	Author      string `json:"author"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Score       int    `json:"score"`
	NumComments int    `json:"num_comments"`
	Downs       int    `json:"downs"`
	Ups         int    `json:"ups"`
	Over18      bool   `json:"over_18"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
	Image       io.Reader
}

const (
	UserAgent = "BrutBot Golang Reddit Reader 1.0"
	Limit     = 100
)

func getItems(subreddit string) ([]Item, error) {
	url := fmt.Sprintf("http://reddit.com/r/%s/top.json?limit=%dx&t=week", subreddit, Limit)
	fmt.Printf("fetching %s\n", url)

	// Create a request and add the proper headers.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)

	// Handle the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	data := new(redditResponseJSON)

	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	items := make([]Item, len(data.Data.Children))
	for i, child := range data.Data.Children {
		items[i] = child.Data
	}

	return items, err
}

var itemsMap = make(map[string]chan Item)

func loadImages(items []Item, ch chan<- Item) {
	for {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(len(items))

		item := items[n]

		if item.Image == nil {
			resp, err := http.Get(item.URL)

			if err != nil {
				continue
			}

			if resp.StatusCode == http.StatusOK {
				item.Image = resp.Body
			}
		}

		fmt.Printf("Pushing [%q, Body = %q, Image = %t] into ch\n", item.Title, item.Body, item.Image != nil)
		ch <- item
	}
}

func GetRandImage(subreddit string) (Item, error) {
	ch, ok := itemsMap[subreddit]
	if !ok {
		var err error

		ch = make(chan Item)
		items, err := getItems(subreddit)

		if err != nil {
			return Item{}, err
		}

		go loadImages(items, ch)
		itemsMap[subreddit] = ch
	}

	item := <-ch

	return item, nil
}
