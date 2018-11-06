package hotcompiler

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"gobible/logmanager/cli/http/services/util"
)

// 热编译
func HotCompile()  {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	files := make([]string,0)

	files = append(files, "main.go")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
					//files = append(files, "main.go")
					Autobuild(files,util.GetCurrentDirectory())
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./dao/mysql")
	if err != nil {
		log.Fatal(err)
	}
	<-done

}
