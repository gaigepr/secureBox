//package fsmonitor
package main

import (
	"code.google.com/p/go.exp/fsnotify"

	"fmt"
	"os"
	"path/filepath"
)



type Handler func([]string, []string)//watcher *fsnotify.Watcher)

func isMember(item string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return true
		}
	}
	return false
}

func collectPaths(paths []string) []string {
	newPaths := make([]string, 1, 2)

	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				newPaths = append(paths, path)
			}
			return err
		})

		if err != nil {
			fmt.Println("Error walking tree: ", err, " on file: ", path)
		}
	}
	return newPaths
}

func MonitorFileSystem(paths []string, excludes []string, handleEvents Handler) {
	var watchedCount int = 0
	paths = collectPaths(paths)
	excludes = collectPaths(excludes)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// report err up?
		fmt.Println("Error creating watcher: ", err)
	}

	// establish watches
	for _, path := range paths {
		if !isMember(path, excludes) {
			err = watcher.Watch(path)
			if err != nil {
				// report error up?
				fmt.Println("Error: ", err, "  establishing watch on: ", path)
			}
			watchedCount++
		}
	}
	fmt.Println("Directories watched: ", watchedCount)

	// call event handler
	handleEvents(paths, excludes)

}

func main() {
	a := []string{"/home/gaige/go/src/github.com/gaigepr/fsmonitor/"}
	b := []string{"/home/gaige/go/src/github.com/gaigepr/fsmonitor/.git/"}
	MonitorFileSystem(a, b, func(c, d []string) {
		fmt.Println("Into event handler!")
	})
}
