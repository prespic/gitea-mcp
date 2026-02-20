package params

import (
	"fmt"
	"strconv"
)

// GetIndex extracts an index parameter from MCP tool arguments.
// It accepts both numeric (float64 from JSON) and string representations.
// This provides better UX for LLM callers that may naturally use strings
// for identifiers like issue/PR numbers.
func GetIndex(args map[string]interface{}, key string) (int64, error) {
	val, exists := args[key]
	if !exists {
		return 0, fmt.Errorf("%s is required", key)
	}

	// Try float64 (JSON number type)
	if f, ok := val.(float64); ok {
		return int64(f), nil
	}

	// Try string and parse to integer
	if s, ok := val.(string); ok {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid integer (got %q)", key, s)
		}
		return i, nil
	}

	return 0, fmt.Errorf("%s must be a number or numeric string", key)
}
