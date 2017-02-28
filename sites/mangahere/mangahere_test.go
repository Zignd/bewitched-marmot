package mangahere

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Search_ShouldReturnData_WhenLookingForKnownMangaTitle(t *testing.T) {
	query := "Noragami"

	list, err := Search(query)
	if err != nil {
		t.Errorf("Search(\"%s\") failed: %v", query, err)
		return
	}

	if len(list) == 0 {
		t.Errorf("Search(\"%s\") = %s, expected at least one item", query, spew.Sprint(list))
		return
	}
}

func Test_Search_ShouldNotReturnData_WhenLookingForNonexistentMangaTitle(t *testing.T) {
	query := "Qwerty!@#$1234"

	list, err := Search(query)
	if err != nil {
		t.Errorf("Search(\"%s\") failed: %v", query, err)
		return
	}

	if len(list) > 0 {
		t.Errorf("Search(\"%s\") = %s, returned something when it shouldn't have", query, spew.Sprint(list))
		return
	}
}

func Test_GetManga_ShouldReturnDetailedData_WhenLookingForExistingURL(t *testing.T) {
	url := "http://www.mangahere.co/manga/noragami/"

	detailedManga, err := GetManga(url)

	if err != nil {
		t.Errorf("GetManga(\"%s\") failed: %v", url, err)
		return
	}

	if detailedManga.Title != "Noragami" {
		t.Errorf("GetManga(\"%s\") returned detailedManga.Title = \"%s\" but the expected was \"Noragami\"", detailedManga.Title, url)
		return
	}

	if !(len(detailedManga.Chapters) >= 1 && detailedManga.Chapters[len(detailedManga.Chapters)-1].Title == "People Who Wear Sportswear") {
		t.Errorf("GetManga(\"%s\") returned CompactChapter.Title = \"%s\" for the first chapter but the expected was \"People Who Wear Sportswear\"", detailedManga.Chapters[len(detailedManga.Chapters)-1].Title, url)
		return
	}
}

func Test_GetChapter_ShouldReturnDetailedData_WhenLookingForExistingURL(t *testing.T) {
	url := "http://www.mangahere.co/manga/noragami/c001/"

	detailedChapter, err := GetChapter(url)
	if err != nil {
		t.Errorf("GetChapter(\"%s\") failed: %v", url, err)
		return
	}

	if len(detailedChapter.PagesURLs) != 69 {
		t.Errorf("GetChapter(\"%s\") did not return the expected amount of pages", url)
		return
	}
}
