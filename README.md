fsmonitor
=========

A cross platform filesystem monitor in golang.

fsmonitor depends on the following packages:
    [queue from gocommons](https://github.com/hishboy/gocommons)
    [fsnotify](https://github.com/howeyc/fsnotify)

usage:
    import (
            "github.com/gaigepr/fsmonitor"
            "fmt"
            "time"
    )

    func main() {
            pathsToWatch := []string{"/path/to/watch1/", "/path/to/watch2/"}
            excludedFromWatch := []string{"/path/to/watch1/no_watch_2/", "/path/to/watch2/no_watch_1/"}
            MonitorFileSystem(pathsToWatch, excludeFromWatch, func(eventQueue *lang.Queue) {
                    for {
	    		    time.Sleep(1 * time.Second)
    			    fmt.Println(eventQueue.Poll())
		    }
            })
    }