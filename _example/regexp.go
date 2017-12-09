package main

import (
	"fmt"
	"regexp"
)

func main() {
	fmt.Println(DeSensitization(`-p 111111`))
}

func DeSensitization(s string) string {
	re := regexp.MustCompile("([-\"\b]?)(password|-p)([ \t=:\b\"]+)([a-zA-Z0-9-_@\\*]{3,})(\"?)")
	s = re.ReplaceAllString(s, "$1$2$3******$5")
	return s
}
