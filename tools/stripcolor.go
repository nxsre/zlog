package tools

import "regexp"

// https://github.com/marthjod/recolor
func StripAnsiColor(str string) string {
	var (
		replacements = []struct {
			regex *regexp.Regexp
			repl  string
		}{
			{
				regex: regexp.MustCompile(`\x1b\[[0-9;]*m`),
				repl:  "",
			},
			{
				regex: regexp.MustCompile(`\[39m`),
				repl:  "",
			},
		}
	)

	for _, replacement := range replacements {
		if replacement.regex.MatchString(str) {
			str = replacement.regex.ReplaceAllString(str, replacement.repl)
		}
	}
	return str
}
