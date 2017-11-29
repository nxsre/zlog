package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/soopsio/zlog/tools"
)

const src = "déjà vu" + // precomposed unicode
	"\n\000\037 \041\176\177\200\377\n" + // various boundary cases
	"as⃝df̅" // unicode combining characters

func main() {
	color.Blue("Prints %s in blue.", "text")
	fmt.Println("source text:")
	fmt.Println(src)
	fmt.Println("\nas bytes, stripped of control codes:")
	fmt.Println(tools.StripCtlFromBytes(src))
	fmt.Println("\nas bytes, stripped of control codes and extended characters:")
	fmt.Println(tools.StripCtlAndExtFromBytes(src))
	fmt.Println("\nas UTF-8, stripped of control codes:")
	fmt.Println(tools.StripCtlFromUTF8(src))
	fmt.Println("\nas UTF-8, stripped of control codes and extended characters:")
	fmt.Println(tools.StripCtlAndExtFromUTF8(src))
	fmt.Println("\nas decomposed and stripped Unicode:")
	fmt.Println(tools.StripCtlAndExtFromUnicode(src))
}
