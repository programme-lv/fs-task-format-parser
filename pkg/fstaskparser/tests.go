package fstaskparser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func readTestsDir(srcDirPath string, fnameToID map[string]int) ([]test, error) {
	dir := filepath.Join(srcDirPath, "tests")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading tests directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	tests := make([]test, 0, len(entries)/2)

	for i := 0; i < len(entries); i += 2 {
		inPath := filepath.Join(dir, entries[i].Name())
		ansPath := filepath.Join(dir, entries[i+1].Name())

		inFilename := entries[i].Name()
		ansFilename := entries[i+1].Name()

		inFilenameBase := strings.TrimSuffix(inFilename, filepath.Ext(inFilename))
		ansFilenameBase := strings.TrimSuffix(ansFilename, filepath.Ext(ansFilename))

		if inFilenameBase != ansFilenameBase {
			return nil, fmt.Errorf("input and answer file base names do not match: %s, %s", inFilenameBase, ansFilenameBase)
		}

		// sometimes the test answer is stored as .out, sometimes as .ans
		if strings.Contains(inFilename, ".ans") || strings.Contains(ansFilename, ".in") {
			// swap the file paths
			inPath, ansPath = ansPath, inPath
		}

		input, err := os.ReadFile(inPath)
		if err != nil {
			return nil, fmt.Errorf("error reading input file: %w", err)
		}

		answer, err := os.ReadFile(ansPath)
		if err != nil {
			return nil, fmt.Errorf("error reading answer file: %w", err)
		}

		// check if mapping to id exists
		if _, ok := fnameToID[inFilenameBase]; !ok {
			return nil, fmt.Errorf("mapping from filename to id does not exist: %s", inFilenameBase)
		}

		tests = append(tests, test{
			ID:     fnameToID[inFilenameBase],
			Input:  input,
			Answer: answer,
		})
	}

	return tests, nil
}

func readExamplesDir(srcDirPath string) ([]example, error) {
	dir := filepath.Join(srcDirPath, "examples")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading examples directory: %w", err)
	}
	// tests are to be read exactly like examples

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	examples := make([]example, 0, len(entries)/2)

	for i := 0; i < len(entries); i += 2 {
		inPath := filepath.Join(dir, entries[i].Name())
		ansPath := filepath.Join(dir, entries[i+1].Name())

		inFilename := entries[i].Name()
		ansFilename := entries[i+1].Name()

		inFilenameBase := strings.TrimSuffix(inFilename, filepath.Ext(inFilename))
		ansFilenameBase := strings.TrimSuffix(ansFilename, filepath.Ext(ansFilename))

		if inFilenameBase != ansFilenameBase {
			return nil, fmt.Errorf("input and answer file base names do not match: %s, %s", inFilenameBase, ansFilenameBase)
		}

		// sometimes the test answer is stored as .out, sometimes as .ans
		if strings.Contains(inFilename, ".ans") || strings.Contains(ansFilename, ".in") {
			// swap the file paths
			inPath, ansPath = ansPath, inPath
		}

		input, err := os.ReadFile(inPath)
		if err != nil {
			return nil, fmt.Errorf("error reading input file: %w", err)
		}

		answer, err := os.ReadFile(ansPath)
		if err != nil {
			return nil, fmt.Errorf("error reading answer file: %w", err)
		}

		examples = append(examples, example{
			ID:     (i / 2) + 1,
			Input:  input,
			Output: answer,
			Name:   &inFilenameBase,
		})
	}

	return examples, nil
}

func readTestIDOverwrite(specVers string, tomlContent []byte) (map[string]int, error) {
	semVerCmpRes, err := getCmpSemVersionsResult(specVers, "v2.3.0")
	if err != nil {
		return nil, fmt.Errorf("error comparing sem versions: %w", err)
	}

	if semVerCmpRes < 0 {
		log.Printf("warning: skipping reading test id overwrite (spec version: %s)\n", specVers)
		// return empty map
		return make(map[string]int), nil
	}

	tomlStruct := struct {
		TestIDOverwrite map[string]int `toml:"test_id_overwrite"`
	}{}

	err = toml.Unmarshal(tomlContent, &tomlStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the test id overwrite: %w", err)
	}

	return tomlStruct.TestIDOverwrite, nil
}

func readTestFNamesSorted(dirPath string) ([]string, error) {
	fnames, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading test filenames: %w", err)
	}

	sort.Slice(fnames, func(i, j int) bool {
		return fnames[i].Name() < fnames[j].Name()
	})

	if len(fnames)%2 != 0 {
		return nil, fmt.Errorf("odd number of test filenames")
	}

	res := make([]string, 0, len(fnames)/2)
	for i := 0; i < len(fnames); i += 2 {
		a_name := fnames[i].Name()
		// remove extension
		a_name = a_name[:len(a_name)-len(filepath.Ext(a_name))]

		b_name := fnames[i+1].Name()
		// remove extension
		b_name = b_name[:len(b_name)-len(filepath.Ext(b_name))]

		if a_name != b_name {
			return nil, fmt.Errorf("input and answer file base names do not match: %s, %s", a_name, b_name)
		}

		res = append(res, a_name)
	}

	return res, nil
}
