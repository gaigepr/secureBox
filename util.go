package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func IndexOf(element string, array []string) int {
	for i := 0; i < len(array); i++ {
		if array[i] == element {
			return i
		}
	}
	return -1
}

func CollectPaths(paths []string) []string {
	// paths to be returned
	collectedPaths := make([]string, 1, 1)

	for _, thisPath := range paths {
		err := filepath.Walk(thisPath,
			// Function arg for filepath.Walk
			func(path string, info os.FileInfo, err error) error {
				if info == nil {
					fmt.Println("File or directory does not exist.", path)
				} else if info.IsDir() {
					collectedPaths = append(collectedPaths, path)
				}
				return nil
			})

		if err != nil {
			fmt.Println(err)
		}
	}
	return collectedPaths
}
