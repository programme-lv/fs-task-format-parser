package fstaskparser

type Task struct {
	problemTomlContent []byte

	// specificationVersion string
	// srcDirPath           string
	problemTags          []string
	problemAuthors       []string
	mdStatements         []MDStatement
	taskName             string
	originOlympiad       string
	difficultyOneToFive  int
	memoryMegabytes      int
	cpuTimeSeconds       float64
	examples             []Example
	exampleFilenameToID  map[string]int
	visibleInputSubtasks []int

	/*
		=== TESTS ===
		1) read test filenames, sort them lexiographically
		2) initialize a map from test filename to its ID
		3) initialize a map from test ID to its filename
		4) override the map with test filename-id dictionary in problem.toml
		5) override the map with test id-filename dictionary in problem.toml
		6) read tests into memory
	*/

	testFnamesSorted []string
	testFilenameToID map[string]int
	testIDOverwrite  map[string]int
	testIDToFilename map[int]string
	tests            []Test

	/*
		=== TEST GROUPS ===
		1) read all group IDs from problem.toml
		2) read wchich groups are public
		3) reach how many points each group has
		4) read to which subtask each group belongs
		5) read which test ids belong to each group
		6) read which filenames belong to each group
		7) append to groups test id(filename) for names
	*/

	testGroupIDs   []int
	isTGroupPublic map[int]bool
	tGroupPoints   map[int]int
	tGroupToStMap  map[int]int
	tGroupTestIDs  map[int][]int
	tGroupFnames   map[int][]string
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
	// ID is the order in which the file comes in lexicographical order
	// OR overriden by the filename-testID dictionary in problem.toml
	// TODO: create the filename-testID dictionary
	ID     int
	Input  []byte
	Answer []byte
}

type Example struct {
	// ID is the order in which the file comes in lexicographical order
	// OR overriden by the filename-exampleID dictionary in problem.toml
	// TODO: create the filename-exampleID dictionary
	ID     int
	Input  []byte
	Output []byte
	Name   *string
}

func NewTask(taskName string) (*Task, error) {
	t := Task{
		problemTomlContent:   []byte{},
		problemTags:          []string{},
		problemAuthors:       []string{},
		mdStatements:         []MDStatement{},
		taskName:             taskName,
		originOlympiad:       "",
		difficultyOneToFive:  0,
		memoryMegabytes:      256,
		cpuTimeSeconds:       1.0,
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

	return &t, nil
}
