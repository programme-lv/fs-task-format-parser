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
		tests:                []Test{},
		mdStatements:         []MDStatement{},
		taskName:             "",
		originOlympiad:       "",
		difficultyOneToFive:  0,
		memoryMegabytes:      0,
		cpuTimeSeconds:       0,
		testGroups:           []TestGroup{},
		examples:             []Example{},
		tGroupToStMap:        map[int]int{},
		isTGroupPublic:       map[int]bool{},
		tGroupPoints:         map[int]int{},
		visibleInputSubtasks: []int{},
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

	t.tests, err = readTestsDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading tests directory: %w", err)
	}

	t.examples, err = readExamplesDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading examples directory: %w", err)
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

func readCPUTimeLimitInSeconds(specVers string, tomlContent string) (float64, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.2")
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
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.2")
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

func readTestsDir(srcDirPath string) ([]Test, error) {
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

		tests = append(tests, Test{
			ID:     (i / 2) + 1,
			Input:  input,
			Answer: answer,
			Name:   &inFilenameBase,
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