package main

import (
	"log"
	"path/filepath"

	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
)

func main() {
	pathToTask := filepath.Join(".", "testdata", "kvadrputekl")
	task, err := fstaskparser.Read(pathToTask)
	if err != nil {
		log.Fatalf("failed to read task: %v", err)
	}
	newPath := filepath.Join(".", "output")
	err = task.Store(newPath)
	if err != nil {
		log.Fatalf("failed to store task: %v", err)
	}
}
