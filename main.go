package main

import (
	"fmt"
	"time"

	"github.com/hishboy/gocommons/lang"
)

func main() {
	var eventQueue *lang.Queue = lang.NewQueue()

	paths := []string{"/home/gaige/testing/"}
	excludes := []string{"/home/gaige/testing/no_watch/"}

	// Create watch, add directories to it
	watchedCount, watcher := SetupWatch(paths, excludes)
	fmt.Println("Directories watched: ", watchedCount)

	go func() {
		for {
			ev := <-watcher.Event
			if ev != nil {
				eventQueue.Push(ev)
			}
		}
	}()

	for {
		time.Sleep(time.Millisecond * 1000)
		fmt.Println(eventQueue.Poll())
	}

}
