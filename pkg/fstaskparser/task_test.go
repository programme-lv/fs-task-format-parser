package fstaskparser_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var prjRootPath = filepath.Join(".", "..", "..")
var testTaskPath = filepath.Join(prjRootPath, "testdata", "kvadrputekl")

func TestReadingWritingTests(t *testing.T) {
	parsedTask, err := fstaskparser.Read(testTaskPath)
	require.NoErrorf(t, err, "failed to read task: %v", err)

	parsedTests := parsedTask.GetTests()
	require.Equal(t, 6, len(parsedTests))

	parsedTestNames := []string{}
	for i := 0; i < 6; i++ {
		parsedTestNames = append(parsedTestNames, *parsedTests[i].Name)
	}
	expectedTestNames := []string{"kp01a", "kp01b", "kp01c", "kp02a", "kp02b", "kp02c"}
	assert.Equal(t, expectedTestNames, parsedTestNames)

	parsedIDs := []int{}
	for i := 0; i < 6; i++ {
		parsedIDs = append(parsedIDs, parsedTests[i].ID)
	}
	expectedIDs := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, expectedIDs, parsedIDs)

	parsedInputs := []string{}
	for i := 0; i < 6; i++ {
		parsedInputs = append(parsedInputs, string(parsedTests[i].Input))
	}

	testPath := filepath.Join(testTaskPath, "tests")
	expectedInputs := []string{}
	for i := 0; i < 6; i++ {
		inPath := filepath.Join(testPath, fmt.Sprintf("%s.in", *parsedTests[i].Name))

		in, err := os.ReadFile(inPath)
		require.NoErrorf(t, err, "failed to read input file: %v", err)

		expectedInputs = append(expectedInputs, string(in))
	}
	assert.Equal(t, expectedInputs, parsedInputs)

	parsedAnswers := []string{}
	for i := 0; i < 6; i++ {
		parsedAnswers = append(parsedAnswers, string(parsedTests[i].Answer))
	}
	expectedAnsers := []string{}
	for i := 0; i < 6; i++ {
		ansPath := filepath.Join(testPath, fmt.Sprintf("%s.out", *parsedTests[i].Name))

		ans, err := os.ReadFile(ansPath)
		require.NoErrorf(t, err, "failed to read answer file: %v", err)

		expectedAnsers = append(expectedAnsers, string(ans))
	}

	assert.Equal(t, expectedAnsers, parsedAnswers)

	tmpDirectory, err := os.MkdirTemp("", "fstaskparser-test-")
	require.NoErrorf(t, err, "failed to create temporary directory: %v", err)
	defer os.RemoveAll(tmpDirectory)

	outputDirectory := filepath.Join(tmpDirectory, "kvadrputekl")

	t.Logf("Created directory for output: %s", outputDirectory)

	err = parsedTask.Store(outputDirectory)
	require.NoErrorf(t, err, "failed to store task: %v", err)

	storedTask, err := fstaskparser.Read(outputDirectory)
	require.NoErrorf(t, err, "failed to read task: %v", err)

	storedTestNames := []string{}
	for i := 0; i < 6; i++ {
		storedTestNames = append(storedTestNames, *storedTask.GetTests()[i].Name)
	}
	assert.Equal(t, expectedTestNames, storedTestNames)

	storedIDs := []int{}
	for i := 0; i < 6; i++ {
		storedIDs = append(storedIDs, storedTask.GetTests()[i].ID)
	}
	assert.Equal(t, expectedIDs, storedIDs)

	storedInputs := []string{}
	for i := 0; i < 6; i++ {
		storedInputs = append(storedInputs, string(storedTask.GetTests()[i].Input))
	}
	assert.Equal(t, expectedInputs, storedInputs)

	storedAnswers := []string{}
	for i := 0; i < 6; i++ {
		storedAnswers = append(storedAnswers, string(storedTask.GetTests()[i].Answer))
	}
	assert.Equal(t, expectedAnsers, storedAnswers)
}

func TestReadingWritingEvaluationConstraints(t *testing.T) {
	parsedTask, err := fstaskparser.Read(testTaskPath)
	require.NoErrorf(t, err, "failed to read task: %v", err)
	assert.Equal(t, 0.5, parsedTask.GetCPUTimeLimitInSeconds())
	assert.Equal(t, 256, parsedTask.GetMemoryLimitInMegabytes())

	tmpDirectory, err := os.MkdirTemp("", "fstaskparser-test-")
	require.NoErrorf(t, err, "failed to create temporary directory: %v", err)
	defer os.RemoveAll(tmpDirectory)

	outputDirectory := filepath.Join(tmpDirectory, "kvadrputekl")
	t.Logf("Created directory for output: %s", outputDirectory)

	err = parsedTask.Store(outputDirectory)
	require.NoErrorf(t, err, "failed to store task: %v", err)

	storedTask, err := fstaskparser.Read(outputDirectory)
	require.NoErrorf(t, err, "failed to read task: %v", err)
	assert.Equal(t, 0.5, storedTask.GetCPUTimeLimitInSeconds())
	assert.Equal(t, 256, storedTask.GetMemoryLimitInMegabytes())
}

// func TestStoringComplexTask(t *testing.T) {
// 	testDir := filepath.Join(".", "..", "..", "testdata", "kvadrputekl")
// 	task, err := fstaskparser.Read(testDir)
// 	if err != nil {
// 		t.Fatalf("failed to read task: %v", err)
// 	}

// 	originalCPUTimeLimitInSeconds := task.GetCPUTimeLimitInSeconds()
// 	assert.Equal(t, 0.5, originalCPUTimeLimitInSeconds)

// 	originalMemoryLimitInMegabytes := task.GetMemoryLimitInMegabytes()
// 	assert.Equal(t, 256, originalMemoryLimitInMegabytes)

// 	// Compare with real values on disk
// 	problemTomlPath := filepath.Join(testDir, "problem.toml")
// 	problemTomlContent, err := os.ReadFile(problemTomlPath)
// 	require.NoErrorf(t, err, "failed to read problem.toml: %v", err)
// 	assert.Equal(t, string(task.GetProblemTomlContent()), string(problemTomlContent))

// 	// Compare examples directory
// 	examplesDir := filepath.Join(testDir, "examples")
// 	exampleFiles, err := os.ReadDir(examplesDir)
// 	require.NoErrorf(t, err, "failed to read examples directory: %v", err)

// 	for _, exampleFile := range exampleFiles {
// 		exampleFilePath := filepath.Join(examplesDir, exampleFile.Name())
// 		exampleFileContent, err := os.ReadFile(exampleFilePath)
// 		require.NoErrorf(t, err, "failed to read example file: %v", err)

// 		var exampleContent []byte
// 		for _, example := range task.GetExamples() {
// 			if exampleFile.Name() == *example.Name {
// 				exampleContent = example.Input
// 				break
// 			}
// 		}
// 		assert.Equal(t, string(exampleContent), string(exampleFileContent))
// 	}
// 	tmpDirectory, err := os.MkdirTemp("", "fstaskparser-test-")
// 	if err != nil {
// 		t.Fatalf("failed to create temporary directory: %v", err)
// 	}
// 	defer os.RemoveAll(tmpDirectory)

// 	outputDirectory := filepath.Join(tmpDirectory, "kvadrputekl")

// 	t.Logf("Created directory for output: %s", outputDirectory)

// 	err = task.Store(outputDirectory)
// 	if err != nil {
// 		t.Fatalf("failed to store task: %v", err)
// 	}

// 	task2, err := fstaskparser.Read(outputDirectory)
// 	require.NoErrorf(t, err, "failed to read task: %v", err)

// 	writtenCPUTimeLimitInSeconds := task2.GetCPUTimeLimitInSeconds()
// 	assert.Equal(t, 0.5, writtenCPUTimeLimitInSeconds)

// 	writtenMemoryLimitInMegabytes := task2.GetMemoryLimitInMegabytes()
// 	assert.Equal(t, originalMemoryLimitInMegabytes, writtenMemoryLimitInMegabytes)

// 	// Compare examples directory
// 	examplesDir2 := filepath.Join(outputDirectory, "examples")
// 	exampleFiles2, err := os.ReadDir(examplesDir2)
// 	require.NoErrorf(t, err, "failed to read examples directory: %v", err)

// 	for _, exampleFile := range exampleFiles2 {
// 		exampleFilePath := filepath.Join(examplesDir2, exampleFile.Name())
// 		exampleFileContent, err := os.ReadFile(exampleFilePath)
// 		require.NoErrorf(t, err, "failed to read example file: %v", err)

// 		var exampleContent []byte
// 		for _, example := range task2.GetExamples() {
// 			if exampleFile.Name() == *example.Name {
// 				exampleContent = example.Input
// 				break
// 			}
// 		}
// 		assert.Equal(t, string(exampleContent), string(exampleFileContent))
// 	}

// 	// TODO: compare all the values to originals
// }

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
