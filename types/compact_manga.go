package types

// CompactManga represents a manga in a compact view
type CompactManga struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
