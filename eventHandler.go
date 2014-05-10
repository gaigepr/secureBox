package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"code.google.com/p/go.exp/fsnotify"
)

type Event struct {
	Name      string
	FilePath  string
	IsDir     bool
	EventType int
	EventTime time.Time
}

func PackageEvent(event *fsnotify.FileEvent) Event {
	eventType := func() int {
		switch {
		case event.IsCreate():
			return fsnotify.FSN_CREATE
		case event.IsModify():
			return fsnotify.FSN_MODIFY
		case event.IsDelete():
			return fsnotify.FSN_DELETE
		case event.IsRename():
			return fsnotify.FSN_RENAME
		}
		return 1
	}()

	var isDir bool
	if !event.IsDelete() && !event.IsRename() {
		info, err := os.Lstat(event.Name)
		if err != nil {
			fmt.Println("Error in PackageEvent checking on file: ", err)
		}
		isDir = info.IsDir()
	}

	re, err := regexp.Compile(".+/(.+)$")
	if err != nil {
		fmt.Println("Error compiling regexp")
	}
	result := re.FindStringSubmatch(event.Name)

	return Event{
		result[1],
		result[0],
		isDir,
		eventType,
		time.Now(),
	}
}

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

		// This reading of the file causes a modify event
		//ReadAndEncrypt(event.FilePath)
	}
}
