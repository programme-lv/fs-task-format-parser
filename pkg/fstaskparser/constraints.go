package fstaskparser

import (
	"fmt"
	"log"

	"github.com/pelletier/go-toml/v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func readCPUTimeLimitInSeconds(specVers string, tomlContent string) (float64, error) {
	log.Printf("Reading CPU time limit for specification version: %s\n", specVers)
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		log.Printf("Error comparing semversions: %v\n", err)
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		log.Printf("Unsupported specification version: %s\n", specVers)
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
		log.Printf("Failed to unmarshal the CPU time limit: %v\n", err)
		return 0, fmt.Errorf("failed to unmarshal the cpu time limit: %w", err)
	}

	log.Printf("Successfully read CPU time limit: %f seconds\n", tomlStruct.Constraints.CPUTimeLimitInSeconds)
	return tomlStruct.Constraints.CPUTimeLimitInSeconds, nil
}

func readMemoryLimitInMegabytes(specVers string, tomlContent string) (int, error) {
	log.Printf("Reading memory limit for specification version: %s\n", specVers)
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		log.Printf("Error comparing semversions: %v\n", err)
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		log.Printf("Unsupported specification version: %s\n", specVers)
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
		log.Printf("Failed to unmarshal the memory limit: %v\n", err)
		return 0, fmt.Errorf("failed to unmarshal the memory limit: %w", err)
	}

	log.Printf("Successfully read memory limit: %d megabytes\n", tomlStruct.Constraints.MemoryLimitInMegabytes)
	return tomlStruct.Constraints.MemoryLimitInMegabytes, nil
}
