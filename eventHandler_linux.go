package main

import (
	"fmt"
	"time"

	"code.google.com/p/go.exp/inotify"
)

const (
	DIR_CREATE     = inotify.IN_CREATE + inotify.IN_ISDIR
	DIR_DELETE     = inotify.IN_DELETE + inotify.IN_ISDIR
	DIR_MOVE_FROM  = inotify.IN_MOVED_FROM + inotify.IN_ISDIR
	DIR_MOVE_TO    = inotify.IN_MOVED_TO + inotify.IN_ISDIR
	FILE_CREATE    = inotify.IN_CREATE
	FILE_MODIFY    = inotify.IN_CLOSE_WRITE
	FILE_DELETE    = inotify.IN_DELETE
	FILE_MOVE_FROM = inotify.IN_MOVED_FROM
	FILE_MOVE_TO   = inotify.IN_MOVED_TO
)

func EventHandler(watcher *inotify.Watcher) {
	var moveFromEvent *inotify.Event

	for {
		select {
		case event := <-watcher.Event:
			switch {
			default:
				fmt.Println(event, time.Now())

			case event.Mask == DIR_CREATE:
				watcher.Watch(event.Name)
				fmt.Println(event.String(), time.Now())

			case event.Mask == FILE_CREATE:
				fmt.Println(event.String(), time.Now())

			case event.Mask == FILE_MODIFY:
				fmt.Println(event.String(), time.Now())

			case event.Mask == DIR_DELETE:
				watcher.RemoveWatch(event.Name)
				fmt.Println(event.String(), time.Now())

			case event.Mask == FILE_DELETE:
				fmt.Println(event.String(), time.Now())

			case event.Mask == DIR_MOVE_FROM:
				moveFromEvent = event
				watcher.RemoveWatch(event.Name)
				fmt.Println(event.String(), time.Now())

			case event.Mask == DIR_MOVE_TO:
				if event.Cookie == moveFromEvent.Cookie {
					watcher.Watch(event.Name)
					fmt.Println("\t", event.String(), time.Now())
				}

			case event.Mask == FILE_MOVE_FROM:
				moveFromEvent = event
				fmt.Println(event.String(), time.Now())

			case event.Mask == FILE_MOVE_TO:
				if event.Cookie == moveFromEvent.Cookie {
					fmt.Println("\t", event.String(), time.Now())
				}

			}

		case err := <-watcher.Error:
			fmt.Println("WATCHER ERROR: ", err)
		}
	}
}
