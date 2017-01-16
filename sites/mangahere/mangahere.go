package mangahere

import (
	"fmt"
	"net/http"
	"net/url"

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

	result := []*types.SearchResultItem{}

	if doc.Find("body > section > article > div > div.result_search").Children().Not(".directory_footer").Length() == 0 {
		return result, nil
	}

	doc.Find("body > section > article > div > div.result_search").Children().Not(".directory_footer").Each(func(index int, selection *goquery.Selection) {
		item := &types.SearchResultItem{}
		item.Name = selection.Find("a.manga_info").Text()
		item.URL, _ = selection.Find("a.manga_info").Attr("href")
		result = append(result, item)
	})

	return result, nil
}
