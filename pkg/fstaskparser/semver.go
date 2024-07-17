package fstaskparser

import (
	"fmt"
	"strconv"
	"strings"
)

func largerOrEqualSemVersionThan(a, b string) (bool, error) {
	cmpRes, err := getCmpSemVersionsResult(a, b)
	if err != nil {
		return false, fmt.Errorf("error comparing sem versions: %w", err)
	}

	return cmpRes >= 0, nil
}

// func smallerOrEqualSemVersionThan(a, b string) (bool, error) {
// 	cmpRes, err := getCmpSemVersionsResult(a, b)
// 	if err != nil {
// 		return false, fmt.Errorf("error comparing sem versions: %w", err)
// 	}

// 	return cmpRes <= 0, nil
// }

// getCmpSemVersionsResult compares two semantic versions and returns the result of the comparison.
//
// Parameters:
// - a: the first semantic version to compare.
// - b: the second semantic version to compare.
//
// Returns:
// - int: the result of the comparison. 1 if a > b, -1 if a < b, 0 if a == b.
// - error: an error if there was an issue converting the semantic versions to int slices for comparison.
func getCmpSemVersionsResult(a, b string) (int, error) {
	aInts, err := formatStringAsIntSliceForComparision(a)
	if err != nil {
		msg := fmt.Errorf("error converting to int slice for a: %w", err)
		return 0, msg
	}

	bInts, err := formatStringAsIntSliceForComparision(b)
	if err != nil {
		msg := fmt.Errorf("error converting to int slice for b: %w", err)
		return 0, msg

	}

	cmpRes := compareIntSlices(aInts, bInts)

	return cmpRes, nil
}

func formatStringAsIntSliceForComparision(s string) ([]int, error) {
	// if starts with a "v", remove it
	if s[0] == 'v' {
		s = s[1:]
	}

	// if ends with a ".0", remove it
	for strings.HasSuffix(s, ".0") {
		s = s[:len(s)-2]
	}

	if s == "" {
		return nil, fmt.Errorf("empty string when converting to int slice")
	}

	// split by "."
	parts := strings.Split(s, ".")
	res := make([]int, len(parts))
	for i, part := range parts {
		var err error
		res[i], err = strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("error converting version part to int: %w", err)
		}
	}

	return res, nil
}

// 0 if a == b
// 1 if a > b
// -1 if a < b
func compareIntSlices(a, b []int) int {
	for i := 0; i < len(a) || i < len(b); i++ {
		if i >= len(a) { // a is shorter and equal in the common part
			return -1
		}
		if i >= len(b) { // b is shorter and equal in the common part
			return 1
		}
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
