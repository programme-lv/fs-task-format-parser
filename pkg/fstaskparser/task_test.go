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
		filename := parsedTask.GetTestFilenameFromID(parsedTests[i].ID)
		fnameNotPtr := ""
		if filename != nil {
			fnameNotPtr = *filename
		}
		parsedTestNames = append(parsedTestNames, fnameNotPtr)
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
		filename := parsedTask.GetTestFilenameFromID(parsedTests[i].ID)
		fnameNotPtr := ""
		if filename != nil {
			fnameNotPtr = *filename
		}
		inPath := filepath.Join(testPath, fmt.Sprintf("%s.in", fnameNotPtr))

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
		filename := parsedTask.GetTestFilenameFromID(parsedTests[i].ID)
		fnameNotPtr := ""
		if filename != nil {
			fnameNotPtr = *filename
		}
		ansPath := filepath.Join(testPath, fmt.Sprintf("%s.out", fnameNotPtr))

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
	tests := storedTask.GetTests()
	for i := 0; i < 6; i++ {
		filename := storedTask.GetTestFilenameFromID(tests[i].ID)
		filenameNotPtr := ""
		if filename != nil {
			filenameNotPtr = *filename
		}
		storedTestNames = append(storedTestNames, filenameNotPtr)
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

func TestReadingWritingExamples(t *testing.T) {
	parsedTask, err := fstaskparser.Read(testTaskPath)
	require.NoErrorf(t, err, "failed to read task: %v", err)

	parsedExamples := parsedTask.GetExamples()
	require.Equal(t, 1, len(parsedExamples))

	parsedExampleNames := []string{}
	for i := 0; i < 1; i++ {
		parsedExampleNames = append(parsedExampleNames, *parsedExamples[i].Name)
	}
	expectedExampleNames := []string{"kp00"}
	assert.Equal(t, expectedExampleNames, parsedExampleNames)

	parsedInputs := []string{}
	for i := 0; i < 1; i++ {
		parsedInputs = append(parsedInputs, string(parsedExamples[i].Input))
	}

	examplePath := filepath.Join(testTaskPath, "examples")
	expectedInputs := []string{}
	for i := 0; i < 1; i++ {
		inPath := filepath.Join(examplePath, fmt.Sprintf("%s.in", *parsedExamples[i].Name))

		in, err := os.ReadFile(inPath)
		require.NoErrorf(t, err, "failed to read input file: %v", err)

		expectedInputs = append(expectedInputs, string(in))
	}
	assert.Equal(t, expectedInputs, parsedInputs)

	parsedOutputs := []string{}
	for i := 0; i < 1; i++ {
		parsedOutputs = append(parsedOutputs, string(parsedExamples[i].Output))
	}
	expectedOutputs := []string{}
	for i := 0; i < 1; i++ {
		outPath := filepath.Join(examplePath, fmt.Sprintf("%s.out", *parsedExamples[i].Name))

		out, err := os.ReadFile(outPath)
		require.NoErrorf(t, err, "failed to read output file: %v", err)

		expectedOutputs = append(expectedOutputs, string(out))
	}

	assert.Equal(t, expectedOutputs, parsedOutputs)

	tmpDirectory, err := os.MkdirTemp("", "fstaskparser-test-")
	require.NoErrorf(t, err, "failed to create temporary directory: %v", err)
	defer os.RemoveAll(tmpDirectory)

	outputDirectory := filepath.Join(tmpDirectory, "kvadrputekl")

	t.Logf("Created directory for output: %s", outputDirectory)

	err = parsedTask.Store(outputDirectory)
	require.NoErrorf(t, err, "failed to store task: %v", err)

	storedTask, err := fstaskparser.Read(outputDirectory)
	require.NoErrorf(t, err, "failed to read task: %v", err)

	storedExampleNames := []string{}
	for i := 0; i < 1; i++ {
		storedExampleNames = append(storedExampleNames, *storedTask.GetExamples()[i].Name)
	}
	assert.Equal(t, expectedExampleNames, storedExampleNames)

	storedInputs := []string{}
	for i := 0; i < 1; i++ {
		storedInputs = append(storedInputs, string(storedTask.GetExamples()[i].Input))
	}
	assert.Equal(t, expectedInputs, storedInputs)

	storedOutputs := []string{}
	for i := 0; i < 1; i++ {
		storedOutputs = append(storedOutputs, string(storedTask.GetExamples()[i].Output))
	}
	assert.Equal(t, expectedOutputs, storedOutputs)
}

func TestReadingWritingMetadata(t *testing.T) {
	parsedTask, err := fstaskparser.Read(testTaskPath)
	require.NoErrorf(t, err, "failed to read task: %v", err)

	// Set metadata using setters
	parsedTask.SetTaskName("Kvadrātveida putekļsūcējs")
	parsedTask.SetProblemTags([]string{"math", "geometry"})
	parsedTask.SetTaskAuthors([]string{"Author1", "Author2"})
	parsedTask.SetOriginOlympiad("LIO")
	parsedTask.SetDifficultyOneToFive(3)

	// Verify the set metadata using getters
	assert.Equal(t, "Kvadrātveida putekļsūcējs", parsedTask.GetTaskName())
	assert.Equal(t, []string{"math", "geometry"}, parsedTask.GetProblemTags())
	assert.Equal(t, []string{"Author1", "Author2"}, parsedTask.GetTaskAuthors())
	assert.Equal(t, "LIO", parsedTask.GetOriginOlympiad())
	assert.Equal(t, 3, parsedTask.GetDifficultyOneToFive())

	tmpDirectory, err := os.MkdirTemp("", "fstaskparser-test-")
	require.NoErrorf(t, err, "failed to create temporary directory: %v", err)
	defer os.RemoveAll(tmpDirectory)

	outputDirectory := filepath.Join(tmpDirectory, "kvadrputekl")
	t.Logf("Created directory for output: %s", outputDirectory)

	err = parsedTask.Store(outputDirectory)
	require.NoErrorf(t, err, "failed to store task: %v", err)

	storedTask, err := fstaskparser.Read(outputDirectory)
	require.NoErrorf(t, err, "failed to read task: %v", err)

	// Verify the stored metadata using getters
	assert.Equal(t, "Kvadrātveida putekļsūcējs", storedTask.GetTaskName())
	assert.Equal(t, []string{"math", "geometry"}, storedTask.GetProblemTags())
	assert.Equal(t, []string{"Author1", "Author2"}, storedTask.GetTaskAuthors())
	assert.Equal(t, "LIO", storedTask.GetOriginOlympiad())
	assert.Equal(t, 3, storedTask.GetDifficultyOneToFive())
}

func TestReadingWritingTestGroups(t *testing.T) {
	parsedTask, err := fstaskparser.Read(testTaskPath)
	assert.NoErrorf(t, err, "failed to read task: %v", err)

	parsedTestGroups := parsedTask.GetTestGroupIDs()
	require.Equal(t, 2, len(parsedTestGroups))

	expectedTestGroups := []int{1, 2}

	assert.Equal(t, expectedTestGroups, parsedTestGroups)

	tmpDirectory, err := os.MkdirTemp("", "fstaskparser-test-")
	require.NoErrorf(t, err, "failed to create temporary directory: %v", err)
	defer os.RemoveAll(tmpDirectory)

	outputDirectory := filepath.Join(tmpDirectory, "kvadrputekl")
	t.Logf("Created directory for output: %s", outputDirectory)

	err = parsedTask.Store(outputDirectory)
	require.NoErrorf(t, err, "failed to store task: %v", err)
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
