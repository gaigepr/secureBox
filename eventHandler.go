package main

import (
	"fmt"

	"code.google.com/p/go.exp/fsnotify"
)

func EventHandler(eventChan chan Event, watcher *fsnotify.Watcher) {
	fmt.Println("In event handler")
	for {
		event := <-eventChan
		switch {
		case event.EventType == fsnotify.FSN_CREATE:
			watcher.Watch(event.FilePath)
			if event.IsDir {
				allPaths := CollectPaths([]string{event.FilePath})
				//fmt.Println(allPaths[1:])
				//time.Sleep(1 * time.Millisecond)
				for _, path := range allPaths[1:] {
					err := watcher.Watch(path)
					if err != nil {
						fmt.Println("Error: ", err, "  establishing watch on: ", path)
					}
				}
			}
			fmt.Printf("EVENT %s: %s \t%v\n", "CREATE", event.FilePath, event.EventTime)

		case event.EventType == fsnotify.FSN_MODIFY:
			fmt.Printf("EVENT %s: %s \t%v\n", "MODIFY", event.FilePath, event.EventTime)

		case event.EventType == fsnotify.FSN_DELETE:
			if event.IsDir {
				err := watcher.RemoveWatch(event.FilePath)
				if err != nil {
					fmt.Println("Error removing watchg (delete): ", err)
				}
			}
			fmt.Printf("EVENT %s: %s \t%v\n", "DELETE", event.FilePath, event.EventTime)

		case event.EventType == fsnotify.FSN_RENAME:
			if event.IsDir {
				err := watcher.RemoveWatch(event.FilePath)
				if err != nil {
					fmt.Println("Error removing watch (rename): ", err)
				}
			}
			fmt.Printf("EVENT %s: %s \t%v\n", "RENAME", event.FilePath, event.EventTime)

		}
	}
}
