package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSocks5String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected socks5Config
		errMsg   string
	}{
		{
			name:  "socks5 with host and port only",
			input: "socks5:banana.island.com:8080",
			expected: socks5Config{
				address:  "banana.island.com:8080",
				username: "",
				password: "",
			},
		},
		{
			name:  "socks5:// with host and port only",
			input: "socks5://mango.island.com:9090",
			expected: socks5Config{
				address:  "mango.island.com:9090",
				username: "",
				password: "",
			},
		},
		{
			name:  "socks5:// with colon-separated credentials",
			input: "socks5://apple.fiji.net:1234:tiger123:secretpaw",
			expected: socks5Config{
				address:  "apple.fiji.net:1234",
				username: "tiger123",
				password: "secretpaw",
			},
		},
		{
			name:  "socks5 with colon-separated credentials",
			input: "socks5:orange.bali.org:5678:elephant789:junglepass",
			expected: socks5Config{
				address:  "orange.bali.org:5678",
				username: "elephant789",
				password: "junglepass",
			},
		},
		{
			name:  "socks5:// with @ format",
			input: "socks5://dolphin456:oceanwave@grape.hawaii.net:3333",
			expected: socks5Config{
				address:  "grape.hawaii.net:3333",
				username: "dolphin456",
				password: "oceanwave",
			},
		},
		{
			name:  "host:port with credentials",
			input: "kiwi.newzealand.co.nz:8080:kiwiuser:kiwipass",
			expected: socks5Config{
				address:  "kiwi.newzealand.co.nz:8080",
				username: "kiwiuser",
				password: "kiwipass",
			},
		},
		{
			name:  "host:port without credentials",
			input: "kiwi.newzealand.co.nz:8080",
			expected: socks5Config{
				address:  "kiwi.newzealand.co.nz:8080",
				username: "",
				password: "",
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: socks5Config{},
			errMsg:   "empty input string",
		},
		{
			name:     "invalid format - missing port",
			input:    "socks5:hostname",
			expected: socks5Config{},
			errMsg:   "invalid format",
		},
		{
			name:     "invalid @ format",
			input:    "socks5://user:pass@host:invalid",
			expected: socks5Config{},
			errMsg:   "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseSocks5String(tt.input)

			if tt.errMsg != "" {
				require.EqualError(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected.address, result.address)
			require.Equal(t, tt.expected.username, result.username)
			require.Equal(t, tt.expected.password, result.password)
		})
	}
}
