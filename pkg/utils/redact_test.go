package utils

import (
	"encoding/json"
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
		"data:application/pdf;base64,JVBERi0xLjQKJeLjz9MKMyAwIG9iago8PC9MZW5ndGggNiAwIFIvRmlsdGVyIC9GbGF0ZURlY29kZT4+CnN0cmVhbQp4nLWT",
		"data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7D",
	},
	"sliceInt":    []int{1, 2, 3},
	"sliceString": []string{"a", "b", "c"},
	"headers": map[string]any{
		"Authorization": "Bearer {token}",
		"x-api-key":     "x-api-key-1",
	},
	"nested": map[string]any{
		"age":  29,
		"name": "John Doe",
		"nested": map[string]any{
			"password":  "12345678-3",
			"x-api-key": "x-api-key-2",
		},
	},
	"sliceNested": []map[string]any{
		{
			"name":      "John Doe",
			"sensitive": "12345678-4",
			"otherArray": []map[string]any{
				{
					"test": "test",
				},
			},
			"age": 29,
		},
		{
			"name":      "John Doe",
			"sensitive": "12345678-4",
			"otherArray": []map[string]any{
				{
					"test": "test",
				},
			},
			"age": 29,
		},
	},
}

var expectedMap = map[string]any{
	"age":               28,
	"name":              "John Doe",
	"document":          redactText,
	"email":             redactText,
	"password":          redactText,
	"passwordConfirm":   redactText,
	"techs":             []string{redactText, redactText, redactText},
	"nilValue":          nil,
	"pointerValue":      "pointer",
	"fileBase64":        redactText,
	"fileBase64AsArray": [3]string{redactText, redactText, redactText},
	"sliceInt":          []int{1, 2, 3},
	"sliceString":       []string{"a", redactText, "c"},
	"headers": map[string]any{
		"Authorization": "Bearer {token}",
		"x-api-key":     redactText,
	},
	"nested": map[string]any{
		"age":  29,
		"name": "John Doe",
		"nested": map[string]any{
			"x-api-key": "x-api-key-2",
			"password":  redactText,
		},
	},
	"sliceNested": []map[string]any{
		{
			"name":      "John Doe",
			"sensitive": redactText,
			"otherArray": []map[string]any{
				{
					"test": "test",
				},
			},
			"age": 29,
		},
		{
			"name":      redactText,
			"sensitive": "12345678-4",
			"otherArray": []map[string]any{
				{
					"test": redactText,
				},
			},
			"age": 29,
		},
	},
}

func TestRedactKeys(t *testing.T) {
	comparator := make(map[string]any)

	expectedBytes, _ := json.Marshal(expectedMap)
	_ = json.Unmarshal(expectedBytes, &comparator)

	keys := []string{
		"email", "password", "passwordConfirm", "headers.x-api-key",
		"document", "sliceNested.0.sensitive", "sliceString.1", "techs",
		"sliceNested.1.name", "sliceNested.1.otherArray.0.test",
	}

	result := RedactKeys(originalMap, keys)

	assert.Equal(t, originalMap["password"], "12345678-1")
	assert.Equal(t, result, comparator)
}
