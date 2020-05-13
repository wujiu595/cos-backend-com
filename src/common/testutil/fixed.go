package testutil

import "regexp"

var (
	defaultPlaceHolderPrefix = `repl-placeholder-`
	DefaultFixedFilters      = []*FixedFilter{
		NewFixedFilter(
			`"(id|createdAt|updatedAt)": "[^"]*"(,?)`,
			func(ph string) string {
				return `"$1": "` + ph + `"$2`
			},
		),
	}
)

type FixedFilter struct {
	regex *regexp.Regexp
	repl  func(ph string) string
}

func NewFixedFilter(regex string, repl func(string) string) *FixedFilter {
	return &FixedFilter{regexp.MustCompile(regex), repl}
}
