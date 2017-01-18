package mangahere

import "testing"
import "github.com/davecgh/go-spew/spew"

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
