package fstaskparser

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

func readCPUTimeLimitInSeconds(specVers string, tomlContent string) (float64, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return 0, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type constraintsStruct struct {
		CPUTimeLimitInSeconds float64 `toml:"cpu_time_seconds"`
	}

	tomlStruct := struct {
		Constraints constraintsStruct `toml:"constraints"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal the cpu time limit: %w", err)
	}

	return tomlStruct.Constraints.CPUTimeLimitInSeconds, nil
}

func readMemoryLimitInMegabytes(specVers string, tomlContent string) (int, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return 0, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type constraintsStruct struct {
		MemoryLimitInMegabytes int `toml:"memory_megabytes"`
	}

	tomlStruct := struct {
		Constraints constraintsStruct `toml:"constraints"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal the memory limit: %w", err)
	}

	return tomlStruct.Constraints.MemoryLimitInMegabytes, nil
}
