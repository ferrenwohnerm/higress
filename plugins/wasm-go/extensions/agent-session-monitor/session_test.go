package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeSessionID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal session id",
			input:    "abc-123-xyz",
			expected: "abc-123-xyz",
		},
		{
			name:     "trims whitespace",
			input:    "  session-456  ",
			expected: "session-456",
		},
		{
			name:     "truncates long id",
			input:    string(make([]byte, 200)),
			expected: string(make([]byte, 128)),
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeSessionID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  PluginConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: PluginConfig{
				SessionIDHeader:       "X-Session-ID",
				MaxSessionDurationSec: 3600,
				EnableMetrics:         true,
			},
			wantErr: false,
		},
		{
			name: "negative max duration",
			config: PluginConfig{
				SessionIDHeader:       "X-Session-ID",
				MaxSessionDurationSec: -1,
			},
			wantErr: true,
		},
		{
			name: "zero duration allowed",
			config: PluginConfig{
				SessionIDHeader:       "X-Session-ID",
				MaxSessionDurationSec: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
