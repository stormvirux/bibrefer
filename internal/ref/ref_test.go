package ref

import (
	"regexp"
	"strings"
	"testing"
)

func BenchmarkLoopingName(b *testing.B) {
	authors := []string{"Thaha Mohammed ", "Riaz Ahmed Sheikh ", "Timothy Anthony Wellington "}
	for i := 0; i < b.N; i++ {
		var newAuthor strings.Builder
		newAuthor.Grow(100)
		// r2 := regexp.MustCompile(`^\s*([\p{Lu}\\{}])[\p{L}-.'\\{}]+,?[\t\v ]+((?:[\p{L}-.'\\{}]*[\t \v]*)+)?$`)
		for i := 0; i < len(authors); i++ {
			author := strings.Split(strings.TrimSpace(authors[i]), " ")
			newAuthor.Reset()
			for j := 0; j < len(author)-1; j++ {
				newAuthor.WriteString(author[j][0:1])
				newAuthor.WriteString(". ")
			}
			authors[i] = newAuthor.String() + author[len(author)-1] // r2.ReplaceAllString(authors[i], "$1. $2")
		}
	}
}
func BenchmarkRegexName(b *testing.B) {
	authors := []string{"Thaha Mohammed ", "Riaz Ahmed Sheikh ", "Timothy Anthony Wellington "}
	for i := 0; i < b.N; i++ {
		r2 := regexp.MustCompile(`^\s*([\p{Lu}\\{}])[\p{L}-.'\\{}]+,?[\t\v ]+((?:[\p{L}-.'\\{}]*[\t \v]*)+)?$`)
		for j := 0; j < len(authors); j++ {
			authors[j] = r2.ReplaceAllString(authors[j], "$1. $2")
		}
	}
}
