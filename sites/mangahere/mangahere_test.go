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

func Test_GetDetailedManga_ShouldReturnDetailedData_WhenLookingForExistingURL(t *testing.T) {
	url := "http://www.mangahere.co/manga/noragami/"

	detailedManga, err := GetDetailedManga(url)

	if err != nil {
		t.Errorf("GetDetailedManga(\"%s\") failed: %v", url, err)
		return
	}

	if detailedManga.Name != "Noragami" {
		t.Errorf("GetDetailedManga(\"%s\"), expected value for Name was \"Noragami\"", url)
		return
	}

	if !(len(detailedManga.Chapters) >= 1 && detailedManga.Chapters[len(detailedManga.Chapters)-1].Name == "People Who Wear Sportswear") {
		t.Errorf("GetDetailedManga(\"%s\"), expected Name of the first Chapter was \"People Who Wear Sportswear\"", url)
		return
	}
}
