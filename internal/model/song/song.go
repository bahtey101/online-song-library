package song

import "strings"

type Song struct {
	ID          uint64   `json:"id"`
	Group       string   `json:"group"`
	Song        string   `json:"song"`
	ReleaseDate string   `json:"releaseDate"`
	Verses      []string `json:"text"`
	Link        string   `json:"link"`
}

func SplitIntoVerses(text string) []string {
	return strings.Split(text, "\n\n")
}
