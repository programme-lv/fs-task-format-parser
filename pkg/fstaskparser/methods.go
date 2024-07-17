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

func (t *Task) GetTestGroups() []TestGroup {
	return t.testGroups
}

func (t *Task) GetPublicTestGroups() []TestGroup {
	var publicTestGroups []TestGroup
	for i, tg := range t.testGroups {
		if t.isTGroupPublic[i] {
			publicTestGroups = append(publicTestGroups, tg)
		}
	}
	return publicTestGroups
}

func (t *Task) GetTestIDFilename(testID int) *string {
	filename, ok := t.testIDToFilename[testID]
	if !ok {
		return nil
	}
	return &filename
}
