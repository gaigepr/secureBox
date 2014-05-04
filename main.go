package main

import (
	"fmt"
	"os"
	"os/signal"

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
			select {
			case ev := <-watcher.Event:
				eventQueue.Push(PackageEvent(ev))
			}
		}
	}()

	go EventHandler(eventQueue)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got Signal: ", s)
}
