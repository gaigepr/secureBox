package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
)

type Configuration struct {
	Username     string
	Password     string
	WatchRoots   []string
	ExcludeRoots []string
}

func main() {
	// Command line arg(s)
	if !(len(os.Args) > 1) {
		panic("Not enough command line args: Need path to config file!")
	}

	// Load or create the config file
	config := Configuration{}
	file, err := os.Open(os.Args[1])
	if err != nil {
		// If the config file is not found we need to make it
		fmt.Println("Could not find config file or something: ", err)
		file, err = os.Create(os.Args[1]);
		if err != nil {
			fmt.Println("ERROR creating file: ", err)
		}
	}

	// Decode config file into Configuration struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("ERROR parsing json: ", err)
	}
	fmt.Println(config)

	// Create watch and add directories to it
	watchedCount, watcher := SetupWatch(config.WatchRoots, config.ExcludeRoots)
	fmt.Println("Directories watched: ", watchedCount)

	go EventHandler(watcher)

	// Wait for SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	for {
		select {
		case s := <-c:
			fmt.Println("Got Signal: ", s)
			return
		}
	}
}
