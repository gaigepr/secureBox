//package fsmonitor
package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"github.com/hishboy/gocommons/lang"

	"fmt"
	"os"
	"path/filepath"
	"time"
)



type Handler func([]string, []string, *lang.Queue)

func isMember(item string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return true
		}
	}
	return false
}

func collectPaths(paths []string) []string {
	newPaths := make([]string, 1, 1)

	for _, thisPath := range paths {
		err := filepath.Walk(thisPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				newPaths = append(newPaths, path)
		}
			return nil
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(newPaths)
	return newPaths
}

func MonitorFileSystem(paths []string, excludes []string, handleEvents Handler) {
	var watchedCount int = 0
	paths = collectPaths(paths)
	excludes = collectPaths(excludes)

	eventQueue := lang.NewQueue()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// report err up?
		fmt.Println("Error creating watcher: ", err)
	}

	// establish watches
	for _, path := range paths {
		if !(isMember(path, excludes)) {
			err = watcher.Watch(path)
			if err != nil {
				// report error up?
				fmt.Println("Error: ", err, "  establishing watch on: ", path)
			}
			watchedCount++
		}
	}

	fmt.Println("Directories watched: ", watchedCount)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				eventQueue.Push(ev)
				//fmt.Println(ev)
			case err := <-watcher.Error:
				fmt.Println(err)
			}
		}
	}()

	// call event handler
	handleEvents(paths, excludes, eventQueue)
}

func main() {
	a := []string{"/home/gaige/Dropbox/school/"}
	b := []string{"/home/gaige/Dropbox/school/2013-2014/cs_301/"}
	MonitorFileSystem(a, b, func(c, d []string, eventQueue *lang.Queue) {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println(eventQueue.Poll())
		}
	})
}
