package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

const redactText = "[Redacted]"

func RedactKeys(data map[string]any, keys []string) map[string]any {
	return redact(copyData(data), mergeKeys(keys), "")
}

func copyData(data map[string]any) map[string]any {
	var result map[string]any

	dataAsBytes, _ := json.Marshal(data)
	_ = json.Unmarshal(dataAsBytes, &result)

	return result
}

func convertKeysToLower(keys []string) {
	for i, key := range keys {
		keys[i] = strings.ToLower(key)
	}
}

func mergeKeys(keys []string) []string {
	envKeys := env.GetRedactKeys()
	envKeys = append(envKeys, keys...)

	convertKeysToLower(envKeys)

	return envKeys
}

func isBase64(str string) bool {
	return regexp.MustCompile(`^data:([a-z]+\/[a-z]+(;[a-z]+=[a-z]+)?)?(;base64)?,([a-zA-Z0-9+/]+={0,2})+$`).MatchString(str)
}

func redact(data map[string]any, keys []string, previousKey string) map[string]any {
	if len(data) == 0 {
		return data
	}

	for key, value := range data {
		if value == nil {
			continue
		}

		nextKey := strings.ToLower(key)
		if previousKey != "" {
			nextKey = strings.ToLower(fmt.Sprintf("%s.%s", previousKey, key))
		}

		valueKind := reflect.TypeOf(value).Kind()
		isRedacted := slices.Contains(keys, nextKey) ||
			(valueKind == reflect.String && isBase64(value.(string))) ||
			slices.Contains(keys, key)

		if isRedacted {
			if valueKind == reflect.Slice {
				for i := range value.([]any) {
					data[key].([]any)[i] = redactText
				}

				continue
			}

			data[key] = redactText
			continue
		}

		switch valueKind {
		case reflect.Slice, reflect.Array:
			for i, v := range value.([]any) {
				checkKey := fmt.Sprintf("%s.%d", nextKey, i)

				if str, ok := v.(string); ok && isBase64(str) || slices.Contains(keys, checkKey) {
					v = redactText
				}

				if valueAsMap, ok := v.(map[string]any); ok {
					v = redact(valueAsMap, keys, checkKey)
				}

				value.([]any)[i] = v
			}
		case reflect.Map:
			if input, ok := value.(map[string]any); ok {
				data[key] = redact(input, keys, nextKey)
			}
		}
	}

	return data
}
