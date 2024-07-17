package fstaskparser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func Read(dirPath string) (*Task, error) {
	t := Task{
		problemTomlContent:   []byte{},
		problemTags:          []string{},
		problemAuthors:       []string{},
		mdStatements:         []MDStatement{},
		taskName:             "",
		originOlympiad:       "",
		difficultyOneToFive:  0,
		memoryMegabytes:      0,
		cpuTimeSeconds:       0,
		examples:             []Example{},
		exampleFilenameToID:  map[string]int{},
		visibleInputSubtasks: []int{},
		testGroupIDs:         []int{},
		isTGroupPublic:       map[int]bool{},
		tGroupPoints:         map[int]int{},
		tGroupToStMap:        map[int]int{},
		testFnamesSorted:     []string{},
		testFilenameToID:     map[string]int{},
		testIDOverwrite:      map[string]int{},
		testIDToFilename:     map[int]string{},
		tests:                []Test{},
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

	return &t, nil
}

func readTestIDOverwrite(specVers string, tomlContent []byte) (map[string]int, error) {
	semVerCmpRes, err := getCmpSemVersionsResult(specVers, "v2.3.0")
	if err != nil {
		return nil, fmt.Errorf("error comparing sem versions: %w", err)
	}

	if semVerCmpRes < 0 {
		log.Println("warning: unsupported specification version (too old):", proglvFSTaskFormatSpecVersOfScript)
		// return empty map
		return make(map[string]int), nil
	}

	tomlStruct := struct {
		TestIDOverwrite map[string]int `toml:"test_id_overwrite"`
	}{}

	err = toml.Unmarshal(tomlContent, &tomlStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the test id overwrite: %w", err)
	}

	return tomlStruct.TestIDOverwrite, nil
}

func readTestFNamesSorted(dirPath string) ([]string, error) {
	fnames, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading test filenames: %w", err)
	}

	sort.Slice(fnames, func(i, j int) bool {
		return fnames[i].Name() < fnames[j].Name()
	})

	if len(fnames)%2 != 0 {
		return nil, fmt.Errorf("odd number of test filenames")
	}

	res := make([]string, 0, len(fnames)/2)
	for i := 0; i < len(fnames); i += 2 {
		a_name := fnames[i].Name()
		// remove extension
		a_name = a_name[:len(a_name)-len(filepath.Ext(a_name))]

		b_name := fnames[i+1].Name()
		// remove extension
		b_name = b_name[:len(b_name)-len(filepath.Ext(b_name))]

		if a_name != b_name {
			return nil, fmt.Errorf("input and answer file base names do not match: %s, %s", a_name, b_name)
		}

		res = append(res, a_name)
	}

	return res, nil
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

func readCPUTimeLimitInSeconds(specVers string, tomlContent string) (float64, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return 0, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type constraintsStruct struct {
		CPUTimeLimitInSeconds float64 `toml:"cpu_time_seconds"`
	}

	tomlStruct := struct {
		Constraints constraintsStruct `toml:"constraints"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal the cpu time limit: %w", err)
	}

	return tomlStruct.Constraints.CPUTimeLimitInSeconds, nil
}

func readMemoryLimitInMegabytes(specVers string, tomlContent string) (int, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return 0, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type constraintsStruct struct {
		MemoryLimitInMegabytes int `toml:"memory_megabytes"`
	}

	tomlStruct := struct {
		Constraints constraintsStruct `toml:"constraints"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal the memory limit: %w", err)
	}

	return tomlStruct.Constraints.MemoryLimitInMegabytes, nil
}

func readProblemTags(specVers string, tomlContent string) ([]string, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return nil, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return nil, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type metadataStruct struct {
		ProblemTags []string `toml:"problem_tags"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	return tomlStruct.Metadata.ProblemTags, nil
}

func readProblemAuthors(specVers string, tomlContent string) ([]string, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return nil, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return nil, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type metadataStruct struct {
		ProblemAuthors []string `toml:"task_authors"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	return tomlStruct.Metadata.ProblemAuthors, nil
}

func readOriginOlympiad(specVers string, tomlContent string) (string, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return "", fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return "", fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type metadataStruct struct {
		OriginOlympiad *string `toml:"origin_olympiad"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	res := ""
	if tomlStruct.Metadata.OriginOlympiad != nil {
		res = *tomlStruct.Metadata.OriginOlympiad
	}
	return res, nil
}

func readDifficultyOneToFive(specVers string, tomlContent string) (int, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return 0, fmt.Errorf("unsupported specification version: %s", specVers)
	}
	type metadataStruct struct {
		DifficultyFrom1To5 *int `toml:"difficulty_1_to_5"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	res := 0
	if tomlStruct.Metadata.DifficultyFrom1To5 != nil {
		res = *tomlStruct.Metadata.DifficultyFrom1To5
	}

	return res, nil
}

func readTestsDir(srcDirPath string, fnameToID map[string]int) ([]Test, error) {
	dir := filepath.Join(srcDirPath, "tests")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading tests directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	tests := make([]Test, 0, len(entries)/2)

	for i := 0; i < len(entries); i += 2 {
		inPath := filepath.Join(dir, entries[i].Name())
		ansPath := filepath.Join(dir, entries[i+1].Name())

		inFilename := entries[i].Name()
		ansFilename := entries[i+1].Name()

		inFilenameBase := strings.TrimSuffix(inFilename, filepath.Ext(inFilename))
		ansFilenameBase := strings.TrimSuffix(ansFilename, filepath.Ext(ansFilename))

		if inFilenameBase != ansFilenameBase {
			return nil, fmt.Errorf("input and answer file base names do not match: %s, %s", inFilenameBase, ansFilenameBase)
		}

		// sometimes the test answer is stored as .out, sometimes as .ans
		if strings.Contains(inFilename, ".ans") || strings.Contains(ansFilename, ".in") {
			// swap the file paths
			inPath, ansPath = ansPath, inPath
		}

		input, err := os.ReadFile(inPath)
		if err != nil {
			return nil, fmt.Errorf("error reading input file: %w", err)
		}

		answer, err := os.ReadFile(ansPath)
		if err != nil {
			return nil, fmt.Errorf("error reading answer file: %w", err)
		}

		// check if mapping to id exists
		if _, ok := fnameToID[inFilenameBase]; !ok {
			return nil, fmt.Errorf("mapping from filename to id does not exist: %s", inFilenameBase)
		}

		tests = append(tests, Test{
			ID:     fnameToID[inFilenameBase],
			Input:  input,
			Answer: answer,
		})
	}

	return tests, nil
}

func readExamplesDir(srcDirPath string) ([]Example, error) {
	dir := filepath.Join(srcDirPath, "examples")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading examples directory: %w", err)
	}
	// tests are to be read exactly like examples

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	examples := make([]Example, 0, len(entries)/2)

	for i := 0; i < len(entries); i += 2 {
		inPath := filepath.Join(dir, entries[i].Name())
		ansPath := filepath.Join(dir, entries[i+1].Name())

		inFilename := entries[i].Name()
		ansFilename := entries[i+1].Name()

		inFilenameBase := strings.TrimSuffix(inFilename, filepath.Ext(inFilename))
		ansFilenameBase := strings.TrimSuffix(ansFilename, filepath.Ext(ansFilename))

		if inFilenameBase != ansFilenameBase {
			return nil, fmt.Errorf("input and answer file base names do not match: %s, %s", inFilenameBase, ansFilenameBase)
		}

		// sometimes the test answer is stored as .out, sometimes as .ans
		if strings.Contains(inFilename, ".ans") || strings.Contains(ansFilename, ".in") {
			// swap the file paths
			inPath, ansPath = ansPath, inPath
		}

		input, err := os.ReadFile(inPath)
		if err != nil {
			return nil, fmt.Errorf("error reading input file: %w", err)
		}

		answer, err := os.ReadFile(ansPath)
		if err != nil {
			return nil, fmt.Errorf("error reading answer file: %w", err)
		}

		examples = append(examples, Example{
			ID:     (i / 2) + 1,
			Input:  input,
			Output: answer,
			Name:   &inFilenameBase,
		})
	}

	return examples, nil
}
