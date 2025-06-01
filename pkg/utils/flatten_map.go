package utils

import (
	"fmt"
	"strconv"
)

func FlattenMap(flatMap map[string]any) map[string]any {
	if flatMap == nil {
		return nil
	}

	if len(flatMap) == 0 {
		return make(map[string]any)
	}

	result := make(map[string]any)
	flattenNestedMap(result, "", flatMap)
	return result
}

func buildPrefixedKey(prefix string, key string) string {
	if prefix != "" {
		return fmt.Sprintf("%s.%s", prefix, key)
	}

	return key
}

func flattenNestedMap(target map[string]any, prefix string, result any) {
	if result == nil {
		target[prefix] = nil
		return
	}

	if mapValue, ok := result.(map[string]any); ok {
		if len(mapValue) == 0 {
			target[prefix] = make(map[string]any)
		}

		for key, value := range mapValue {
			flattenNestedMap(target, buildPrefixedKey(prefix, key), value)
		}

		return
	}

	if arrayValue, ok := result.([]any); ok {
		if len(arrayValue) == 0 {
			target[prefix] = make([]any, 0)
		}

		for i, value := range arrayValue {
			flattenNestedMap(target, buildPrefixedKey(prefix, strconv.Itoa(i)), value)
		}

		return
	}

	target[prefix] = result
}
