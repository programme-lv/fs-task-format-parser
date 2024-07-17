package fstaskparser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	if specVers[0] == 'v' {
		specVers = specVers[1:]
	}

	t.taskName, err = readTaskName(specVers, string(problemTomlContent))
	if err != nil {
		return nil, fmt.Errorf("error reading task name: %w", err)
	}

	t.tests, err = readTestsDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading tests directory: %w", err)
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

func largerOrEqualSemVersionThan(a, b string) (bool, error) {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < len(aParts) || i < len(bParts); i++ {
		if i >= len(aParts) { // a is shorter and equal in the common part
			return false, nil
		}
		if i >= len(bParts) { // b is shorter and equal in the common part
			return true, nil
		}
		// cast to ints
		a_i, err := strconv.Atoi(aParts[i])
		if err != nil {
			return false, fmt.Errorf("error converting version part to int: %w", err)
		}
		b_i, err := strconv.Atoi(bParts[i])
		if err != nil {
			return false, fmt.Errorf("error converting version part to int: %w", err)
		}
		if a_i < b_i {
			return false, nil
		}
		if a_i > b_i {
			return true, nil
		}
	}
	return true, nil
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
