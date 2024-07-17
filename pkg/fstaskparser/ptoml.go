package fstaskparser

import (
	"bytes"
	"log"

	"github.com/pelletier/go-toml/v2"
)

type ProblemTOML struct {
	Specification   string           `toml:"specification"`
	TaskName        string           `toml:"task_name"`
	Metadata        PTomlMetadata    `toml:"metadata"`
	Constraints     PTomlConstraints `toml:"constraints"`
	TestGroups      []PTomlTestGroup `toml:"test_groups"`
	VisInpSTs       []int            `toml:"visible_input_subtasks"`
	TestIDOverwrite map[string]int   `toml:"test_id_overwrite"`
}

type PTomlMetadata struct {
	ProblemTags        []string `toml:"problem_tags"`
	DifficultyFrom1To5 int      `toml:"difficulty_1_to_5"`
	TaskAuthors        []string `toml:"task_authors"`
	OriginOlympiad     string   `toml:"origin_olympiad"`
}

type PTomlConstraints struct {
	MemoryMegabytes int     `toml:"memory_megabytes"`
	CPUTimeSeconds  float64 `toml:"cpu_time_seconds"`
}

// PTomlTestGroup is a structure to store groups used in LIO test format
type PTomlTestGroup struct {
	GroupID    int      `toml:"group_id"`
	Points     int      `toml:"points"`
	Public     bool     `toml:"public"`
	Subtask    *int     `toml:"subtask,omitempty"` // nil if subtask id not found
	TestIDs    []int    `toml:"test_ids,omitempty"`
	TestFnames []string `toml:"test_filenames,omitempty"` // either one is fine
}

func (task *Task) encodeProblemTOML() ([]byte, error) {
	t := ProblemTOML{
		Specification: proglvFSTaskFormatSpecVersOfScript,
		TaskName:      task.taskName,
		Metadata: PTomlMetadata{
			ProblemTags:        task.problemTags,
			DifficultyFrom1To5: task.difficultyOneToFive,
			TaskAuthors:        task.problemAuthors,
			OriginOlympiad:     task.originOlympiad,
		},
		Constraints: PTomlConstraints{
			MemoryMegabytes: task.memoryMegabytes,
			CPUTimeSeconds:  task.cpuTimeSeconds,
		},
		TestGroups: []PTomlTestGroup{},
		VisInpSTs:  []int{},
	}
	t.Specification = proglvFSTaskFormatSpecVersOfScript

	// fill test groups
	tGroupsWithSTAssigned := 0
	for _, tg := range task.testGroups {
		ptomlTestGroup := PTomlTestGroup{
			GroupID: tg.GroupID,
			Points:  0,
			Public:  false,
			Subtask: nil,
			TestIDs: tg.TestIDs,
		}

		tGroupPoints, ok := task.tGroupPoints[tg.GroupID]
		if ok {
			ptomlTestGroup.Points = tGroupPoints
		} else {
			log.Fatalf("Group %d has no points assigned\n", tg.GroupID)
		}

		isPublic, ok := task.isTGroupPublic[tg.GroupID]
		if ok {
			ptomlTestGroup.Public = isPublic
		}

		tGroupSt, ok := task.tGroupToStMap[tg.GroupID]
		if ok {
			ptomlTestGroup.Subtask = &tGroupSt
			tGroupsWithSTAssigned++
		}

		t.TestGroups = append(t.TestGroups, ptomlTestGroup)
	}
	if tGroupsWithSTAssigned != 0 && tGroupsWithSTAssigned != len(task.testGroups) {
		log.Fatalf("Some test groups have subtasks assigned, while others don't\n")
	}

	// fill visible input subtasks
	t.VisInpSTs = append(t.VisInpSTs, task.visibleInputSubtasks...)

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
