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
	originOlympiad       *string
	difficultyOneToFive  *int
	memoryMegabytes      int
	cpuTimeSeconds       float64
	testGroups           []TestGroup
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
	ID     int
	Input  []byte
	Answer []byte
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
		originOlympiad:       new(string),
		difficultyOneToFive:  new(int),
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
