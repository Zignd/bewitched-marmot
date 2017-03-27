package types

// DetailedManga represents a manga in a detailed view
type DetailedManga struct {
	CompactManga
	Chapters []*CompactChapter `json:"chapters"`
}
