package types

// CompactChapter represents a chapter in a compact view
type CompactChapter struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}
