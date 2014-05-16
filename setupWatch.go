package main

import (
	"fmt"

	"code.google.com/p/go.exp/fsnotify"
)

func SetupWatch(paths []string, excludes []string) (int, *fsnotify.Watcher) {
	var watchedCount int

	paths = CollectPaths(paths)
	excludes = CollectPaths(excludes)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error establishing watcher: ", err)
	}

	// establish watches
	for _, path := range paths {
		if IndexOf(path, excludes) == -1 {
			err = watcher.Watch(path)
			if err != nil {
				fmt.Println("Error: ", err, "  establishing watch on: ", path)
			}
			watchedCount++
		}
	}
	return watchedCount, watcher
}
