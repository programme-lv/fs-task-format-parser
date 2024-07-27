package fstaskparser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const proglvFSTaskFormatSpecVersOfScript = "v2.4.0"

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

	err = task.storeMdStatements(filepath.Join(dirPath, "statements", "md"))
	if err != nil {
		log.Printf("Error storing Markdown statements: %v\n", err)
		return fmt.Errorf("error storing Markdown statements: %w", err)
	}
	log.Println("Markdown statements written successfully")

	err = task.storeAssets(filepath.Join(dirPath, "assets"))
	if err != nil {
		log.Printf("Error storing assets: %v\n", err)
		return fmt.Errorf("error storing assets: %w", err)
	}
	log.Println("assets written successfully")

	log.Printf("Task successfully stored in directory: %s\n", dirPath)
	return nil
}

func (task *Task) storeAssets(assetDir string) error {
	err := os.MkdirAll(assetDir, 0755)
	if err != nil {
		log.Printf("Error creating assets directory: %v\n", err)
		return fmt.Errorf("error creating assets directory: %w", err)
	}
	log.Println("Assets directory created successfully")

	for _, v := range task.assets {
		// v.Content
		// v.RelativePath
		path := filepath.Join(assetDir, v.RelativePath)
		err = os.WriteFile(path, v.Content, 0644)
		if err != nil {
			log.Printf("Error writing asset: %v\n", err)
			return fmt.Errorf("error writing asset: %w", err)
		}
		log.Printf("Asset written: %s\n", path)
	}
	return nil
}

func (task *Task) storeMdStatements(mdStatementDir string) error {
	err := os.MkdirAll(mdStatementDir, 0755)
	if err != nil {
		log.Printf("Error creating Markdown statements directory: %v\n", err)
		return fmt.Errorf("error creating Markdown statements directory: %w", err)
	}
	log.Println("Markdown statements directory created successfully")

	for _, v := range task.mdStatements {
		// create language directory
		dirPath := filepath.Join(mdStatementDir, *v.Language)
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Printf("Error creating Markdown statement directory: %v\n", err)
			return fmt.Errorf("error creating Markdown statement directory: %w", err)
		}
		log.Printf("Markdown statement directory created: %s\n", dirPath)

		inputPath := filepath.Join(dirPath, "input.md")
		outputPath := filepath.Join(dirPath, "output.md")
		storyPath := filepath.Join(dirPath, "story.md")
		scoringPath := filepath.Join(dirPath, "scoring.md")
		notesPath := filepath.Join(dirPath, "notes.md")

		if v.Input != "" {
			err = os.WriteFile(inputPath, []byte(v.Input), 0644)
			if err != nil {
				log.Printf("Error writing Markdown statement: %v\n", err)
				return fmt.Errorf("error writing Markdown statement: %w", err)
			}
			log.Printf("Markdown statement written to: %s\n", inputPath)
		}

		if v.Output != "" {
			err = os.WriteFile(outputPath, []byte(v.Output), 0644)
			if err != nil {
				log.Printf("Error writing Markdown statement: %v\n", err)
				return fmt.Errorf("error writing Markdown statement: %w", err)
			}
			log.Printf("Markdown statement written to: %s\n", outputPath)
		}

		if v.Story != "" {
			err = os.WriteFile(storyPath, []byte(v.Story), 0644)
			if err != nil {
				log.Printf("Error writing Markdown statement: %v\n", err)
				return fmt.Errorf("error writing Markdown statement: %w", err)
			}
			log.Printf("Markdown statement written to: %s\n", storyPath)
		}

		if v.Scoring != nil {
			err = os.WriteFile(scoringPath, []byte(*v.Scoring), 0644)
			if err != nil {
				log.Printf("Error writing Markdown statement: %v\n", err)
				return fmt.Errorf("error writing Markdown statement: %w", err)
			}
			log.Printf("Markdown statement written to: %s\n", scoringPath)
		}

		if v.Notes != nil {
			err = os.WriteFile(notesPath, []byte(*v.Notes), 0644)
			if err != nil {
				log.Printf("Error writing Markdown statement: %v\n", err)
				return fmt.Errorf("error writing Markdown statement: %w", err)
			}
			log.Printf("Markdown statement written to: %s\n", notesPath)
		}
	}

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
