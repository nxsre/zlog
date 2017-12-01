package tools

import "regexp"

// 删除 xterm icon 和 title 的设置字符 \e]0; 和 \a 中间的所有字符串
func StripXtermCtl(str string) string {
	var (
		replacements = []struct {
			regex *regexp.Regexp
			repl  string
		}{
			{
				regex: regexp.MustCompile(`\x1b\]0;[^\a]+\a`),
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
