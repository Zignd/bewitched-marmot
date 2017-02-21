package mangahere

import (
	"fmt"
	"net/http"
	"net/url"

	"time"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/zignd/bewitched-marmot/types"
)

const host = "www.mangahere.co"

var baseURL = fmt.Sprintf("http://%s", host)

// Search queries the site
func Search(query string) ([]*types.CompactManga, error) {
	searchURL := fmt.Sprintf("%s/search.php?name=%s", baseURL, url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Search(\"%s\") could not create request", searchURL)
	}

	req.Header.Set("referer", baseURL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Search(\"%s\") failed to perform query", query)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, errors.Wrapf(err, "Search(\"%s\") failed to parse HTML document", searchURL)
	}
	if wasThrottled(doc) == true {
		duration := time.Duration(10) * time.Second
		time.Sleep(duration)
		return Search(query)
	}

	result := []*types.CompactManga{}

	docResultItems := doc.Find("body > section > article > div > div.result_search").Children().Not(".directory_footer").FilterFunction(func(index int, selection *goquery.Selection) bool {
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

func wasThrottled(doc *goquery.Document) bool {
	was := false

	doc.Find("body > section > article > div > div.result_search > dl").Children().Each(func(index int, selection *goquery.Selection) {
		was = (selection.Text() == "Sorry you have just searched, please try 5 seconds later.")
	})

	return was
}

// GetManga retrieves a manga
func GetManga(mangaURL string) (*types.DetailedManga, error) {
	req, err := http.NewRequest("GET", mangaURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "GetDetailedManga(%s) could not create a request", mangaURL)
	}

	req.Header.Set("referer", baseURL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "GetDetailedManga(%s) failed to retrieve the manga page", mangaURL)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, errors.Wrapf(err, "GetDetailedManga(%s) failed to parse HTML document", mangaURL)
	}

	detailedManga := &types.DetailedManga{}
	detailedManga.URL = mangaURL

	// DetailedManga.Name
	name, exists := doc.
		Find("head > meta[property=\"og:title\"]").
		Attr("content")
	if exists == false {
		return nil, errors.Errorf("GetDetailedManga(%s) failed to parse the manga name", mangaURL)
	}
	detailedManga.Name = name

	// DetailedManga.Description
	description := doc.Find("#show").Text()
	if description == "" {
		return nil, errors.Errorf("GetDetailedManga(%s) failed to parse the manga description", mangaURL)
	}
	detailedManga.Description = description

	// DetailedManga.Chapters
	chapters := []*types.CompactChapter{}
	var hasErr bool
	ulsSelection := doc.Find("#main > article > div > div.manga_detail > div.detail_list > ul").Eq(0).Children()
	chaptersNumbers := ulsSelection.Length()

	ulsSelection.Each(func(index int, chapterSelection *goquery.Selection) {
		chapter := &types.CompactChapter{}
		chapter.Number = chaptersNumbers
		chaptersNumbers = chaptersNumbers - 1

		wholeText := chapterSelection.Find("span.left").Text()
		aText := chapterSelection.Find("span.left > a").Text()
		chapter.Name = strings.Trim(strings.Replace(wholeText, aText, "", 1), "\n ")

		url, exists := chapterSelection.Find("a").Attr("href")
		if exists == false {
			hasErr = true
			return
		}
		chapter.URL = url

		chapters = append(chapters, chapter)
	})
	if hasErr {
		return nil, errors.Errorf("GetDetailedManga(%s) could not retrieve a URL for a chapter", mangaURL)
	}
	detailedManga.Chapters = chapters

	return detailedManga, nil
}

// GetChapter retrieves a chapter
func GetChapter(chapterURL string) (*types.DetailedChapter, error) {
	detailedChapter := &types.DetailedChapter{}
	detailedChapter.URL = chapterURL
	detailedChapter.PagesURLs = []string{}

	req, err := http.NewRequest("GET", chapterURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "GetChapter(\"%s\") could not create request", chapterURL)
	}

	req.Header.Add("referer", baseURL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "GetChapter(\"%s\") failed to retrieve the page", chapterURL)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, errors.Wrapf(err, "GetChapter(\"%s\") failed to parse HTML document", chapterURL)
	}

	pageURL, exists := doc.Find("#image").Attr("src")
	if !exists {
		return nil, errors.Errorf("GetChapter(\"%s\") failed to find the image URL on the HTML document", chapterURL)
	}

	detailedChapter.PagesURLs = append(detailedChapter.PagesURLs, pageURL)

	docOptions := doc.Find("div.go_page:nth-child(3) > span:nth-child(3) > select:nth-child(2)").Children()
	if docOptions.Find("option").Length() > 0 {
		return nil, errors.Errorf("GetChapter(\"%s\") failed to find the select element with the pages links on the HTML document", chapterURL)
	}

	docOptionsValues := []string{}

	var hasErr bool
	docOptions.Not("[selected=\"selected\"]").Each(func(index int, docOption *goquery.Selection) {
		value, exists := docOption.Attr("value")
		if !exists {
			hasErr = true
			return
		}
		docOptionsValues = append(docOptionsValues, value)
	})
	if hasErr {
		return nil, errors.Errorf("GetChapter(\"%s\") failed to find the pages links on the HTML document", chapterURL)
	}

	for _, docOptionValue := range docOptionsValues {
		req, err := http.NewRequest("GET", docOptionValue, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChapter(\"%s\") could not create request", chapterURL)
		}

		req.Header.Add("referer", baseURL)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChapter(\"%s\") failed to retrieve the page", chapterURL)
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChapter(\"%s\") failed to parse HTML document", chapterURL)
		}

		pageURL, exists := doc.Find("#image").Attr("src")
		if !exists {
			return nil, errors.Errorf("GetChapter(\"%s\") failed to find the image URL on the HTML document", chapterURL)
		}

		detailedChapter.PagesURLs = append(detailedChapter.PagesURLs, pageURL)
	}

	return detailedChapter, nil
}
