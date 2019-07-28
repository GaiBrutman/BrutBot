// reddit.go is responsible for fetching and loading Reddit posts
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

// A struct that represents a Reddit fetch response
type redditResponseJSON struct {
	Data struct {
		Children []struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// A struct that represents a Reddit post
type Post struct {
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
	UserAgent = "BrutBot Golang Reddit Reader 1.0" // The fetch UserAgent
	Limit     = 100                                // number of Reddit posts fetched
)

// Fetches [Limit] posts from a given subreddit
// Returns a slice of fetched Posts
func getPosts(subreddit string) ([]Post, error) {
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

	// Decode the JSON into the response struct
	data := new(redditResponseJSON)

	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	// Put the response in a slice of Posts
	posts := make([]Post, len(data.Data.Children))
	for i, child := range data.Data.Children {
		posts[i] = child.Data
	}

	return posts, err
}

// maps a subreddit to a Post channel of the subreddit's fetched posts
var postsMap = make(map[string]chan Post)

// Loads the image of a random Reddit post from posts and writes the post into ch
func loadImages(posts []Post, ch chan<- Post) {
	for {
		// Generate pseudo-random slice index
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(len(posts))

		post := posts[n]

		// fetch the posts image (if necessary)
		if post.Image == nil {
			resp, err := http.Get(post.URL)

			if err != nil {
				continue
			}

			if resp.StatusCode == http.StatusOK {
				post.Image = resp.Body
			}
		}

		// Push post into ch
		fmt.Printf("Pushing [%q, Body = %q, Image = %t] into ch\n", post.Title, post.Body, post.Image != nil)
		ch <- post
	}
}

func GetRandImage(subreddit string) (Post, error) {
	ch, ok := postsMap[subreddit]
	if !ok {
		var err error

		ch = make(chan Post)
		posts, err := getPosts(subreddit)

		if err != nil {
			return Post{}, err
		}

		go loadImages(posts, ch)
		postsMap[subreddit] = ch
	}

	post := <-ch

	return post, nil
}
