package fstaskparser

import (
	"fmt"
	"os"
	"path/filepath"
)

const proglvFSTaskFormatSpecVersOfScript = "v2.3.0"

func (task *Task) Store(dirPath string) error {
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory already exists: %s", dirPath)
	}

	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	pToml, err := task.encodeProblemTOML()
	if err != nil {
		return fmt.Errorf("error encoding problem.toml: %w", err)
	}

	err = os.WriteFile(filepath.Join(dirPath, "problem.toml"), pToml, 0644)
	if err != nil {
		return fmt.Errorf("error writing problem.toml: %w", err)
	}

	// create tests directory
	testsDirPath := filepath.Join(dirPath, "tests")
	err = os.Mkdir(testsDirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating tests directory: %w", err)
	}

	for i, t := range task.tests {
		// create input file {name}.in
		// create answer file {name}.out
		// use name for {name} if it exists
		// otherwise use padded id (3 digits)

		var inPath string = ""
		var ansPath string = ""

		if fname, ok := task.testIDToFilename[t.ID]; ok {
			inPath = filepath.Join(testsDirPath, fname+".in")
			ansPath = filepath.Join(testsDirPath, fname+".out")
		} else {
			inName := fmt.Sprintf("%03d.in", i+1)
			ansName := fmt.Sprintf("%03d.out", i+1)
			inPath = filepath.Join(testsDirPath, inName)
			ansPath = filepath.Join(testsDirPath, ansName)
		}

		err = os.WriteFile(inPath, t.Input, 0644)
		if err != nil {
			return fmt.Errorf("error writing input file: %w", err)
		}

		err = os.WriteFile(ansPath, t.Answer, 0644)
		if err != nil {
			return fmt.Errorf("error writing answer file: %w", err)
		}
	}

	examplesDirPath := filepath.Join(dirPath, "examples")
	err = os.Mkdir(examplesDirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating examples directory: %w", err)
	}

	for i, e := range task.examples {
		// treat tests and examples the same way

		var inPath string = ""
		var ansPath string = ""

		if e.Name != nil {
			inPath = filepath.Join(examplesDirPath, *e.Name+".in")
			ansPath = filepath.Join(examplesDirPath, *e.Name+".out")
		} else {
			inName := fmt.Sprintf("%03d.in", i+1)
			ansName := fmt.Sprintf("%03d.out", i+1)
			inPath = filepath.Join(examplesDirPath, inName)
			ansPath = filepath.Join(examplesDirPath, ansName)
		}

		err = os.WriteFile(inPath, e.Input, 0644)
		if err != nil {
			return fmt.Errorf("error writing input file: %w", err)
		}

		err = os.WriteFile(ansPath, e.Output, 0644)
		if err != nil {
			return fmt.Errorf("error writing answer file: %w", err)
		}
	}

	return nil
}
