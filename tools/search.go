package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SearchResponse struct {
	Items []struct {
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
		Pagemap struct {
			Metatags []struct {
				OgImage string `json:"og:image"`
			} `json:"metatags"`
			CseThumbnail []struct {
				Src string `json:"src"`
			} `json:"cse_thumbnail"`
		} `json:"pagemap"`
		DisplayLink string `json:"displayLink"`
	} `json:"items"`
}

type SearchResult struct {
	ImageLinks   []string
	Snippets     []string
	RelatedLinks []string
}

func Search(query string) (*SearchResult, error) {
	apiKey := ""
	cx := "f69337833a4504233"
	escapedTerm := url.QueryEscape(query)
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s", apiKey, cx, escapedTerm)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	result := &SearchResult{}
	relatedCount := 0

	ignoreDomains := []string{"ads", "spam"}
	ignore := func(domain string) bool {
		for _, d := range ignoreDomains {
			if strings.Contains(domain, d) {
				return true
			}
		}
		return false
	}

	for _, item := range searchResp.Items {
		// Image link extraction
		img := ""
		if len(item.Pagemap.Metatags) > 0 && item.Pagemap.Metatags[0].OgImage != "" {
			img = item.Pagemap.Metatags[0].OgImage
		} else if len(item.Pagemap.CseThumbnail) > 0 && item.Pagemap.CseThumbnail[0].Src != "" {
			img = item.Pagemap.CseThumbnail[0].Src
		}
		if img != "" && !ignore(img) {
			result.ImageLinks = append(result.ImageLinks, img)
		}

		// Snippet extraction
		if item.Snippet != "" {
			result.Snippets = append(result.Snippets, item.Snippet)
		}

		// Related links (max 2, not from ignoreDomains)
		if relatedCount < 2 && !ignore(item.DisplayLink) {
			result.RelatedLinks = append(result.RelatedLinks, item.Link)
			relatedCount++
		}
	}

	return result, nil
}
