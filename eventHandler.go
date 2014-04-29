package main

import(
	"fmt"
	"time"

	"github.com/hishboy/gocommons/lang"
)

func EventHandler(eventQueue *lang.Queue) {
	fmt.Println("In event handler")
	for {
		time.Sleep(time.Millisecond * 1000)
		fmt.Println(eventQueue.Poll())
	}
}
