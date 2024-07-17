package fstaskparser

type Task struct {
	problemTomlContent []byte
	// specificationVersion string
	// srcDirPath           string
	problemTags          []string
	problemAuthors       []string
	tests                []Test // DONE
	mdStatements         []MDStatement
	taskName             string
	originOlympiad       string
	difficultyOneToFive  int
	memoryMegabytes      int     // DONE
	cpuTimeSeconds       float64 // DONE
	testGroups           []TestGroup
	examples             []Example
	tGroupToStMap        map[int]int
	isTGroupPublic       map[int]bool
	tGroupPoints         map[int]int
	visibleInputSubtasks []int
}

func (t *Task) GetTaskName() string {
	return t.taskName
}

func (t *Task) GetProblemTags() []string {
	return t.problemTags
}

func (t *Task) GetProblemAuthors() []string {
	return t.problemAuthors
}

func (t *Task) GetMDStatements() []MDStatement {
	return t.mdStatements
}

func (t *Task) GetOriginOlympiad() string {
	return t.originOlympiad
}

func (t *Task) GetDifficultyOneToFive() int {
	return t.difficultyOneToFive
}

func (t *Task) GetTestGroups() []TestGroup {
	return t.testGroups
}

func (t *Task) GetExamples() []Example {
	return t.examples
}

func (t *Task) GetTGroupToStMap() map[int]int {
	return t.tGroupToStMap
}

func (t *Task) GetIsTGroupPublic() map[int]bool {
	return t.isTGroupPublic
}

func (t *Task) GetTGroupPoints() map[int]int {
	return t.tGroupPoints
}

func (t *Task) GetVisibleInputSubtasks() []int {
	return t.visibleInputSubtasks
}

func (t *Task) GetProblemTomlContent() []byte {
	return t.problemTomlContent
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
