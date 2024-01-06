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

func RedactKeys(metadata map[string]any, keys []string) map[string]any {
	if len(keys) == 0 {
		return metadata
	}

	metadataAsBytes, _ := json.Marshal(metadata)
	var metadataAsMap map[string]any
	_ = json.Unmarshal(metadataAsBytes, &metadataAsMap)

	for key, value := range metadataAsMap {
		if value == nil {
			continue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			if isBase64(value.(string)) {
				metadataAsMap[key] = censorText
			}
		case reflect.Slice, reflect.Array:
			for i, v := range value.([]interface{}) {
				if valueAsString, ok := v.(string); ok && isBase64(valueAsString) {
					v = censorText
				}

				if valueAsMap, ok := v.(map[string]any); ok {
					v = RedactKeys(valueAsMap, keys)
				}

				value.([]interface{})[i] = v
			}
		case reflect.Map:
			if input, ok := value.(map[string]any); ok {
				metadataAsMap[key] = RedactKeys(input, keys)
			}
		}

		for _, k := range keys {
			if strings.ToLower(key) == strings.ToLower(k) {
				metadataAsMap[key] = censorText
			}
		}
	}

	return metadataAsMap
}
