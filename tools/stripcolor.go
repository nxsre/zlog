package tools

import "regexp"

func StripAnsiColor(str string) string {
	var (
		replacements = []struct {
			regex *regexp.Regexp
			repl  string
		}{
			{
				regex: regexp.MustCompile(`\[(3[1-8]|0)m`),
				repl:  "\033[${1}m",
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
