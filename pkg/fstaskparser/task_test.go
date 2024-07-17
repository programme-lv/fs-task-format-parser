package fstaskparser_test

import (
	"path/filepath"
	"testing"

	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
)

func TestStoringComplexTask(t *testing.T) {
	testDir := filepath.Join(".", "..", "..", "testdata", "kvadrputekl")
	task, err := fstaskparser.Read(testDir)
	if err != nil {
		t.Fatalf("failed to read task: %v", err)
	}

	err = task.Store(filepath.Join(".", "output"))
	if err != nil {
		t.Fatalf("failed to store task: %v", err)
	}
}
