package fstaskparser

type Task struct {
	problemTomlContent []byte
	// specificationVersion string
	// srcDirPath           string
	problemTags          []string
	problemAuthors       []string
	tests                []Test
	mdStatements         []MDStatement
	taskName             string
	originOlympiad       string
	difficultyOneToFive  int
	memoryMegabytes      int
	cpuTimeSeconds       float64
	testGroups           []TestGroup
	examples             []Example
	tGroupToStMap        map[int]int
	isTGroupPublic       map[int]bool
	tGroupPoints         map[int]int
	visibleInputSubtasks []int
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
	Name   *string
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
		tests:                []Test{},
		mdStatements:         []MDStatement{},
		taskName:             taskName,
		originOlympiad:       "",
		difficultyOneToFive:  0,
		memoryMegabytes:      256,
		cpuTimeSeconds:       1.0,
		testGroups:           []TestGroup{},
		tGroupToStMap:        map[int]int{},
		isTGroupPublic:       map[int]bool{},
		tGroupPoints:         map[int]int{},
		visibleInputSubtasks: []int{},
	}

	return &t, nil
}

func (t *Task) GetTests() []Test {
	return t.tests
}
