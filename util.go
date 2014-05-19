package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"code.google.com/p/go.exp/inotify"
)

type Event struct {
	Name      string
	FilePath  string
	IsDir     bool
	EventType uint32
	EventTime time.Time
}

func PackageEvent(event *inotify.Event) Event {
	var isDir bool
	if event.Mask == inotify.IN_ISDIR {
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
		event.Mask,
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
