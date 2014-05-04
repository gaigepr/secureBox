package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"code.google.com/p/go.exp/fsnotify"
	"github.com/hishboy/gocommons/lang"
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

func EventHandler(eventQueue *lang.Queue) {
	fmt.Println("In event handler")
	for {
		if eventQueue.Len() > 0 {
			// Type assert to pullthe struct out of the interface
			event := eventQueue.Poll().(Event)
			fmt.Println(event)
			time.Sleep(time.Millisecond * 500)
			if !event.IsDir {
				// This reading of the file causes a modify event
				ReadAndEncrypt(event.FilePath)
			}
		} else {
			time.Sleep(time.Millisecond * 50)
		}
	}
}
