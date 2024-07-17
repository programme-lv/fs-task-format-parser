package fstaskparser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pelletier/go-toml/v2"
)

const proglvFSTaskFormatSpecVersion = "2.3"

type Task struct {
	problemTomlContent   []byte
	specificationVersion string
	srcDirPath           string
	problemTags          []string
	problemAuthors       []string
	tests                []Test
	mdStatements         []MDStatement
	taskName             string
	originOlympiad       *string
	difficultyOneToFive  *int
	memoryMegabytes      int
	cpuTimeSeconds       float64
	testGroups           []TestGroup
	tGroupToStMap        map[int]int
	isTGroupPublic       map[int]bool
	tGroupPoints         map[int]int
	visibleInputSubtasks []int
	testScoringType      string // "individual", "group", "subtask"
}

type TestGroup struct {
	GroupID int
	TestIDs []int
}

type MDStatement struct {
	Language *string
	Story    string
	Input    string
	Output   string
	Notes    *string
	Scoring  *string
}

// tests are executed in order of ID
type Test struct {
	ID     int
	Input  []byte
	Answer []byte
	Name   *string
}

func (t *Task) Read(dirPath string) error {
	problemTomlPath := filepath.Join(dirPath, "problem.toml")
	problemTomlContent, err := os.ReadFile(problemTomlPath)
	if err != nil {
		return fmt.Errorf("error reading problem.toml: %w", err)
	}

	t.problemTomlContent = problemTomlContent

	var specVersStruct struct {
		Specification string `toml:"specification"`
	}

	err = toml.Unmarshal(problemTomlContent, &specVersStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the specification: %w", err)
	}

	t.specificationVersion = specVersStruct.Specification

	err = t.readTestsDir()
	if err != nil {
		return fmt.Errorf("error reading tests directory: %w", err)
	}

	return nil
}

func (t *Task) readTestsDir() error {
	dir := filepath.Join(t.srcDirPath, "tests")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading tests directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i := 0; i < len(entries); i += 2 {
		inPath := filepath.Join(dir, entries[i].Name())
		ansPath := filepath.Join(dir, entries[i+1].Name())

		input, err := os.ReadFile(inPath)
		if err != nil {
			return fmt.Errorf("error reading input file: %w", err)
		}

		answer, err := os.ReadFile(ansPath)
		if err != nil {
			return fmt.Errorf("error reading answer file: %w", err)
		}

		filename := entries[i].Name()

		t.tests = append(t.tests, Test{
			ID:     (i / 2) + 1,
			Input:  input,
			Answer: answer,
			Name:   &filename,
		})
	}

	return nil
}

func (t *Task) GetTests() []Test {
	return t.tests
}

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
