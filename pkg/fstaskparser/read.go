package fstaskparser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Read(dirPath string) (*task, error) {
	log.Printf("Starting to read directory: %s\n", dirPath)

	t := task{
		problemTomlContent:   []byte{},
		problemTags:          []string{},
		problemAuthors:       []string{},
		mdStatements:         []mDStatement{},
		taskName:             "",
		originOlympiad:       "",
		difficultyOneToFive:  0,
		memoryMegabytes:      0,
		cpuTimeSeconds:       0,
		examples:             []example{},
		exampleFilenameToID:  map[string]int{},
		visibleInputSubtasks: []int{},
		testFnamesSorted:     []string{},
		testFilenameToID:     map[string]int{},
		testIDOverwrite:      map[string]int{},
		testIDToFilename:     map[int]string{},
		tests:                []test{},
		testGroupIDs:         []int{},
		isTGroupPublic:       map[int]bool{},
		tGroupPoints:         map[int]int{},
		tGroupToStMap:        map[int]int{},
		tGroupTestIDs:        map[int][]int{},
		tGroupFnames:         map[int][]string{},
	}

	problemTomlPath := filepath.Join(dirPath, "problem.toml")
	log.Printf("Reading problem.toml from: %s\n", problemTomlPath)
	problemTomlContent, err := os.ReadFile(problemTomlPath)
	if err != nil {
		log.Printf("Error reading problem.toml: %v\n", err)
		return nil, fmt.Errorf("error reading problem.toml: %w", err)
	}

	t.problemTomlContent = problemTomlContent
	log.Println("problem.toml content read successfully")

	var specVersStruct struct {
		Specification string `toml:"specification"`
	}

	err = toml.Unmarshal(problemTomlContent, &specVersStruct)
	if err != nil {
		log.Printf("Failed to unmarshal the specification: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal the specification: %w", err)
	}

	specVers := specVersStruct.Specification
	if len(specVers) == 0 {
		log.Println("Empty specification found")
		return nil, fmt.Errorf("empty specification")
	}
	if specVers[0] == 'v' {
		specVers = specVers[1:]
	}

	log.Printf("Specification version: %s\n", specVers)

	semVersCmpRes, err := getCmpSemVersionsResult(specVers, proglvFSTaskFormatSpecVersOfScript)
	if err != nil {
		log.Printf("Error comparing sem versions: %v\n", err)
		return nil, fmt.Errorf("error comparing sem versions: %w", err)
	}

	if semVersCmpRes > 0 {
		log.Printf("Unsupported specification version (too new): %s\n", specVers)
		return nil, fmt.Errorf("unsupported specification version (too new): %s", specVers)
	}

	if semVersCmpRes < 0 {
		log.Printf("Warning: outdated specification version (too old): %s\n", specVers)
	}

	t.taskName, err = readTaskName(specVers, string(problemTomlContent))
	if err != nil {
		log.Printf("Error reading task name: %v\n", err)
		return nil, fmt.Errorf("error reading task name: %w", err)
	}

	t.cpuTimeSeconds, err = readCPUTimeLimitInSeconds(specVers, string(problemTomlContent))
	if err != nil {
		log.Printf("Error reading CPU time limit: %v\n", err)
		return nil, fmt.Errorf("error reading cpu time limit: %w", err)
	}

	t.memoryMegabytes, err = readMemoryLimitInMegabytes(specVers, string(problemTomlContent))
	if err != nil {
		log.Printf("Error reading memory limit: %v\n", err)
		return nil, fmt.Errorf("error reading memory limit: %w", err)
	}

	t.problemTags, err = readProblemTags(specVers, string(problemTomlContent))
	if err != nil {
		log.Printf("Error reading problem tags: %v\n", err)
		return nil, fmt.Errorf("error reading problem tags: %w", err)
	}

	t.problemAuthors, err = readProblemAuthors(specVers, string(problemTomlContent))
	if err != nil {
		log.Printf("Error reading problem authors: %v\n", err)
		return nil, fmt.Errorf("error reading problem authors: %w", err)
	}

	t.originOlympiad, err = readOriginOlympiad(specVers, string(problemTomlContent))
	if err != nil {
		log.Printf("Error reading origin olympiad: %v\n", err)
		return nil, fmt.Errorf("error reading origin olympiad: %w", err)
	}

	t.difficultyOneToFive, err = readDifficultyOneToFive(specVers, string(problemTomlContent))
	if err != nil {
		log.Printf("Error reading difficulty: %v\n", err)
		return nil, fmt.Errorf("error reading difficulty: %w", err)
	}

	log.Println("Reading test filenames from the tests directory")
	t.testFnamesSorted, err = readTestFNamesSorted(filepath.Join(dirPath, "tests"))
	if err != nil {
		log.Printf("Error reading test filenames: %v\n", err)
		return nil, fmt.Errorf("error reading test filenames: %w", err)
	}

	for i, fname := range t.testFnamesSorted {
		t.testFilenameToID[fname] = i + 1
		t.testIDToFilename[i+1] = fname
	}

	log.Println("Reading test ID overwrite")
	t.testIDOverwrite, err = readTestIDOverwrite(specVers, problemTomlContent)
	if err != nil {
		log.Printf("Error reading test ID overwrite: %v\n", err)
		return nil, fmt.Errorf("error reading test id overwrite: %w", err)
	}

	for k, v := range t.testIDOverwrite {
		t.testIDToFilename[v] = k
		t.testFilenameToID[k] = v
	}

	spottedFnames := make(map[int]bool)
	for _, fname := range t.testFnamesSorted {
		if _, ok := spottedFnames[t.testFilenameToID[fname]]; ok {
			log.Printf("Duplicate filename for ID: %s\n", fname)
			return nil, fmt.Errorf("duplicate filename for ID: %s", fname)
		}
		spottedFnames[t.testFilenameToID[fname]] = true
	}

	spottedIDs := make(map[string]bool)
	for _, id := range t.testIDToFilename {
		if _, ok := spottedIDs[id]; ok {
			log.Printf("Duplicate ID for filename: %s\n", id)
			return nil, fmt.Errorf("duplicate ID for filename: %s", id)
		}
		spottedIDs[id] = true
	}

	log.Println("Reading tests directory")
	t.tests, err = readTestsDir(dirPath, t.testFilenameToID)
	if err != nil {
		log.Printf("Error reading tests directory: %v\n", err)
		return nil, fmt.Errorf("error reading tests directory: %w", err)
	}

	log.Println("Reading examples directory")
	t.examples, err = readExamplesDir(dirPath)
	if err != nil {
		log.Printf("Error reading examples directory: %v\n", err)
		return nil, fmt.Errorf("error reading examples directory: %w", err)
	}

	log.Println("Reading test group IDs")
	t.testGroupIDs, err = readTestGroupIDs(specVers, problemTomlContent)
	if err != nil {
		log.Printf("Error reading test group IDs: %v\n", err)
		return nil, fmt.Errorf("error reading test group IDs: %w", err)
	}

	log.Println("Reading is test group public")
	t.isTGroupPublic, err = readIsTGroupPublic(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		log.Printf("Error reading is test group public: %v\n", err)
		return nil, fmt.Errorf("error reading is test group public: %w", err)
	}

	log.Println("Reading test group points")
	t.tGroupPoints, err = readTGroupPoints(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		log.Printf("Error reading test group points: %v\n", err)
		return nil, fmt.Errorf("error reading test group points: %w", err)
	}

	log.Println("Reading test group to subtask map")
	t.tGroupToStMap, err = readTGroupToStMap(specVers, problemTomlContent)
	if err != nil {
		log.Printf("Error reading test group to subtask map: %v\n", err)
		return nil, fmt.Errorf("error reading test group to subtask map: %w", err)
	}

	log.Println("Reading test group test IDs")
	t.tGroupTestIDs, err = readTGroupTestIDs(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		log.Printf("Error reading test group test IDs: %v\n", err)
		return nil, fmt.Errorf("error reading test group test IDs: %w", err)
	}

	log.Println("Reading test group filenames")
	t.tGroupFnames, err = readTGroupFnames(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		log.Printf("Error reading test group filenames: %v\n", err)
		return nil, fmt.Errorf("error reading test group filenames: %w", err)
	}

	for k, v := range t.tGroupFnames {
		for _, fname := range v {
			t.tGroupTestIDs[k] = append(t.tGroupTestIDs[k], t.testFilenameToID[fname])
		}
	}

	idsSpotted := make(map[int]bool)
	for _, v := range t.testGroupIDs {
		for _, id := range t.tGroupTestIDs[v] {
			if _, ok := idsSpotted[id]; ok {
				log.Printf("Duplicate test ID in test group: %d\n", id)
				return nil, fmt.Errorf("duplicate test ID in test group: %d", id)
			}
			idsSpotted[id] = true
		}
	}

	log.Println("Successfully read and parsed task")
	return &t, nil
}

func readTaskName(specVers string, tomlContent string) (string, error) {
	log.Printf("Reading task name for specification version: %s\n", specVers)
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.2")
	if err != nil {
		log.Printf("Error comparing semversions: %v\n", err)
		return "", fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		log.Printf("Unsupported specification version: %s\n", specVers)
		return "", fmt.Errorf("unsupported specification version: %s", specVers)
	}

	tomlStruct := struct {
		TaskName string `toml:"task_name"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		log.Printf("Failed to unmarshal the task name: %v\n", err)
		return "", fmt.Errorf("failed to unmarshal the task name: %w", err)
	}

	log.Printf("Successfully read task name: %s\n", tomlStruct.TaskName)
	return tomlStruct.TaskName, nil
}
