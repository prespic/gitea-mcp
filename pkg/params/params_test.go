package params

import (
	"strings"
	"testing"
)

func TestGetIndex(t *testing.T) {
	tests := []struct {
		name      string
		args      map[string]interface{}
		key       string
		wantIndex int64
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid float64",
			args:      map[string]interface{}{"index": float64(123)},
			key:       "index",
			wantIndex: 123,
			wantErr:   false,
		},
		{
			name:      "valid string",
			args:      map[string]interface{}{"index": "456"},
			key:       "index",
			wantIndex: 456,
			wantErr:   false,
		},
		{
			name:      "valid string with large number",
			args:      map[string]interface{}{"index": "999999"},
			key:       "index",
			wantIndex: 999999,
			wantErr:   false,
		},
		{
			name:    "missing parameter",
			args:    map[string]interface{}{},
			key:     "index",
			wantErr: true,
			errMsg:  "index is required",
		},
		{
			name:    "invalid string (not a number)",
			args:    map[string]interface{}{"index": "abc"},
			key:     "index",
			wantErr: true,
			errMsg:  "must be a valid integer",
		},
		{
			name:    "invalid string (decimal)",
			args:    map[string]interface{}{"index": "12.34"},
			key:     "index",
			wantErr: true,
			errMsg:  "must be a valid integer",
		},
		{
			name:    "invalid type (bool)",
			args:    map[string]interface{}{"index": true},
			key:     "index",
			wantErr: true,
			errMsg:  "must be a number or numeric string",
		},
		{
			name:    "invalid type (map)",
			args:    map[string]interface{}{"index": map[string]string{"foo": "bar"}},
			key:     "index",
			wantErr: true,
			errMsg:  "must be a number or numeric string",
		},
		{
			name:      "custom key name",
			args:      map[string]interface{}{"pr_index": "789"},
			key:       "pr_index",
			wantIndex: 789,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, err := GetIndex(tt.args, tt.key)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetIndex() expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("GetIndex() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("GetIndex() unexpected error = %v", err)
				return
			}
			if gotIndex != tt.wantIndex {
				t.Errorf("GetIndex() = %v, want %v", gotIndex, tt.wantIndex)
			}
		})
	}
}
