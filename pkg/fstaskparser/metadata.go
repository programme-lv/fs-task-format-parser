package fstaskparser

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

func readProblemTags(specVers string, tomlContent string) ([]string, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return nil, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return nil, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type metadataStruct struct {
		ProblemTags []string `toml:"problem_tags"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	return tomlStruct.Metadata.ProblemTags, nil
}

func readProblemAuthors(specVers string, tomlContent string) ([]string, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return nil, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return nil, fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type metadataStruct struct {
		ProblemAuthors []string `toml:"task_authors"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	return tomlStruct.Metadata.ProblemAuthors, nil
}

func readOriginOlympiad(specVers string, tomlContent string) (string, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return "", fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return "", fmt.Errorf("unsupported specification version: %s", specVers)
	}

	type metadataStruct struct {
		OriginOlympiad *string `toml:"origin_olympiad"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	res := ""
	if tomlStruct.Metadata.OriginOlympiad != nil {
		res = *tomlStruct.Metadata.OriginOlympiad
	}
	return res, nil
}

func readDifficultyOneToFive(specVers string, tomlContent string) (int, error) {
	cmpres, err := largerOrEqualSemVersionThan(specVers, "2.0")
	if err != nil {
		return 0, fmt.Errorf("error comparing semversions: %w", err)
	}
	if !cmpres {
		return 0, fmt.Errorf("unsupported specification version: %s", specVers)
	}
	type metadataStruct struct {
		DifficultyFrom1To5 *int `toml:"difficulty_1_to_5"`
	}

	tomlStruct := struct {
		Metadata metadataStruct `toml:"metadata"`
	}{}

	err = toml.Unmarshal([]byte(tomlContent), &tomlStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal the problem tags: %w", err)
	}

	res := 0
	if tomlStruct.Metadata.DifficultyFrom1To5 != nil {
		res = *tomlStruct.Metadata.DifficultyFrom1To5
	}

	return res, nil
}
