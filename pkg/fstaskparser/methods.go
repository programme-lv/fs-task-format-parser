package fstaskparser

import "log"

func (t *task) GetCPUTimeLimitInSeconds() float64 {
	return t.cpuTimeSeconds
}

func (t *task) GetMemoryLimitInMegabytes() int {
	return t.memoryMegabytes
}

func (t *task) GetTests() []test {
	return t.tests
}

// creates a new test and returns its ID
func (t *task) AddTest(input []byte, answer []byte) int {
	// find the minimum positive excluded id from tests

	// we assign this test the number but it may not correspond to lex order
	// that is the responsibility of the persistence layer

	mex := 1
	found := true
	for found {
		found = false
		for i := 0; i < len(t.tests); i++ {
			if t.tests[i].ID == mex {
				found = true
				mex++
			}
		}
	}

	t.tests = append(t.tests, test{
		ID:     mex,
		Input:  input,
		Answer: answer,
	})

	return mex
}

func (t *task) AssignFilenameToTest(filename string, testID int) {
	_, ok1 := t.testIDToFilename[testID]
	_, ok2 := t.testFilenameToID[filename]
	if ok1 || ok2 {
		log.Fatalf("test with ID %d or filename %s already exists", testID, filename)
	}

	t.testIDToFilename[testID] = filename
	t.testFilenameToID[filename] = testID
}

func (t *task) GetExamples() []example {
	return t.examples
}

func (t *task) GetTaskName() string {
	return t.taskName
}

func (t *task) SetTaskName(name string) {
	t.taskName = name
}

func (t *task) GetProblemTags() []string {
	return t.problemTags
}

func (t *task) SetProblemTags(tags []string) {
	t.problemTags = tags
}

func (t *task) GetTaskAuthors() []string {
	return t.problemAuthors
}

func (t *task) SetTaskAuthors(authors []string) {
	t.problemAuthors = authors
}

func (t *task) GetOriginOlympiad() string {
	return t.originOlympiad
}

func (t *task) SetOriginOlympiad(origin string) {
	t.originOlympiad = origin
}

func (t *task) GetDifficultyOneToFive() int {
	return t.difficultyOneToFive
}

func (t *task) SetDifficultyOneToFive(difficulty int) {
	t.difficultyOneToFive = difficulty
}

type TestGroupInfo struct {
	GroupID int
	Points  int
	Public  bool
	TestIDs []int
	Subtask int
}

func (t *task) GetInfoOnTestGroup(id int) TestGroupInfo {
	return TestGroupInfo{
		GroupID: id,
		Points:  t.tGroupPoints[id],
		Public:  t.isTGroupPublic[id],
		TestIDs: t.tGroupTestIDs[id],
		Subtask: t.tGroupToStMap[id],
	}
}

func (t *task) GetTestGroupIDs() []int {
	return t.testGroupIDs
}

func (t *task) GetTestFilenameFromID(testID int) string {
	filename, ok := t.testIDToFilename[testID]
	if !ok {
		return ""
	}
	return filename
}
