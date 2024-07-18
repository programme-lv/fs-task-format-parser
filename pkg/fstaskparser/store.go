package fstaskparser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const proglvFSTaskFormatSpecVersOfScript = "v2.3.0"

func (task *Task) Store(dirPath string) error {
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

	err = task.storeProblemToml(filepath.Join(dirPath, "problem.toml"))
	if err != nil {
		log.Printf("Error storing problem.toml: %v\n", err)
		return fmt.Errorf("error storing problem.toml: %w", err)
	}
	log.Println("problem.toml written successfully")

	err = task.storeTests(filepath.Join(dirPath, "tests"))
	if err != nil {
		log.Printf("Error storing tests: %v\n", err)
		return fmt.Errorf("error storing tests: %w", err)
	}
	log.Println("tests written successfully")

	err = task.storeExamples(filepath.Join(dirPath, "examples"))
	if err != nil {
		log.Printf("Error storing examples: %v\n", err)
		return fmt.Errorf("error storing examples: %w", err)
	}
	log.Println("examples written successfully")

	err = task.storePDFStatements(filepath.Join(dirPath, "statements", "pdf"))
	if err != nil {
		log.Printf("Error storing PDF statements: %v\n", err)
		return fmt.Errorf("error storing PDF statements: %w", err)
	}
	log.Println("PDF statements written successfully")

	log.Printf("Task successfully stored in directory: %s\n", dirPath)
	return nil
}

func (task *Task) storePDFStatements(pdfStatementsDir string) error {
	err := os.MkdirAll(pdfStatementsDir, 0755)
	if err != nil {
		log.Printf("Error creating PDF statements directory: %v\n", err)
		return fmt.Errorf("error creating PDF statements directory: %w", err)
	}
	log.Println("PDF statements directory created successfully")

	for k, v := range task.pdfStatements {
		// k is language, v is content
		fname := fmt.Sprintf("%s.pdf", k)
		fpath := filepath.Join(pdfStatementsDir, fname)
		err = os.WriteFile(fpath, []byte(v), 0644)
		if err != nil {
			log.Printf("Error writing PDF statement: %v\n", err)
			return fmt.Errorf("error writing PDF statement: %w", err)
		}
		log.Printf("PDF statement written to: %s\n", fpath)
	}

	return nil
}

func (task *Task) storeProblemToml(problemTomlPath string) error {
	pToml, err := task.encodeProblemTOML()
	if err != nil {
		log.Printf("Error encoding problem.toml: %v\n", err)
		return fmt.Errorf("error encoding problem.toml: %w", err)
	}
	err = os.WriteFile(problemTomlPath, pToml, 0644)
	if err != nil {
		log.Printf("Error writing problem.toml: %v\n", err)
		return fmt.Errorf("error writing problem.toml: %w", err)
	}
	log.Println("problem.toml written successfully")
	return nil
}

func (task *Task) storeTests(testsDirPath string) error {
	var err error
	err = os.Mkdir(testsDirPath, 0755)
	if err != nil {
		log.Printf("Error creating tests directory: %v\n", err)
		return fmt.Errorf("error creating tests directory: %w", err)
	}
	log.Println("tests directory created successfully")

	for _, t := range task.tests {
		fname := task.getTestToBeWrittenFname(t.ID)
		inPath := filepath.Join(testsDirPath, fname+".in")
		ansPath := filepath.Join(testsDirPath, fname+".out")

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

	return nil
}

func (task *Task) getTestToBeWrittenFname(id int) string {
	if fname, ok := task.testIDToFilename[id]; ok {
		return fname
	} else {
		return fmt.Sprintf("%03d", id)
	}
}

func (task *Task) getTestIDByFilenameOverwriteMap() map[string]int {
	res := map[string]int{}

	type testWrittenFilename struct {
		ID    int
		Fname string
	}
	order := []testWrittenFilename{}
	for _, t := range task.tests {
		order = append(order, testWrittenFilename{
			ID:    t.ID,
			Fname: task.getTestToBeWrittenFname(t.ID),
		})
	}

	sort.Slice(order, func(i, j int) bool {
		return order[i].Fname < order[j].Fname
	})

	for i, o := range order {
		received := i + 1
		actual := o.ID
		if received != actual {
			res[o.Fname] = received
		}
	}

	return res
}

func (task *Task) storeExamples(examplesDirPath string) error {
	var err error
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
	return nil
}
