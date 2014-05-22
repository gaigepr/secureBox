package main

import (
	"fmt"

	"code.google.com/p/go.exp/inotify"
)

func SetupWatch(paths []string, excludes []string) (int, *inotify.Watcher) {
	// How many directories are being watched
	var watchedCount int

	// Collect all subdirs of the watch and exclude roots
	paths = CollectPaths(paths)
	excludes = CollectPaths(excludes)

	watcher, err := inotify.NewWatcher()
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
