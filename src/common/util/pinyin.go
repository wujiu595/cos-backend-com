package util

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

func ConvertToPinyin(v string) (full, short string) {
	parts := pinyin.LazyConvert(v, nil)
	for _, part := range parts {
		short += string(part[0])
	}
	full = strings.Join(parts, " ")
	return
}
