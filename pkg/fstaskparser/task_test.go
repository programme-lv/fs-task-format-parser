package fstaskparser_test

import (
	"path/filepath"
	"testing"

	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoringComplexTask(t *testing.T) {
	testDir := filepath.Join(".", "..", "..", "testdata", "kvadrputekl")
	task, err := fstaskparser.Read(testDir)
	if err != nil {
		t.Fatalf("failed to read task: %v", err)
	}

	originalCPUTimeLimitInSeconds, err := task.GetCPUTimeInSeconds()
	require.NoErrorf(t, err, "failed to get cpu time limit: %v", err)

	assert.Equal(t, 0.5, originalCPUTimeLimitInSeconds)

	err = task.Store(filepath.Join(".", "output"))
	if err != nil {
		t.Fatalf("failed to store task: %v", err)
	}

	task2, err := fstaskparser.Read(filepath.Join(".", "output"))
	require.NoErrorf(t, err, "failed to read task: %v", err)

	writtenCPUTimeLimitInSeconds, err := task2.GetCPUTimeInSeconds()
	require.NoErrorf(t, err, "failed to get cpu time limit: %v", err)

	assert.Equal(t, 0.5, writtenCPUTimeLimitInSeconds)
}

/*
kvadrputekl problem.toml

specification = '2.2'
task_name = 'Kvadrātveida putekļsūcējs'
visible_input_subtasks = [1]

[metadata]
  problem_tags = []
  difficulty_1_to_5 = 3
  task_authors = []
  origin_olympiad = 'LIO'

[constraints]
  memory_megabytes = 256
  cpu_time_seconds = 0.5

[[test_groups]]
  group_id = 1
  points = 3
  subtask = 1
  public = true
  test_filenames = ['kp01a.in', 'kp01b.in', 'kp01c.in', 'kp01a.out', 'kp01b.out', 'kp01c.out']

[[test_groups]]
  group_id = 2
  points = 8
  subtask = 2
  public = true
  test_filenames = ['kp02a.in', 'kp02b.in', 'kp02c.in', 'kp02a.out', 'kp02b.out', 'kp02c.out']
*/
