package mangahere

import "testing"

func TestSearch(t *testing.T) {
	query := "Noragami"

	list, error := Search(query)
	if error != nil {
		t.Errorf("Search(\"%v\") failed: %v", query, error)
		return
	}

	if len(list) == 0 {
		t.Errorf("Search(\"%v\") returned nothing. At least one item was expected", query)
	}
}
