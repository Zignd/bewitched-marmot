package types

// DetailedManga represents a manga in a detailed view
type DetailedManga struct {
	Name         string
	Description  string
	URL          string
	ChaptersURLs []string
}
