package params

import (
	"fmt"
	"strconv"
)

// GetString extracts a required string parameter from MCP tool arguments.
func GetString(args map[string]any, key string) (string, error) {
	val, ok := args[key].(string)
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}
	return val, nil
}

// GetOptionalString extracts an optional string parameter with a default value.
func GetOptionalString(args map[string]any, key, defaultVal string) string {
	if val, ok := args[key].(string); ok {
		return val
	}
	return defaultVal
}

// GetStringSlice extracts an optional string slice parameter from MCP tool arguments.
func GetStringSlice(args map[string]any, key string) []string {
	val, ok := args[key]
	if !ok {
		return nil
	}
	sliceVal, ok := val.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(sliceVal))
	for _, item := range sliceVal {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// GetPagination extracts page and pageSize parameters, returning them as ints.
func GetPagination(args map[string]any, defaultPageSize int64) (page, pageSize int) {
	return int(GetOptionalInt(args, "page", 1)), int(GetOptionalInt(args, "pageSize", defaultPageSize))
}

// ToInt64 converts a value to int64, accepting both float64 (JSON number) and
// string representations. Returns false if the value cannot be converted.
func ToInt64(val any) (int64, bool) {
	switch v := val.(type) {
	case float64:
		return int64(v), true
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}

// GetIndex extracts a required integer parameter from MCP tool arguments.
// It accepts both numeric (float64 from JSON) and string representations.
// This provides better UX for LLM callers that may naturally use strings
// for identifiers like issue/PR numbers.
func GetIndex(args map[string]any, key string) (int64, error) {
	val, exists := args[key]
	if !exists {
		return 0, fmt.Errorf("%s is required", key)
	}

	if i, ok := ToInt64(val); ok {
		return i, nil
	}

	if s, ok := val.(string); ok {
		return 0, fmt.Errorf("%s must be a valid integer (got %q)", key, s)
	}

	return 0, fmt.Errorf("%s must be a number or numeric string", key)
}

// GetOptionalInt extracts an optional integer parameter from MCP tool arguments.
// Returns defaultVal if the key is missing or the value cannot be parsed.
// Accepts both float64 (JSON number) and string representations.
func GetOptionalInt(args map[string]any, key string, defaultVal int64) int64 {
	val, exists := args[key]
	if !exists {
		return defaultVal
	}
	if i, ok := ToInt64(val); ok {
		return i
	}
	return defaultVal
}
