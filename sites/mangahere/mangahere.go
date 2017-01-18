package mangahere

import (
	"fmt"
	"net/http"
	"net/url"

	"time"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Zignd/bewitched-marmot/types"
)

const host = "www.mangahere.co"

var baseURL = fmt.Sprintf("http://%s", host)

// Search queries the site
func Search(query string) ([]*types.CompactManga, error) {
	searchURL := fmt.Sprintf("%s/search.php?name=%s", baseURL, url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Search(\"%s\") could not create request: %v", searchURL, err)
	}

	req.Header.Set("referer", baseURL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Search(\"%s\") failed to perform query: %v", query, err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, fmt.Errorf("Search(\"%s\") failed to parse HTML document: %v", searchURL, err)
	}
	if wasThrottled(doc) == true {
		duration := time.Duration(10) * time.Second
		time.Sleep(duration)
		return Search(query)
	}

	result := []*types.CompactManga{}

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
		item := &types.CompactManga{}
		item.Name = selection.Find("a.manga_info").Text()
		item.URL, _ = selection.Find("a.manga_info").Attr("href")
		result = append(result, item)
	})

	return result, nil
}

// GetDetailedManga returns detailed data for a manga based on its URL
func GetDetailedManga(mangaPageURL string) (*types.DetailedManga, error) {
	req, err := http.NewRequest("GET", mangaPageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("GetDetailedManga(%s) could not create a request: %v", mangaPageURL, err)
	}

	req.Header.Set("referer", baseURL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetDetailedManga(%s) failed to retrieve the manga page: %v", mangaPageURL, err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, fmt.Errorf("GetDetailedManga(%s) failed to parse HTML document: %v", mangaPageURL, err)
	}

	detailedManga := &types.DetailedManga{}
	detailedManga.URL = mangaPageURL

	// DetailedManga.Name
	name, exists := doc.
		Find("head > meta[property=\"og:title\"]").
		Attr("content")
	if exists == false {
		return nil, fmt.Errorf("GetDetailedManga(%s) failed to parse the manga name", mangaPageURL)
	}
	detailedManga.Name = name

	// DetailedManga.Description
	description := doc.Find("#show").Text()
	if description == "" {
		return nil, fmt.Errorf("GetDetailedManga(%s) failed to parse the manga description", mangaPageURL)
	}
	detailedManga.Description = description

	// DetailedManga.Chapters
	chapters := []*types.CompactChapter{}
	var errChapter error
	doc.
		Find("#main > article > div > div.manga_detail > div.detail_list > ul").
		Eq(0).
		Children().
		Each(func(index int, chapterSelection *goquery.Selection) {
			chapter := &types.CompactChapter{}

			wholeText := chapterSelection.Find("span.left").Text()
			aText := chapterSelection.Find("span.left > a").Text()
			chapter.Name = strings.Trim(strings.Replace(wholeText, aText, "", 1), "\n ")

			url, exists := chapterSelection.Find("a").Attr("href")
			if exists == false {
				errChapter = fmt.Errorf("GetDetailedManga(%s) could not retrieve a URL for a chapter", mangaPageURL)
			}
			chapter.URL = url

			chapters = append(chapters, chapter)
		})
	if errChapter != nil {
		return nil, errChapter
	}
	detailedManga.Chapters = chapters

	return detailedManga, nil
}

func wasThrottled(doc *goquery.Document) bool {
	was := false

	doc.Find("body > section > article > div > div.result_search > dl").Children().Each(func(index int, selection *goquery.Selection) {
		was = (selection.Text() == "Sorry you have just searched, please try 5 seconds later.")
	})

	return was
}
