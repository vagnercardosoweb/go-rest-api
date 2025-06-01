package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenMap(t *testing.T) {
	t.Run("should return nil when input is nil", func(t *testing.T) {
		result := FlattenMap(nil)
		assert.Nil(t, result)
	})

	t.Run("should return empty map when input is empty", func(t *testing.T) {
		result := FlattenMap(make(map[string]any))
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should flatten simple map", func(t *testing.T) {
		input := map[string]any{
			"name": "John",
			"age":  30,
		}

		expected := map[string]any{
			"name": "John",
			"age":  30,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should flatten nested map", func(t *testing.T) {
		input := map[string]any{
			"user": map[string]any{
				"name": "John",
				"age":  30,
			},
		}

		expected := map[string]any{
			"user.name": "John",
			"user.age":  30,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should flatten deeply nested map", func(t *testing.T) {
		input := map[string]any{
			"user": map[string]any{
				"profile": map[string]any{
					"name": "John",
					"age":  30,
				},
			},
		}

		expected := map[string]any{
			"user.profile.name": "John",
			"user.profile.age":  30,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should flatten array", func(t *testing.T) {
		input := map[string]any{
			"tags": []any{"go", "javascript", "python"},
		}

		expected := map[string]any{
			"tags.0": "go",
			"tags.1": "javascript",
			"tags.2": "python",
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should flatten empty array", func(t *testing.T) {
		input := map[string]any{
			"tags": []any{},
		}

		expected := map[string]any{
			"tags": []any{},
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should flatten empty nested map", func(t *testing.T) {
		input := map[string]any{
			"user": map[string]any{},
		}

		expected := map[string]any{
			"user": map[string]any{},
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should flatten array with nested objects", func(t *testing.T) {
		input := map[string]any{
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

		expected := map[string]any{
			"users.0.name": "John",
			"users.0.age":  30,
			"users.1.name": "Jane",
			"users.1.age":  25,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle nil values", func(t *testing.T) {
		input := map[string]any{
			"name":  "John",
			"phone": nil,
		}

		expected := map[string]any{
			"name":  "John",
			"phone": nil,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle nil value in nested structure", func(t *testing.T) {
		input := map[string]any{
			"user": map[string]any{
				"name":  "John",
				"phone": nil,
			},
		}

		expected := map[string]any{
			"user.name":  "John",
			"user.phone": nil,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle complex mixed structure", func(t *testing.T) {
		input := map[string]any{
			"user": map[string]any{
				"name": "John",
				"tags": []any{"admin", "user"},
				"profile": map[string]any{
					"bio": "Developer",
				},
			},
			"count": 42,
		}

		expected := map[string]any{
			"user.name":        "John",
			"user.tags.0":      "admin",
			"user.tags.1":      "user",
			"user.profile.bio": "Developer",
			"count":            42,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle primitive values", func(t *testing.T) {
		input := map[string]any{
			"string": "test",
			"int":    42,
			"float":  3.14,
			"bool":   true,
			"nil":    nil,
		}

		expected := map[string]any{
			"string": "test",
			"int":    42,
			"float":  3.14,
			"bool":   true,
			"nil":    nil,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle deeply nested nil values", func(t *testing.T) {
		input := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": nil,
				},
			},
		}

		expected := map[string]any{
			"level1.level2.level3": nil,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle array with nil values", func(t *testing.T) {
		input := map[string]any{
			"items": []any{nil, "value", nil},
		}

		expected := map[string]any{
			"items.0": nil,
			"items.1": "value",
			"items.2": nil,
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})

	t.Run("should handle mixed structure with all edge cases", func(t *testing.T) {
		input := map[string]any{
			"empty_map":   map[string]any{},
			"empty_array": []any{},
			"nil_value":   nil,
			"nested": map[string]any{
				"array_with_nil": []any{nil, map[string]any{"key": "value"}},
				"deeply_nested": map[string]any{
					"level": map[string]any{
						"final": "value",
					},
				},
			},
		}

		expected := map[string]any{
			"empty_map":                        map[string]any{},
			"empty_array":                      []any{},
			"nil_value":                        nil,
			"nested.array_with_nil.0":          nil,
			"nested.array_with_nil.1.key":      "value",
			"nested.deeply_nested.level.final": "value",
		}

		result := FlattenMap(input)
		assert.Equal(t, expected, result)
	})
}

func TestBuildPrefixedKey(t *testing.T) {
	t.Run("should return key when prefix is empty", func(t *testing.T) {
		result := buildPrefixedKey("", "name")
		assert.Equal(t, "name", result)
	})

	t.Run("should return prefixed key when prefix is not empty", func(t *testing.T) {
		result := buildPrefixedKey("user", "name")
		assert.Equal(t, "user.name", result)
	})

	t.Run("should handle nested prefixes", func(t *testing.T) {
		result := buildPrefixedKey("user.profile", "name")
		assert.Equal(t, "user.profile.name", result)
	})
}
