package fstaskparser

func (t *Task) GetCPUTimeLimitInSeconds() (float64, error) {
	return t.cpuTimeSeconds, nil
}

func (t *Task) GetMemoryLimitInMegabytes() (int, error) {
	return t.memoryMegabytes, nil
}
