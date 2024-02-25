package utils

import (
	"regexp"
	"strings"
)

var reUpper = regexp.MustCompile(`[A-Z]+`)

func ToKebab(s string) string {
	keb := reUpper.ReplaceAllStringFunc(s, func(s string) string {
		s = strings.ToLower(s)
		if l := len(s); 2 <= l {
			s = s[:l-1] + "-" + string(s[l-1])
		}
		return "-" + s
	})
	if keb[0] == '-' {
		return keb[1:]
	}
	return keb
}
