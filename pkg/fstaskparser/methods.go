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
