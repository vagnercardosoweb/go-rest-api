package utils

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
)

const censorText = "[Redacted]"

func isBase64(str string) bool {
	return regexp.MustCompile(`^data:([a-z]+\/[a-z]+(;[a-z]+=[a-z]+)?)?(;base64)?,([a-zA-Z0-9+/]+={0,2})+$`).MatchString(str)
}

func RedactKeys(data map[string]any, keys []string) map[string]any {
	if len(keys) == 0 {
		return data
	}

	encodedData, _ := json.Marshal(data)

	var redactedData map[string]any
	_ = json.Unmarshal(encodedData, &redactedData)

	for key, value := range redactedData {
		if value == nil {
			continue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			if isBase64(value.(string)) {
				redactedData[key] = censorText
			}
		case reflect.Slice, reflect.Array:
			for i, v := range value.([]any) {
				if valueAsString, ok := v.(string); ok && isBase64(valueAsString) {
					v = censorText
				}

				if valueAsMap, ok := v.(map[string]any); ok {
					v = RedactKeys(valueAsMap, keys)
				}

				value.([]any)[i] = v
			}
		case reflect.Map:
			if input, ok := value.(map[string]any); ok {
				redactedData[key] = RedactKeys(input, keys)
			}
		}

		for _, k := range keys {
			if strings.EqualFold(k, key) {
				redactedData[key] = censorText
			}
		}
	}

	return redactedData
}
