package utils_test

import (
	"testing"

	"github.com/youta-t/flarc/utils"
	"github.com/youta-t/its"
)

func TestToKebab(t *testing.T) {
	theory := func(when string, then its.Matcher[string]) func(*testing.T) {
		return func(t *testing.T) {
			got := utils.ToKebab(when)
			then.Match(got).OrError(t)
		}
	}

	t.Run("snailCase", theory(
		"snailCaseString",
		its.EqEq("snail-case-string"),
	))

	t.Run("CamelCase", theory(
		"CamelCaseString",
		its.EqEq("camel-case-string"),
	))

	t.Run("alllowercase", theory(
		"alllowercasestring",
		its.EqEq("alllowercasestring"),
	))

	t.Run("having UPPER word", theory(
		"receivedHTTPMethod",
		its.EqEq("received-http-method"),
	))
}
