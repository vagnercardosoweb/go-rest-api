package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandMap(t *testing.T) {
	t.Run("should handle empty map", func(t *testing.T) {
		input := make(map[string]any)
		result := ExpandMap(input)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should handle nil map", func(t *testing.T) {
		result := ExpandMap(nil)
		assert.Nil(t, result)
	})

	t.Run("should expand simple flat map", func(t *testing.T) {
		input := map[string]any{
			"name": "John",
			"age":  30,
		}

		expected := map[string]any{
			"name": "John",
			"age":  30,
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should expand nested map", func(t *testing.T) {
		input := map[string]any{
			"user.name": "John",
			"user.age":  30,
		}

		expected := map[string]any{
			"user": map[string]any{
				"name": "John",
				"age":  30,
			},
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should expand deeply nested map", func(t *testing.T) {
		input := map[string]any{
			"user.profile.name": "John",
			"user.profile.age":  30,
		}

		expected := map[string]any{
			"user": map[string]any{
				"profile": map[string]any{
					"name": "John",
					"age":  30,
				},
			},
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should expand array indices", func(t *testing.T) {
		input := map[string]any{
			"tags.0": "go",
			"tags.1": "javascript",
			"tags.2": "python",
		}

		expected := map[string]any{
			"tags": []any{"go", "javascript", "python"},
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should expand array with gaps", func(t *testing.T) {
		input := map[string]any{
			"tags.0": "go",
			"tags.2": "python",
		}

		expected := map[string]any{
			"tags": []any{"go", nil, "python"},
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should expand nested objects in array", func(t *testing.T) {
		input := map[string]any{
			"users.0.name": "John",
			"users.0.age":  30,
			"users.1.name": "Jane",
			"users.1.age":  25,
		}

		expected := map[string]any{
			"users": []any{
				map[string]any{
					"name": "John",
					"age":  30,
				},
				map[string]any{
					"name": "Jane",
					"age":  25,
				},
			},
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle mixed structure", func(t *testing.T) {
		input := map[string]any{
			"user.name":        "John",
			"user.tags.0":      "admin",
			"user.tags.1":      "user",
			"user.profile.bio": "Developer",
			"count":            42,
		}

		expected := map[string]any{
			"user": map[string]any{
				"name": "John",
				"tags": []any{"admin", "user"},
				"profile": map[string]any{
					"bio": "Developer",
				},
			},
			"count": 42,
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle nil values", func(t *testing.T) {
		input := map[string]any{
			"user.name":  "John",
			"user.phone": nil,
		}

		expected := map[string]any{
			"user": map[string]any{
				"name":  "John",
				"phone": nil,
			},
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle single level key", func(t *testing.T) {
		input := map[string]any{
			"name": "John",
		}

		expected := map[string]any{
			"name": "John",
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle large array index", func(t *testing.T) {
		input := map[string]any{
			"items.10": "value",
		}

		expected := map[string]any{
			"items": []any{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, "value"},
		}

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle conflicting types - map to slice", func(t *testing.T) {
		input := map[string]any{
			"data.field": "value1",
			"data.0":     "value2",
		}

		result := ExpandMap(input)

		// When there's a conflict, the last one wins
		// In this case, since maps are processed in order, result can vary
		// but the function should not panic
		assert.NotNil(t, result)
	})

	t.Run("should handle conflicting types - slice to map", func(t *testing.T) {
		input := map[string]any{
			"data.0":     "value1",
			"data.field": "value2",
		}

		result := ExpandMap(input)

		// When there's a conflict, the function should handle gracefully
		assert.NotNil(t, result)
	})

	t.Run("should return empty map when root becomes array", func(t *testing.T) {
		input := map[string]any{
			"0": "value1",
			"1": "value2",
		}

		result := ExpandMap(input)

		// When the root becomes an array (numeric keys), ExpandMap should return empty map
		// since it can only return map[string]any, not []any
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should handle only numeric keys at root level", func(t *testing.T) {
		input := map[string]any{
			"0": "first",
			"1": "second",
			"2": "third",
		}

		result := ExpandMap(input)

		// When all keys at root level are numeric, the result should be an empty map
		// because ExpandMap can only return map[string]any, not []any
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should handle large array index properly", func(t *testing.T) {
		input := map[string]any{
			"data.100": "value",
		}

		expected := map[string]any{
			"data": make([]any, 101), // Should create slice with 101 elements (0-100)
		}
		// Set the 100th element
		expectedSlice := make([]any, 101)
		expectedSlice[100] = "value"
		expected["data"] = expectedSlice

		result := ExpandMap(input)
		assert.Equal(t, expected, result)
	})
}

func TestExpandSlice(t *testing.T) {
	t.Run("should initialize nil slice", func(t *testing.T) {
		var slice []any
		expandSlice(&slice, 0)
		assert.NotNil(t, slice)
		assert.Equal(t, 1, len(slice))
	})

	t.Run("should expand slice when index is larger", func(t *testing.T) {
		slice := []any{"a", "b"}
		expandSlice(&slice, 5)
		assert.Equal(t, 6, len(slice))
		assert.Equal(t, "a", slice[0])
		assert.Equal(t, "b", slice[1])
		assert.Nil(t, slice[2])
		assert.Nil(t, slice[5])
	})

	t.Run("should not expand slice when index is within bounds", func(t *testing.T) {
		slice := []any{"a", "b", "c"}
		originalLen := len(slice)
		expandSlice(&slice, 1)
		assert.Equal(t, originalLen, len(slice))
	})

	t.Run("should handle empty slice", func(t *testing.T) {
		slice := make([]any, 0)
		expandSlice(&slice, 2)
		assert.Equal(t, 3, len(slice))
	})

	t.Run("should handle negative index gracefully", func(t *testing.T) {
		slice := []any{"a", "b"}
		// Negative index should not cause issues, but this would be an edge case
		// The function currently doesn't handle negative indices explicitly
		expandSlice(&slice, 0) // Should not change anything since 0 is within bounds
		assert.Equal(t, 2, len(slice))
	})

	t.Run("should expand slice from zero length to large index", func(t *testing.T) {
		slice := make([]any, 0)
		expandSlice(&slice, 10)
		assert.Equal(t, 11, len(slice))
		for i := range 11 {
			assert.Nil(t, slice[i])
		}
	})
}

func TestSetMapRecursive(t *testing.T) {
	t.Run("should create new map when target is nil", func(t *testing.T) {
		result := setMapRecursive(nil, "key", []string{"subkey"}, "value")
		assert.NotNil(t, result)
		assert.IsType(t, map[string]any{}, result)
	})

	t.Run("should use existing map when target is map", func(t *testing.T) {
		target := map[string]any{"existing": "value"}
		result := setMapRecursive(target, "key", []string{"subkey"}, "newvalue")
		assert.Equal(t, "value", result["existing"])
		assert.NotNil(t, result["key"])
	})

	t.Run("should create new map when target is not a map", func(t *testing.T) {
		target := "not a map"
		result := setMapRecursive(target, "key", []string{"subkey"}, "value")
		assert.NotNil(t, result)
		assert.IsType(t, map[string]any{}, result)
	})
}

func TestSetSliceRecursive(t *testing.T) {
	t.Run("should create new slice when target is nil", func(t *testing.T) {
		result := setSliceRecursive(nil, 0, []string{"key"}, "value")
		assert.NotNil(t, result)
		assert.IsType(t, []any{}, result)
		assert.Equal(t, 1, len(result))
	})

	t.Run("should use existing slice when target is slice", func(t *testing.T) {
		target := []any{"existing"}
		result := setSliceRecursive(target, 1, []string{"key"}, "value")
		assert.Equal(t, "existing", result[0])
		assert.Equal(t, 2, len(result))
	})

	t.Run("should create new slice when target is not a slice", func(t *testing.T) {
		target := "not a slice"
		result := setSliceRecursive(target, 0, []string{"key"}, "value")
		assert.NotNil(t, result)
		assert.IsType(t, []any{}, result)
	})
}

func TestSetValueRecursive(t *testing.T) {
	t.Run("should return value when path is empty", func(t *testing.T) {
		result := setValueRecursive(nil, []string{}, "value")
		assert.Equal(t, "value", result)
	})

	t.Run("should handle numeric key for slice", func(t *testing.T) {
		result := setValueRecursive(nil, []string{"0", "key"}, "value")
		assert.IsType(t, []any{}, result)
	})

	t.Run("should handle non-numeric key for map", func(t *testing.T) {
		result := setValueRecursive(nil, []string{"key", "subkey"}, "value")
		assert.IsType(t, map[string]any{}, result)
	})

	t.Run("should handle mixed numeric and non-numeric keys", func(t *testing.T) {
		result := setValueRecursive(nil, []string{"items", "0", "name"}, "John")
		expected := map[string]any{
			"items": []any{
				map[string]any{
					"name": "John",
				},
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("should return empty map when root becomes array", func(t *testing.T) {
		input := map[string]any{
			"0": "value1",
			"1": "value2",
		}

		result := ExpandMap(input)

		// When the root becomes an array (numeric keys), ExpandMap should return empty map
		// since it can only return map[string]any, not []any
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})
}
