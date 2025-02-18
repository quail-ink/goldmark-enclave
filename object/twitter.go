package object

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TwitterOembedResp struct {
	URL        string `json:"url"`
	AuthorName string `json:"author_name"`
	AuthorURL  string `json:"author_url"`
	HTML       string `json:"html"`
	Type       string `json:"type"`
}

func GetTweetOembedHtml(url, theme string) (string, error) {
	if theme == "dark" {
		theme = "dark"
	} else {
		theme = "light"
	}

	oembedUrl := fmt.Sprintf("https://publish.x.com/oembed?url=%s&theme=%s", url, theme)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, oembedUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Access-Control-Allow-Origin", "*")
	req.Header.Add("Access-Control-Allow-Methods", "GET")
	req.Header.Add("Access-Control-Allow-Headers", "Content-Type")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body TwitterOembedResp
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return body.HTML, nil
}
