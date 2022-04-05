package eco

import (
	"regexp"
	"strings"
)

var regexpMatchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var regexpMatchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// toSnakeCase converts a string to snake case.
func toSnakeCase(in string) (out string) {
	out = regexpMatchFirstCap.ReplaceAllString(in, "${1}_${2}")
	out = regexpMatchAllCap.ReplaceAllString(out, "${1}_${2}")
	out = strings.ToLower(out)
	return
}

// defaultEnvNameTransformerFunc is the default function for naming the environment variables.
func defaultEnvNameTransformerFunc(parts []string, sep string) string {
	return strings.ToUpper(strings.Join(parts, sep))
}
