package fstaskparser

import (
	"fmt"
	"os"
	"path/filepath"
)

const proglvFSTaskFormatSpecVersion = "2.3"

func (t *Task) Store(dirPath string) error {
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory already exists: %s", dirPath)
	}

	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	pToml, err := t.encodeProblemTOML()
	if err != nil {
		return fmt.Errorf("error encoding problem.toml: %w", err)
	}

	err = os.WriteFile(filepath.Join(dirPath, "problem.toml"), pToml, 0644)
	if err != nil {
		return fmt.Errorf("error writing problem.toml: %w", err)
	}

	return nil
}
