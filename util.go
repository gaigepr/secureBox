package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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
		err := filepath.Walk(filepath.Clean(thisPath),
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
