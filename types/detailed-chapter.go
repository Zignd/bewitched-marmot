package types

// DetailedChapter represents a chapter in a detailed view
type DetailedChapter struct {
	CompactChapter
	PagesURLs []string
}
