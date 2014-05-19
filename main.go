package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		TestCrypto()
		return
	}

	paths := []string{"/home/gaige/testing/"}
	excludes := []string{"/home/gaige/testing/no_watch/"}

	// Create watch, add directories to it
	watchedCount, watcher := SetupWatch(paths, excludes)
	fmt.Println("Directories watched: ", watchedCount)

	eventChan := make(chan Event)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				eventChan <- PackageEvent(ev)
			case err := <-watcher.Error:
				fmt.Println(err)
			}
		}
	}()

	go EventHandler(eventChan, watcher)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	for {
		select {
		case s := <-c:
			fmt.Println("Got Signal: ", s)
			return
		}
	}
}
