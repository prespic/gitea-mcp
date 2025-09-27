package operation

import (
	"testing"
)

func TestParseBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantOK    bool
	}{
		{
			name:      "valid token",
			header:    "Bearer validtoken",
			wantToken: "validtoken",
			wantOK:    true,
		},
		{
			name:      "token with spaces trimmed",
			header:    "Bearer   spacedToken ",
			wantToken: "spacedToken",
			wantOK:    true,
		},
		{
			name:      "lowercase bearer should fail",
			header:    "bearer lowercase",
			wantToken: "",
			wantOK:    false,
		},
		{
			name:      "bearer with no token",
			header:    "Bearer ",
			wantToken: "",
			wantOK:    false,
		},
		{
			name:      "bearer with only spaces",
			header:    "Bearer     ",
			wantToken: "",
			wantOK:    false,
		},
		{
			name:      "missing space after Bearer",
			header:    "Bearertoken",
			wantToken: "",
			wantOK:    false,
		},
		{
			name:      "different auth type",
			header:    "Basic dXNlcjpwYXNz",
			wantToken: "",
			wantOK:    false,
		},
		{
			name:      "empty header",
			header:    "",
			wantToken: "",
			wantOK:    false,
		},
		{
			name:      "token with internal spaces",
			header:    "Bearer token with spaces",
			wantToken: "token with spaces",
			wantOK:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotOK := parseBearerToken(tt.header)
			if gotToken != tt.wantToken {
				t.Errorf("parseBearerToken() token = %q, want %q", gotToken, tt.wantToken)
			}
			if gotOK != tt.wantOK {
				t.Errorf("parseBearerToken() ok = %v, want %v", gotOK, tt.wantOK)
			}
		})
	}
}
