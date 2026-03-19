package shared

import "regexp"

var idPattern = regexp.MustCompile(`[0-9a-f]{32}`)

func ExtractID(value string) string {
	return idPattern.FindString(value)
}
