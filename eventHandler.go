package main

import (
	"fmt"

	"code.google.com/p/go.exp/inotify"
)

func EventHandler(eventChan chan Event, watcher *inotify.Watcher) {
	fmt.Println("In event handler")
	for {
		event := <-eventChan
		switch {
		case event.EventType == inotify.IN_CREATE:
			watcher.Watch(event.FilePath)

			if event.IsDir {
				// TODO: As the scan happens create new Events and push them onto the eventChan
				allPaths := CollectPaths([]string{event.FilePath})
				for _, path := range allPaths[1:] {
					err := watcher.Watch(path)
					if err != nil {
						fmt.Println("Error: ", err, "  establishing watch on: ", path)
					}
				}
			}

			// genEncryptionShit(event.IsDir()) -- Make AES key, Add key yo key tree, etc
			// log() -- because logs are cool
			// UploadShit(event) -- upload to server

			fmt.Printf("EVENT %s: %s \t%v\n", "CREATE", event.FilePath, event.EventTime)

		case event.EventType == inotify.IN_CLOSE_WRITE:

			// log() -- because logs are cool
			// UploadShit(event) -- upload to server

			fmt.Printf("EVENT %s: %s \t%v\n", "CLOSE_WRITE", event.FilePath, event.EventTime)

		case event.EventType == inotify.IN_DELETE:
			if event.IsDir {
				err := watcher.RemoveWatch(event.FilePath)
				if err != nil {
					fmt.Println("Error removing watchg (delete): ", err)
				}
			}
			fmt.Printf("EVENT %s: %s \t%v\n", "DELETE", event.FilePath, event.EventTime)

		case event.EventType == inotify.IN_MOVE:
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
