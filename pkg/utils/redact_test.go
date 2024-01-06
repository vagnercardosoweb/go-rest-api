package utils

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pointer = "pointer"
var nilValue *string = nil

var originalMap = map[string]any{
	"age":             28,
	"name":            "John Doe",
	"document":        "000.000.000-00",
	"email":           "johndoe@test.com",
	"password":        "12345678-1",
	"passwordConfirm": "12345678-1",
	"techs":           []string{"Go", "Node.js", "React"},
	"nilValue":        nilValue,
	"pointerValue":    &pointer,
	"fileBase64":      "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7D",
	"fileBase64AsArray": [3]string{
		"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7D",
		"data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7D",
		"data:application/pdf;base64,JVBERi0xLjQKJeLjz9MKMyAwIG9iago8PC9MZW5ndGggNiAwIFIvRmlsdGVyIC9GbGF0ZURlY29kZT4+CnN0cmVhbQp4nLWT",
	},
	"sliceInt":    []int{1, 2, 3},
	"sliceString": []string{"a", "b", "c"},
	"headers": map[string]any{
		"Authorization": "Bearer {token}",
		"x-api-key":     "x-api-key-1",
	},
	"nested": map[string]any{
		"age":       "29",
		"name":      "John Doe",
		"password":  "12345678-2",
		"x-api-key": "x-api-key-2",
		"nested": map[string]any{
			"password":  "12345678-3",
			"x-api-key": "x-api-key-3",
		},
	},
	"sliceNested": []map[string]any{
		{
			"name":     "John Doe",
			"password": "12345678-4",
			"age":      "29",
		},
		{
			"name":     "John Doe",
			"password": "12345678-4",
			"age":      29,
		},
	},
}

var expectedMap = map[string]any{
	"age":               28,
	"name":              "John Doe",
	"document":          censorText,
	"email":             censorText,
	"password":          censorText,
	"passwordConfirm":   censorText,
	"techs":             []string{"Go", "Node.js", "React"},
	"nilValue":          nil,
	"pointerValue":      "pointer",
	"fileBase64":        censorText,
	"fileBase64AsArray": [3]string{censorText, censorText, censorText},
	"sliceInt":          []int{1, 2, 3},
	"sliceString":       []string{"a", "b", "c"},
	"headers": map[string]any{
		"Authorization": "Bearer {token}",
		"x-api-key":     censorText,
	},
	"nested": map[string]any{
		"age":       "29",
		"name":      "John Doe",
		"password":  censorText,
		"x-api-key": censorText,
		"nested": map[string]any{
			"password":  censorText,
			"x-api-key": censorText,
		},
	},
	"sliceNested": []map[string]any{
		{
			"name":     "John Doe",
			"password": censorText,
			"age":      "29",
		},
		{
			"name":     "John Doe",
			"password": censorText,
			"age":      29,
		},
	},
}

var keys = []string{"password", "passwordConfirm", "x-api-key", "document", "email"}

func TestRedactKeysWithKeys(t *testing.T) {
	expectedAsBytes, _ := json.Marshal(expectedMap)
	expectedAsMap := make(map[string]any)
	_ = json.Unmarshal(expectedAsBytes, &expectedAsMap)

	result := RedactKeys(originalMap, keys)

	assert.True(t, reflect.DeepEqual(result, expectedAsMap))
	assert.True(t, reflect.DeepEqual(originalMap, originalMap))
}

func TestRedactKeysWithoutKeys(t *testing.T) {
	result := RedactKeys(originalMap, []string{})
	assert.True(t, reflect.DeepEqual(result, originalMap))
}
