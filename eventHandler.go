package main

import (
	"fmt"
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
			return 1
		case event.IsModify():
			return 2
		case event.IsDelete():
			return 4
		case event.IsRename():
			return 8
		}
		return 1
	}()

	//fmt.Println("EVENTTTTT: ", event)

	return Event{
		event.Name,
		event.Name,
		eventType,
		time.Now(),
	}

}

func EventHandler(eventQueue *lang.Queue) {
	fmt.Println("In event handler")
	for {

		// Events put into the queue lose most of their info.
		// Before pushing onto the queue must make own new
		// struct that conatins relevant info + timestamp etc
		if eventQueue.Len() > 0 {
			// Type assert to pullthe struct out of the interface
			event := eventQueue.Poll().(Event)
			fmt.Println(event)
		}
		time.Sleep(time.Millisecond * 500)
	}
}
