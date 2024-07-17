package fstaskparser

func (t *Task) GetCPUTimeLimitInSeconds() float64 {
	return t.cpuTimeSeconds
}

func (t *Task) GetMemoryLimitInMegabytes() int {
	return t.memoryMegabytes
}

func (t *Task) GetFullTaskName() string {
	return t.taskName
}

func (t *Task) GetTests() []Test {
	return t.tests
}

func (t *Task) GetExamples() []Example {
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

func (t *Task) GetPublicTestGroupIDs() []int {
	res := []int{}

	for i, id := range t.testGroupIDs {
		if t.isTGroupPublic[i] {
			res = append(res, id)
		}
	}

	return res

}

func (t *Task) GetTestFilenameFromID(testID int) *string {
	filename, ok := t.testIDToFilename[testID]
	if !ok {
		return nil
	}
	return &filename
}
