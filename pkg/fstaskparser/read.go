package fstaskparser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func Read(dirPath string) (*task, error) {
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
	problemTomlContent, err := os.ReadFile(problemTomlPath)
	if err != nil {
		return nil, fmt.Errorf("error reading problem.toml: %w", err)
	}

	t.problemTomlContent = problemTomlContent

	var specVersStruct struct {
		Specification string `toml:"specification"`
	}

	err = toml.Unmarshal(problemTomlContent, &specVersStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the specification: %w", err)
	}

	specVers := specVersStruct.Specification
	if len(specVers) == 0 {
		return nil, fmt.Errorf("empty specification")
	}
	if specVers[0] == 'v' {
		specVers = specVers[1:]
	}

	log.Println("specification version:", specVers)

	semVersCmpRes, err := getCmpSemVersionsResult(specVers, proglvFSTaskFormatSpecVersOfScript)
	if err != nil {
		return nil, fmt.Errorf("error comparing sem versions: %w", err)
	}

	// if the semantic version is larger we have a problem
	if semVersCmpRes > 0 {
		return nil, fmt.Errorf("unsupported specification version (too new): %s", specVers)
	}

	if semVersCmpRes < 0 {
		log.Println("warning: unsupported specification version (too old):", specVers)
	}

	t.taskName, err = readTaskName(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading task name: %w", err)
	}

	t.cpuTimeSeconds, err = readCPUTimeLimitInSeconds(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading cpu time limit: %w", err)
	}

	t.memoryMegabytes, err = readMemoryLimitInMegabytes(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading memory limit: %w", err)
	}

	t.problemTags, err = readProblemTags(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading problem tags: %w", err)
	}

	t.problemAuthors, err = readProblemAuthors(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading problem authors: %w", err)
	}

	t.originOlympiad, err = readOriginOlympiad(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading origin olympiad: %w", err)
	}

	t.difficultyOneToFive, err = readDifficultyOneToFive(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading difficulty: %w", err)
	}

	// read test filenames into the slice
	t.testFnamesSorted, err = readTestFNamesSorted(filepath.Join(dirPath, "tests"))
	if err != nil {
		return nil, fmt.Errorf("error reading test filenames: %w", err)
	}

	for i, fname := range t.testFnamesSorted {
		t.testFilenameToID[fname] = i + 1
		t.testIDToFilename[i+1] = fname
	}

	t.testIDOverwrite, err = readTestIDOverwrite(specVers, problemTomlContent)
	if err != nil {
		return nil, fmt.Errorf("error reading test id overwrite: %w", err)
	}

	for k, v := range t.testIDOverwrite {
		t.testIDToFilename[v] = k
		t.testFilenameToID[k] = v
	}

	// iterate over test filenames, make sure no two filenames have the same ID
	// iterate over IDS, make sure no two IDS have the same filename

	spottedFnames := make(map[int]bool) // stores spotted ids when reading filenames
	for _, fname := range t.testFnamesSorted {
		if _, ok := spottedFnames[t.testFilenameToID[fname]]; ok {
			return nil, fmt.Errorf("duplicate filename for ID: %s", fname)
		}
		spottedFnames[t.testFilenameToID[fname]] = true
	}

	spottedIDs := make(map[string]bool) // stores spotted fnames when reading ids
	for _, id := range t.testIDToFilename {
		if _, ok := spottedIDs[id]; ok {
			return nil, fmt.Errorf("duplicate ID for filename: %s", id)
		}
		spottedIDs[id] = true
	}

	t.tests, err = readTestsDir(dirPath, t.testFilenameToID)
	if err != nil {
		return nil, fmt.Errorf("error reading tests directory: %w", err)
	}

	t.examples, err = readExamplesDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading examples directory: %w", err)
	}

	t.testGroupIDs, err = readTestGroupIDs(specVers, problemTomlContent)
	if err != nil {
		return nil, fmt.Errorf("error reading test group IDs: %w", err)
	}

	t.isTGroupPublic, err = readIsTGroupPublic(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		return nil, fmt.Errorf("error reading is test group public: %w", err)
	}

	t.tGroupPoints, err = readTGroupPoints(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		return nil, fmt.Errorf("error reading test group points: %w", err)
	}

	t.tGroupToStMap, err = readTGroupToStMap(specVers, problemTomlContent)
	if err != nil {
		return nil, fmt.Errorf("error reading test group to subtask map: %w", err)
	}

	t.tGroupTestIDs, err = readTGroupTestIDs(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		return nil, fmt.Errorf("error reading test group test IDs: %w", err)
	}

	t.tGroupFnames, err = readTGroupFnames(specVers, problemTomlContent, t.testGroupIDs)
	if err != nil {
		return nil, fmt.Errorf("error reading test group filenames: %w", err)
	}

	// add to ids
	for k, v := range t.tGroupFnames {
		for _, fname := range v {
			t.tGroupTestIDs[k] = append(t.tGroupTestIDs[k], t.testFilenameToID[fname])
		}
	}

	// validate that no two test groups have the same test ID
	idsSpotted := make(map[int]bool)

	for _, v := range t.testGroupIDs {
		for _, id := range t.tGroupTestIDs[v] {
			if _, ok := idsSpotted[id]; ok {
				return nil, fmt.Errorf("duplicate test ID in test group: %d", id)
			}
			idsSpotted[id] = true
		}
	}

	return &t, nil
}

func readTaskName(specVers string, tomlContent string) (string, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.2")
	if err != nil {
		return "", fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return "", fmt.Errorf("unsupported specification version: %s", specVers)
	}

	tomlStruct := struct {
		TaskName string `toml:"task_name"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal the task name: %w", err)
	}

	return tomlStruct.TaskName, nil
}
