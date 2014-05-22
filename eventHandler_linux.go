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
	// TODO: make into a slice to cache multiple cookies
	var moveFromEvent *inotify.Event

	for {
		select {
		case event := <-watcher.Event:
			switch {
			default:
				fmt.Println(event, time.Now())

			case event.Mask == DIR_CREATE:
				// This checks if a newly created directory has children.
				// If it does it adds them to the watch
				// This is to deal with the 'mdkir -p' problem
				paths := CollectPaths([]string{event.Name})
				if len(paths) > 1 {
					for i := 0; i < len(paths); i++ {
						watcher.Watch(paths[i])
					}
				} else {
					watcher.Watch(event.Name)
				}
				fmt.Println(event.String(), time.Now())

			case event.Mask == FILE_CREATE:
				// Create encryption key, put into keymanager
				// Upload as if a modify event
				fmt.Println(event.String(), time.Now())

			case event.Mask == FILE_MODIFY:
				// Encrypt && upload
				fmt.Println(event.String(), time.Now())

			case event.Mask == DIR_DELETE:
				// Signal server for a delete
				watcher.RemoveWatch(event.Name)
				fmt.Println(event.String(), time.Now())

			case event.Mask == FILE_DELETE:
				// Signal server for a delete, if has children, delete them as well.
				// This would present a case where the parent is deleted first and the
				// child delete events come after.
				fmt.Println(event.String(), time.Now())

			case event.Mask == DIR_MOVE_FROM:
				// When a dir is moved this will trigger possibly a lot more move events if there are lots of children
				// This is another reason for the cookie cache, anytime a move event happens, append the cookie.
				// Don't make it a static size.
				moveFromEvent = event
				watcher.RemoveWatch(event.Name)
				fmt.Println(event.String(), time.Now())

			case event.Mask == DIR_MOVE_TO:
				if event.Cookie == moveFromEvent.Cookie {
					watcher.Watch(event.Name)
					fmt.Println("\t", event.String(), time.Now())
				}

			case event.Mask == FILE_MOVE_FROM:
				// Check stored cookie from last move event, if it is the same, follow through with the move.
				// TODO: make this a cookie cache since we can have multiple storage Roots, this could cause misses of move events
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
