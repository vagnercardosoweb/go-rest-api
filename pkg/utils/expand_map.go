package utils

import (
	"strconv"
	"strings"
)

func ExpandMap(flatMap map[string]any) map[string]any {
	if flatMap == nil {
		return flatMap
	}

	var resultRoot any = make(map[string]any)

	for k, v := range flatMap {
		parts := strings.Split(k, ".")
		resultRoot = setValueRecursive(resultRoot, parts, v)
	}

	// If the result is a map, return it as a map[string]any
	if finalResultMap, ok := resultRoot.(map[string]any); ok {
		return finalResultMap
	}

	// If the result is not a map, return an empty map
	return make(map[string]any)
}

func expandSlice(s *[]any, minIndex int) {
	if *s == nil {
		*s = make([]any, 0)
	}

	if minIndex >= len(*s) {
		newSlice := make([]any, minIndex+1)
		copy(newSlice, *s)
		*s = newSlice
	}
}

func setMapRecursive(target any, currentPath string, remainingPath []string, value any) map[string]any {
	var result map[string]any

	if target != nil {
		var ok bool
		result, ok = target.(map[string]any)

		if !ok {
			result = make(map[string]any)
		}
	} else {
		result = make(map[string]any)
	}

	result[currentPath] = setValueRecursive(result[currentPath], remainingPath, value)

	return result
}

func setSliceRecursive(target any, idx int, remainingPath []string, value any) []any {
	var slice []any

	if target != nil {
		var ok bool
		slice, ok = target.([]any)

		if !ok {
			slice = make([]any, 0)
		}
	} else {
		slice = make([]any, 0)
	}

	expandSlice(&slice, idx)
	slice[idx] = setValueRecursive(slice[idx], remainingPath, value)

	return slice
}

func setValueRecursive(target any, path []string, value any) any {
	if len(path) == 0 {
		return value
	}

	key := path[0]
	remainingPath := path[1:]

	if idx, err := strconv.Atoi(key); err == nil {
		return setSliceRecursive(target, idx, remainingPath, value)
	}

	return setMapRecursive(target, key, remainingPath, value)
}
