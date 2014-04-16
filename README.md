fsmonitor
=========

A cross platform filesystem monitor in golang.

fsmonitor depends on the following packages:
    [queue from gocommons](https://github.com/hishboy/gocommons)
    [fsnotify](https://github.com/howeyc/fsnotify)

usage:


func main() {
	a := []string{"/home/user/dir_1/", /home/user/configs_and_videos/}
	b := []string{"/home/user/dir_1/ignore_this_dir/"}
	MonitorFileSystem(a, b, func(eventQueue *lang.Queue) {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println(eventQueue.Poll())
		}
	})
}
