package song

import "strings"

type Song struct {
	ID          uint64
	Group       string
	Song        string
	ReleaseDate string
	Verses      []string
	Link        string
}

func SplitIntoVerses(text string) []string {
	return strings.Split(text, "\n\n")
}
