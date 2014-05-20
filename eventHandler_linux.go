package main

import (
	"fmt"
	"time"

	"code.google.com/p/go.exp/inotify"
)

func EventHandler(watcher *inotify.Watcher) {
	for {
		select {
		case event := <-watcher.Event:
			switch {
			case event.Mask == (inotify.IN_CREATE + inotify.IN_ISDIR):
				fmt.Printf("EVENT %s: %s \t%v\n", "DIR CREATE", event.Name, time.Now())

			case event.Mask == inotify.IN_CREATE:
				fmt.Printf("EVENT %s: %s \t%v\n", "CREATE", event.Name, time.Now())

			case event.Mask == inotify.IN_CLOSE_WRITE:
				fmt.Printf("EVENT %s: %s \t%v\n", "CLOSE_WRITE", event.Name, time.Now())

			case event.Mask == (inotify.IN_DELETE + inotify.IN_ISDIR):
				fmt.Printf("EVENT %s: %s \t%v\n", "DIR DELETE", event.Name, time.Now())

			case event.Mask == inotify.IN_DELETE:
				fmt.Printf("EVENT %s: %s \t%v\n", "DELETE", event.Name, time.Now())

			case event.Mask == inotify.IN_MOVE:
				fmt.Printf("EVENT %s: %s \t%v\n", "MOVE", event.Name, time.Now())

			}

		case err := <-watcher.Error:
			fmt.Println("WATCHER ERROR: ", err)
		}
	}
}
