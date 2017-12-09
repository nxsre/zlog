package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("111")
	var ch = make(chan []byte)
	go func() {

		time.Sleep(5 * time.Second)
		a := <-ch
		fmt.Println(string(a))
	}()

	ch <- []byte("11111")
	fmt.Println("写入成功")
}
