package fstaskparser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const proglvFSTaskFormatSpecVersOfScript = "v2.3.0"

func (task *task) Store(dirPath string) error {
	log.Printf("Starting to store task to directory: %s\n", dirPath)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		log.Printf("Directory already exists: %s\n", dirPath)
		return fmt.Errorf("directory already exists: %s", dirPath)
	}

	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		log.Printf("Error creating directory: %v\n", err)
		return fmt.Errorf("error creating directory: %w", err)
	}

	pToml, err := task.encodeProblemTOML()
	if err != nil {
		log.Printf("Error encoding problem.toml: %v\n", err)
		return fmt.Errorf("error encoding problem.toml: %w", err)
	}

	err = os.WriteFile(filepath.Join(dirPath, "problem.toml"), pToml, 0644)
	if err != nil {
		log.Printf("Error writing problem.toml: %v\n", err)
		return fmt.Errorf("error writing problem.toml: %w", err)
	}
	log.Println("problem.toml written successfully")

	// create tests directory
	testsDirPath := filepath.Join(dirPath, "tests")
	err = os.Mkdir(testsDirPath, 0755)
	if err != nil {
		log.Printf("Error creating tests directory: %v\n", err)
		return fmt.Errorf("error creating tests directory: %w", err)
	}
	log.Println("tests directory created successfully")

	for i, t := range task.tests {
		var inPath string
		var ansPath string

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
			log.Printf("Error writing input file %s: %v\n", inPath, err)
			return fmt.Errorf("error writing input file: %w", err)
		}

		err = os.WriteFile(ansPath, t.Answer, 0644)
		if err != nil {
			log.Printf("Error writing answer file %s: %v\n", ansPath, err)
			return fmt.Errorf("error writing answer file: %w", err)
		}
	}
	log.Println("Test files written successfully")

	// create examples directory
	examplesDirPath := filepath.Join(dirPath, "examples")
	err = os.Mkdir(examplesDirPath, 0755)
	if err != nil {
		log.Printf("Error creating examples directory: %v\n", err)
		return fmt.Errorf("error creating examples directory: %w", err)
	}
	log.Println("examples directory created successfully")

	for i, e := range task.examples {
		var inPath string
		var ansPath string

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
			log.Printf("Error writing input file %s: %v\n", inPath, err)
			return fmt.Errorf("error writing input file: %w", err)
		}

		err = os.WriteFile(ansPath, e.Output, 0644)
		if err != nil {
			log.Printf("Error writing answer file %s: %v\n", ansPath, err)
			return fmt.Errorf("error writing answer file: %w", err)
		}
	}
	log.Println("Example files written successfully")

	log.Printf("Task successfully stored in directory: %s\n", dirPath)
	return nil
}
