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
		if replacement.regex.MatchString(line) {
			line = replacement.regex.ReplaceAllString(line, replacement.repl)
		}
	}
	return line
}
