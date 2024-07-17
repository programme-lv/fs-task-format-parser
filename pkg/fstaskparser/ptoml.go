package fstaskparser

import (
	"bytes"
	"log"

	"github.com/pelletier/go-toml/v2"
)

type ProblemTOML struct {
	Specification string           `toml:"specification"`
	TaskName      string           `toml:"task_name"`
	Metadata      PTomlMetadata    `toml:"metadata"`
	Constraints   PTomlConstraints `toml:"constraints"`
	TestGroups    []PTomlTestGroup `toml:"test_groups"`
	VisInpSTs     []int            `toml:"visible_input_subtasks"`
}

type PTomlMetadata struct {
	ProblemTags        []string `toml:"problem_tags"`
	DifficultyFrom1To5 int      `toml:"difficulty_1_to_5"`
	TaskAuthors        []string `toml:"task_authors"`
	OriginOlympiad     *string  `toml:"origin_olympiad"`
}

type PTomlConstraints struct {
	MemoryMegabytes int     `toml:"memory_megabytes"`
	CPUTimeSeconds  float64 `toml:"cpu_time_seconds"`
}

// PTomlTestGroup is a structure to store groups used in LIO test format
type PTomlTestGroup struct {
	GroupID int   `toml:"group_id"`
	Points  int   `toml:"points"`
	Public  bool  `toml:"public"`
	Subtask int   `toml:"subtask"`
	TestIDs []int `toml:"test_ids"`
	// TestFnames []string `toml:"test_filenames"`
}

func (task *Task) encodeProblemTOML() ([]byte, error) {
	difficultyOneToFive := 0
	if task.difficultyOneToFive != nil {
		difficultyOneToFive = *task.difficultyOneToFive
	}

	t := ProblemTOML{
		Specification: proglvFSTaskFormatSpecVersion,
		TaskName:      task.taskName,
		Metadata: PTomlMetadata{
			ProblemTags:        task.problemTags,
			DifficultyFrom1To5: difficultyOneToFive,
			TaskAuthors:        task.problemAuthors,
			OriginOlympiad:     task.originOlympiad,
		},
		Constraints: PTomlConstraints{
			MemoryMegabytes: task.memoryMegabytes,
			CPUTimeSeconds:  task.cpuTimeSeconds,
		},
		TestGroups: []PTomlTestGroup{}, // TODO: fill test groups
		VisInpSTs:  []int{},
	}

	// fill test groups
	t.VisInpSTs = append(t.VisInpSTs, task.visibleInputSubtasks...)

	// fill visible input subtasks

	t.Specification = "2.2"

	buf := bytes.NewBuffer(make([]byte, 0))
	err := toml.NewEncoder(buf).
		SetTablesInline(false).
		// SetArraysMultiline(true).
		SetIndentTables(true).Encode(t)

	if err != nil {
		log.Fatalf("Failed to marshal the problem.toml: %v\n", err)
	}

	return buf.Bytes(), nil
}
