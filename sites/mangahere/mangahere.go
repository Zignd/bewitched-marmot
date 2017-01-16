package mangahere

import (
	"fmt"
	"net/http"
	"net/url"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Zignd/bewitched-marmot/types"
)

const host = "www.mangahere.co"

var baseURL = fmt.Sprintf("http://%s/", host)

// Search queries Manga Here
func Search(query string) ([]*types.SearchResultItem, error) {
	searchURL := fmt.Sprintf("%s/search.php?name=%s", baseURL, url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request for %s: %v", host, err)
	}

	req.Header.Set("referer", baseURL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform query for %s: %v", query, err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response: %v", err)
	}
	if wasThrottled(doc) == true {
		duration := time.Duration(10) * time.Second
		time.Sleep(duration)
		return Search(query)
	}

	result := []*types.SearchResultItem{}

	docResultItems := doc.
		Find("body > section > article > div > div.result_search").
		Children().
		Not(".directory_footer").
		FilterFunction(func(index int, selection *goquery.Selection) bool {
			return selection.Text() != "No Manga Series"
		})

	if docResultItems.Length() == 0 {
		return result, nil
	}

	docResultItems.Each(func(index int, selection *goquery.Selection) {
		item := &types.SearchResultItem{}
		item.Name = selection.Find("a.manga_info").Text()
		item.URL, _ = selection.Find("a.manga_info").Attr("href")
		result = append(result, item)
	})

	return result, nil
}

func wasThrottled(doc *goquery.Document) bool {
	was := false

	doc.Find("body > section > article > div > div.result_search > dl").Children().Each(func(index int, selection *goquery.Selection) {
		was = (selection.Text() == "Sorry you have just searched, please try 5 seconds later.")
	})

	return was
}
