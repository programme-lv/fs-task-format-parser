package fstaskparser

import (
	"fmt"
	"log"
	"sort"
)

func (t *Task) GetCPUTimeLimitInSeconds() float64 {
	return t.cpuTimeSeconds
}

func (t *Task) SetCPUTimeLimitInSeconds(seconds float64) {
	t.cpuTimeSeconds = seconds
}

func (t *Task) GetMemoryLimitInMegabytes() int {
	return t.memoryMegabytes
}

func (t *Task) SetMemoryLimitInMegabytes(megabytes int) {
	t.memoryMegabytes = megabytes
}

func (t *Task) SwapTestsWithIDs(id1 int, id2 int) {
	id1Filename, id1HasFilename := t.testIDToFilename[id1]
	id2Filename, id2HasFilename := t.testIDToFilename[id2]

	for i := 0; i < len(t.tests); i++ {
		if t.tests[i].ID == id1 {
			t.tests[i].ID = id2
		} else if t.tests[i].ID == id2 {
			t.tests[i].ID = id1
		}
	}
	// testIDToFilename is affected,
	// testFilenameToID is affected,
	// testIDoverwrite is not affected
	// tGroupTestIDs is affected

	delete(t.testFilenameToID, id1Filename)
	delete(t.testIDToFilename, id1)
	delete(t.testFilenameToID, id2Filename)
	delete(t.testIDToFilename, id2)

	if id1HasFilename {
		t.testFilenameToID[id1Filename] = id2
		t.testIDToFilename[id2] = id1Filename
	}
	if id2HasFilename {
		t.testFilenameToID[id2Filename] = id1
		t.testIDToFilename[id1] = id2Filename
	}

	for k, v := range t.tGroupTestIDs {
		for i := 0; i < len(v); i++ {
			if v[i] == id1 {
				v[i] = id2
			} else if v[i] == id2 {
				v[i] = id1
			}
		}
		t.tGroupTestIDs[k] = v
	}
}

func (t *Task) GetTestsSortedByID() []test {
	sort.Slice(t.tests, func(i, j int) bool {
		return t.tests[i].ID < t.tests[j].ID
	})

	return t.tests
}

// creates a new test and returns its ID
func (t *Task) AddTest(input []byte, answer []byte) int {
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

func (t *Task) AssignFilenameToTest(filename string, testID int) {
	_, ok1 := t.testIDToFilename[testID]
	_, ok2 := t.testFilenameToID[filename]
	if ok1 || ok2 {
		log.Fatalf("test with ID %d or filename %s already exists", testID, filename)
	}

	t.testIDToFilename[testID] = filename
	t.testFilenameToID[filename] = testID
}

func (t *Task) GetExamples() []example {
	return t.examples
}

func (t *Task) GetTaskName() string {
	return t.taskName
}

func (t *Task) SetTaskName(name string) {
	t.taskName = name
}

func (t *Task) GetProblemTags() []string {
	return t.problemTags
}

func (t *Task) SetProblemTags(tags []string) {
	t.problemTags = tags
}

func (t *Task) GetTaskAuthors() []string {
	return t.problemAuthors
}

func (t *Task) SetTaskAuthors(authors []string) {
	t.problemAuthors = authors
}

func (t *Task) GetOriginOlympiad() string {
	return t.originOlympiad
}

func (t *Task) SetOriginOlympiad(origin string) {
	t.originOlympiad = origin
}

func (t *Task) GetDifficultyOneToFive() int {
	return t.difficultyOneToFive
}

func (t *Task) SetDifficultyOneToFive(difficulty int) {
	t.difficultyOneToFive = difficulty
}

type TestGroupInfo struct {
	GroupID int
	Points  int
	Public  bool
	TestIDs []int
	Subtask int
}

func (t *Task) GetInfoOnTestGroup(id int) TestGroupInfo {
	return TestGroupInfo{
		GroupID: id,
		Points:  t.tGroupPoints[id],
		Public:  t.isTGroupPublic[id],
		TestIDs: t.tGroupTestIDs[id],
		Subtask: t.tGroupToStMap[id],
	}
}

func (t *Task) GetTestGroupIDs() []int {
	return t.testGroupIDs
}

func (t *Task) GetTestFilenameFromID(testID int) string {
	filename, ok := t.testIDToFilename[testID]
	if !ok {
		return ""
	}
	return filename
}

func (t *Task) testGroupWithIDExists(id int) bool {
	for i := 0; i < len(t.testGroupIDs); i++ {
		if t.testGroupIDs[i] == id {
			return true
		}
	}
	return false
}

func (t *Task) AddTestGroupWithID(groupID int, points int, public bool, testIDs []int, subtask int) error {
	if t.testGroupWithIDExists(groupID) {
		return fmt.Errorf("test group with ID %d already exists", groupID)
	}
	t.testGroupIDs = append(t.testGroupIDs, groupID)
	t.isTGroupPublic[groupID] = public
	t.tGroupPoints[groupID] = points
	t.tGroupTestIDs[groupID] = testIDs
	t.tGroupToStMap[groupID] = subtask

	return nil
}

func (t *Task) testGroupMexPositiveID() int {
	mex := 1
	found := true
	for found {
		found = false
		for i := 0; i < len(t.testGroupIDs); i++ {
			if t.testGroupIDs[i] == mex {
				found = true
				mex++
			}
		}
	}
	return mex
}

func (t *Task) AddTestGroup(points int, public bool, testIDs []int, subtask int) {
	err := t.AddTestGroupWithID(t.testGroupMexPositiveID(), points, public, testIDs, subtask)
	if err != nil {
		log.Fatalf("error adding test group: %v", err)
	}
}
