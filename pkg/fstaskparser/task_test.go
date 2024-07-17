package fstaskparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoringComplexTask(t *testing.T) {
	testDir := filepath.Join(".", "..", "..", "testdata", "kvadrputekl")
	task, err := fstaskparser.Read(testDir)
	if err != nil {
		t.Fatalf("failed to read task: %v", err)
	}

	originalCPUTimeLimitInSeconds, err := task.GetCPUTimeLimitInSeconds()
	require.NoErrorf(t, err, "failed to get cpu time limit: %v", err)
	assert.Equal(t, 0.5, originalCPUTimeLimitInSeconds)

	originalMemoryLimitInMegabytes, err := task.GetMemoryLimitInMegabytes()
	require.NoErrorf(t, err, "failed to get memory limit: %v", err)
	assert.Equal(t, 256, originalMemoryLimitInMegabytes)

	// Create a temporary directory for output
	tmpDirectory, err := os.MkdirTemp("", "fstaskparser-test-")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDirectory)

	outputDirectory := filepath.Join(tmpDirectory, "kvadrputekl")

	t.Logf("Created directory for output: %s", outputDirectory)

	err = task.Store(outputDirectory)
	if err != nil {
		t.Fatalf("failed to store task: %v", err)
	}

	task2, err := fstaskparser.Read(outputDirectory)
	require.NoErrorf(t, err, "failed to read task: %v", err)

	writtenCPUTimeLimitInSeconds, err := task2.GetCPUTimeLimitInSeconds()
	require.NoErrorf(t, err, "failed to get cpu time limit: %v", err)

	assert.Equal(t, 0.5, writtenCPUTimeLimitInSeconds)

	writtenMemoryLimitInMegabytes, err := task2.GetMemoryLimitInMegabytes()
	require.NoErrorf(t, err, "failed to get memory limit: %v", err)

	assert.Equal(t, originalMemoryLimitInMegabytes, writtenMemoryLimitInMegabytes)

	// the tests after writing should have 12 files
	// two of those files should be kp02a.in and kp02a.out

	files, err := os.ReadDir(filepath.Join(outputDirectory, "tests"))
	if err != nil {
		t.Fatalf("failed to read tests directory: %v", err)
	}

	assert.Equal(t, 12, len(files))

	foundKp02aIn := false
	foundKp02aOut := false
	for _, f := range files {
		if f.Name() == "kp02a.in" {
			foundKp02aIn = true
		}
		if f.Name() == "kp02a.out" {
			foundKp02aOut = true
		}
	}

	assert.True(t, foundKp02aIn)
	assert.True(t, foundKp02aOut)
	// Verify problem.toml
	problemTomlPath := filepath.Join(outputDirectory, "problem.toml")
	problemTomlContent, err := os.ReadFile(problemTomlPath)
	require.NoErrorf(t, err, "failed to read problem.toml: %v", err)

	expectedProblemTomlContent := `specification = '2.2'
task_name = 'Kvadrātveida putekļsūcējs'
visible_input_subtasks = [1]

[metadata]
  problem_tags = []
  difficulty_1_to_5 = 3
  task_authors = []
  origin_olympiad = 'LIO'

[constraints]
  memory_megabytes = 256
  cpu_time_seconds = 0.5

[[test_groups]]
  group_id = 1
  points = 3
  subtask = 1
  public = true
  test_filenames = ['kp01a.in', 'kp01b.in', 'kp01c.in', 'kp01a.out', 'kp01b.out', 'kp01c.out']

[[test_groups]]
  group_id = 2
  points = 8
  subtask = 2
  public = true
  test_filenames = ['kp02a.in', 'kp02b.in', 'kp02c.in', 'kp02a.out', 'kp02b.out', 'kp02c.out']
`
	assert.Equal(t, expectedProblemTomlContent, string(problemTomlContent))

	// Verify examples directory
	examplesDir := filepath.Join(outputDirectory, "examples")
	exampleFiles, err := os.ReadDir(examplesDir)
	require.NoErrorf(t, err, "failed to read examples directory: %v", err)

	assert.Equal(t, 2, len(exampleFiles))

	foundExampleIn := false
	foundExampleOut := false
	for _, f := range exampleFiles {
		if f.Name() == "kp00.in" {
			foundExampleIn = true
		}
		if f.Name() == "kp00.out" {
			foundExampleOut = true
		}
	}

	assert.True(t, foundExampleIn)
	assert.True(t, foundExampleOut)
	// Verify problem.toml using methods
	expectedTask := &fstaskparser.Task{
		taskName:             "Kvadrātveida putekļsūcējs",
		problemTags:          []string{},
		difficultyOneToFive:  3,
		problemAuthors:       []string{},
		originOlympiad:       "LIO",
		memoryMegabytes:      256,
		cpuTimeSeconds:       0.5,
		visibleInputSubtasks: []int{1},
		testGroups: []fstaskparser.TestGroup{
			{
				GroupID: 1,
				TestIDs: []int{1, 2, 3, 4, 5, 6},
			},
			{
				GroupID: 2,
				TestIDs: []int{7, 8, 9, 10, 11, 12},
			},
		},
	}

	assert.Equal(t, expectedTask.taskName, task2.taskName)
	assert.Equal(t, expectedTask.problemTags, task2.problemTags)
	assert.Equal(t, expectedTask.difficultyOneToFive, task2.difficultyOneToFive)
	assert.Equal(t, expectedTask.problemAuthors, task2.problemAuthors)
	assert.Equal(t, expectedTask.originOlympiad, task2.originOlympiad)
	assert.Equal(t, expectedTask.memoryMegabytes, task2.memoryMegabytes)
	assert.Equal(t, expectedTask.cpuTimeSeconds, task2.cpuTimeSeconds)
	assert.Equal(t, expectedTask.visibleInputSubtasks, task2.visibleInputSubtasks)
	assert.Equal(t, expectedTask.testGroups, task2.testGroups)

	// Verify examples directory using methods
	examples := task2.GetExamples()
	assert.Equal(t, 2, len(examples))

	foundExampleIn := false
	foundExampleOut := false
	for _, example := range examples {
		if example.Name != nil && *example.Name == "kp00.in" {
			foundExampleIn = true
		}
		if example.Name != nil && *example.Name == "kp00.out" {
			foundExampleOut = true
		}
	}

	assert.True(t, foundExampleIn)
	assert.True(t, foundExampleOut)
}

/*
kvadrputekl problem.toml

specification = '2.2'
task_name = 'Kvadrātveida putekļsūcējs'
visible_input_subtasks = [1]

[metadata]
  problem_tags = []
  difficulty_1_to_5 = 3
  task_authors = []
  origin_olympiad = 'LIO'

[constraints]
  memory_megabytes = 256
  cpu_time_seconds = 0.5

[[test_groups]]
  group_id = 1
  points = 3
  subtask = 1
  public = true
  test_filenames = ['kp01a.in', 'kp01b.in', 'kp01c.in', 'kp01a.out', 'kp01b.out', 'kp01c.out']

[[test_groups]]
  group_id = 2
  points = 8
  subtask = 2
  public = true
  test_filenames = ['kp02a.in', 'kp02b.in', 'kp02c.in', 'kp02a.out', 'kp02b.out', 'kp02c.out']
*/

/*
/home/kp/Programming/PROGLV/fs-task-format-parser > tree ./testdata/kvadrputekl/
./testdata/kvadrputekl/
├── examples
│   ├── kp00.in
│   └── kp00.out
├── problem.toml
├── statements
│   └── pdf
│       └── lv.pdf
└── tests
    ├── kp01a.in
    ├── kp01a.out
    ├── kp01b.in
    ├── kp01b.out
    ├── kp01c.in
    ├── kp01c.out
    ├── kp02a.in
    ├── kp02a.out
    ├── kp02b.in
    ├── kp02b.out
    ├── kp02c.in
    └── kp02c.out

5 directories, 16 files
*/
