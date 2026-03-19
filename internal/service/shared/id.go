package shared

import "regexp"

var idPattern = regexp.MustCompile(`([0-9a-f]{32}|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)

func ExtractID(value string) string {
	return idPattern.FindString(value)
}
