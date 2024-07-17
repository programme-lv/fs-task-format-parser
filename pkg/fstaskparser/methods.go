package fstaskparser

func (t *Task) GetCPUTimeInSeconds() (float64, error) {
	return t.cpuTimeSeconds, nil
}
