fsmonitor
=========

A cross platform filesystem monitor in golang.

fsmonitor depends on the following packages:

[queue from gocommons](https://github.com/hishboy/gocommons) and [fsnotify](https://github.com/howeyc/fsnotify)

usage:

```
package main

import (
	"github.com/gaigepr/fsmonitor"
	"fmt"
	"time"
)

func main() {
	pathsToWatch := []string{"/home/user/dir_1/", "/home/user/configs_and_videos/"}
	PathsToExcludeFromWatch := []string{"/home/user/dir_1/ignore_this_dir/"}
	MonitorFileSystem(pathsToWatch, pathsToExcludeFromWatch, func(eventQueue *lang.Queue) {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println(eventQueue.Poll())
		}
	})
}
```
This example will pull events off the eventQueue every second and print either: 
```
<nil>
``` 
Or the path where the event happened and the event type like so:
```
"/home/user/dir_1//file.txt": MODIFY|ATTRIB
```
