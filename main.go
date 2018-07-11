package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Watching dir:", dir)
	watchDir(dir)
}

func watchDir(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					// Check if modified file is go file.
					if strings.HasSuffix(event.Name, ".go") {
						runTests()
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func runTests() {
	fmt.Println("RUNNING TESTS")
	fmt.Println("--------------")

	cmd := exec.Command("go", "test", "-v")

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr

	err := cmd.Run()

	if err != nil {
		fmt.Println(outErr.String())
		// fmt.Println(out.String())
		// log.Fatal(err)
	}

	fmt.Println(out.String())
}
