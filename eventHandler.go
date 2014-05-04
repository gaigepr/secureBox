package main

import (
	"fmt"
	"regexp"
	"time"

	"code.google.com/p/go.exp/fsnotify"
	"github.com/hishboy/gocommons/lang"
)

type Event struct {
	Name      string
	FilePath  string
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

	re, err := regexp.Compile("(.+/)/(.+)")
	if err != nil {
		fmt.Println("Problem compiling regexp")
	}
	result := re.FindStringSubmatch(event.Name)

	return Event{
		result[2],
		result[1],
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
		}
		time.Sleep(time.Millisecond * 500)
	}
}
