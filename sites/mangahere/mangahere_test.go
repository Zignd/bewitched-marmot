package mangahere

import "testing"

func Test_Search_ShouldReturnData_WhenLookingForKnownMangaTitle(t *testing.T) {
	query := "Noragami"

	list, err := Search(query)
	if err != nil {
		t.Errorf("Search(\"%v\") failed: %v", query, err)
		return
	}

	if len(list) == 0 {
		t.Errorf("Search(\"%v\") returned nothing. At least one item was expected", query)
		return
	}
}

func Test_Search_ShouldNotReturnData_WhenLookingForNonexistentMangaTitle(t *testing.T) {
	query := "Qwerty!@#$1234"

	list, err := Search(query)
	if err != nil {
		t.Errorf("Search(\"%v\") failed: %v", query, err)
		return
	}

	if len(list) > 0 {
		t.Errorf("Search(\"%v\") returned something when it shouldn't have", query)
		return
	}
}
