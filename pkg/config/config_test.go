// Copyright 2025 Liam White
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneMethod(t *testing.T) {
	t.Run("String method", func(t *testing.T) {
		assert.Equal(t, "http", CloneMethodHTTP.String())
		assert.Equal(t, "ssh", CloneMethodSSH.String())
	})

	t.Run("IsValid method", func(t *testing.T) {
		assert.True(t, CloneMethodHTTP.IsValid())
		assert.True(t, CloneMethodSSH.IsValid())
		assert.False(t, CloneMethod("invalid").IsValid())
	})

	t.Run("ParseCloneMethod", func(t *testing.T) {
		tests := []struct {
			input    string
			expected CloneMethod
			hasError bool
		}{
			{"http", CloneMethodHTTP, false},
			{"ssh", CloneMethodSSH, false},
			{"HTTP", CloneMethodHTTP, false}, // case insensitive
			{"SSH", CloneMethodSSH, false},   // case insensitive
			{"invalid", "", true},
			{"", "", true},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				result, err := ParseCloneMethod(tt.input)
				if tt.hasError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, result)
				}
			})
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Hosts)
	assert.Empty(t, cfg.Hosts)
}

func TestGetDefaultConfigPath(t *testing.T) {
	path, err := GetDefaultConfigPath()
	require.NoError(t, err)

	assert.Contains(t, path, ".config/worktree/settings.yaml")
}

func TestConfig_GetAccount(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]HostConfig{
			"github.com": {
				Account:     "testuser",
				CloneMethod: CloneMethodHTTP,
			},
			"enterprise.github.com": {
				Account:     "corpuser",
				CloneMethod: CloneMethodSSH,
			},
		},
	}

	tests := []struct {
		name     string
		domain   string
		expected string
	}{
		{
			name:     "existing domain",
			domain:   "github.com",
			expected: "testuser",
		},
		{
			name:     "existing enterprise domain",
			domain:   "enterprise.github.com",
			expected: "corpuser",
		},
		{
			name:     "empty domain defaults to github.com",
			domain:   "",
			expected: "testuser",
		},
		{
			name:     "non-existing domain",
			domain:   "gitlab.com",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.GetAccount(tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_GetCloneMethod(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]HostConfig{
			"github.com": {
				Account:     "testuser",
				CloneMethod: CloneMethodSSH,
			},
			"enterprise.github.com": {
				Account:     "corpuser",
				CloneMethod: CloneMethodHTTP,
			},
		},
	}

	tests := []struct {
		name     string
		domain   string
		expected CloneMethod
	}{
		{
			name:     "existing domain with SSH",
			domain:   "github.com",
			expected: CloneMethodSSH,
		},
		{
			name:     "existing domain with HTTP",
			domain:   "enterprise.github.com",
			expected: CloneMethodHTTP,
		},
		{
			name:     "empty domain defaults to github.com",
			domain:   "",
			expected: CloneMethodSSH,
		},
		{
			name:     "non-existing domain defaults to HTTP",
			domain:   "gitlab.com",
			expected: CloneMethodHTTP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.GetCloneMethod(tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_GetHostConfig(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]HostConfig{
			"github.com": {
				Account:     "testuser",
				CloneMethod: CloneMethodSSH,
			},
		},
	}

	t.Run("existing host", func(t *testing.T) {
		result := cfg.GetHostConfig("github.com")
		expected := HostConfig{
			Account:     "testuser",
			CloneMethod: CloneMethodSSH,
		}
		assert.Equal(t, expected, result)
	})

	t.Run("non-existing host", func(t *testing.T) {
		result := cfg.GetHostConfig("gitlab.com")
		expected := HostConfig{
			Account:     "",
			CloneMethod: CloneMethodHTTP,
		}
		assert.Equal(t, expected, result)
	})
}

func TestConfig_SetAccount(t *testing.T) {
	cfg := &Config{}

	cfg.SetAccount("github.com", "newuser")
	assert.Equal(t, "newuser", cfg.GetAccount("github.com"))
	assert.Equal(t, CloneMethodHTTP, cfg.GetCloneMethod("github.com"))

	// Setting account should preserve existing clone method
	cfg.SetCloneMethod("github.com", CloneMethodSSH)
	cfg.SetAccount("github.com", "anotheruser")
	assert.Equal(t, "anotheruser", cfg.GetAccount("github.com"))
	assert.Equal(t, CloneMethodSSH, cfg.GetCloneMethod("github.com"))
}

func TestConfig_SetCloneMethod(t *testing.T) {
	cfg := &Config{}

	cfg.SetCloneMethod("github.com", CloneMethodSSH)
	assert.Equal(t, CloneMethodSSH, cfg.GetCloneMethod("github.com"))
	assert.Equal(t, "", cfg.GetAccount("github.com"))

	// Setting clone method should preserve existing account
	cfg.SetAccount("github.com", "testuser")
	cfg.SetCloneMethod("github.com", CloneMethodHTTP)
	assert.Equal(t, "testuser", cfg.GetAccount("github.com"))
	assert.Equal(t, CloneMethodHTTP, cfg.GetCloneMethod("github.com"))
}

func TestConfig_ListAccounts(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]HostConfig{
			"github.com": {
				Account:     "testuser",
				CloneMethod: CloneMethodHTTP,
			},
			"enterprise.github.com": {
				Account:     "corpuser",
				CloneMethod: CloneMethodSSH,
			},
		},
	}

	result := cfg.ListAccounts()
	expected := map[string]string{
		"github.com":            "testuser",
		"enterprise.github.com": "corpuser",
	}

	assert.Equal(t, expected, result)
}

func TestConfig_ListHosts(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]HostConfig{
			"github.com": {
				Account:     "testuser",
				CloneMethod: CloneMethodHTTP,
			},
			"enterprise.github.com": {
				Account:     "corpuser",
				CloneMethod: CloneMethodSSH,
			},
		},
	}

	result := cfg.ListHosts()
	expected := map[string]HostConfig{
		"github.com": {
			Account:     "testuser",
			CloneMethod: CloneMethodHTTP,
		},
		"enterprise.github.com": {
			Account:     "corpuser",
			CloneMethod: CloneMethodSSH,
		},
	}

	assert.Equal(t, expected, result)
}

func TestConfig_GenerateRepositoryURL(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]HostConfig{
			"github.com": {
				Account:     "testuser",
				CloneMethod: CloneMethodHTTP,
			},
			"enterprise.github.com": {
				Account:     "corpuser",
				CloneMethod: CloneMethodSSH,
			},
		},
	}

	tests := []struct {
		name     string
		domain   string
		org      string
		repo     string
		expected string
	}{
		{
			name:     "GitHub HTTP",
			domain:   "github.com",
			org:      "myorg",
			repo:     "myrepo",
			expected: "https://github.com/myorg/myrepo.git",
		},
		{
			name:     "Enterprise SSH",
			domain:   "enterprise.github.com",
			org:      "myorg",
			repo:     "myrepo",
			expected: "git@enterprise.github.com:myorg/myrepo.git",
		},
		{
			name:     "Unknown domain defaults to HTTP",
			domain:   "gitlab.com",
			org:      "myorg",
			repo:     "myrepo",
			expected: "https://gitlab.com/myorg/myrepo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.GenerateRepositoryURL(tt.domain, tt.org, tt.repo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_GenerateUserRepositoryURL(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]HostConfig{
			"github.com": {
				Account:     "testuser",
				CloneMethod: CloneMethodHTTP,
			},
			"enterprise.github.com": {
				Account:     "corpuser",
				CloneMethod: CloneMethodSSH,
			},
		},
	}

	tests := []struct {
		name     string
		domain   string
		repo     string
		expected string
	}{
		{
			name:     "GitHub HTTP",
			domain:   "github.com",
			repo:     "myrepo",
			expected: "https://github.com/testuser/myrepo.git",
		},
		{
			name:     "Enterprise SSH",
			domain:   "enterprise.github.com",
			repo:     "myrepo",
			expected: "git@enterprise.github.com:corpuser/myrepo.git",
		},
		{
			name:     "Unknown domain with no account",
			domain:   "gitlab.com",
			repo:     "myrepo",
			expected: "https://gitlab.com//myrepo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.GenerateUserRepositoryURL(tt.domain, tt.repo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_Migration(t *testing.T) {
	t.Run("migrates legacy accounts to hosts", func(t *testing.T) {
		cfg := &Config{
			Accounts: map[string]string{
				"github.com":            "testuser",
				"enterprise.github.com": "corpuser",
			},
		}

		err := cfg.migrate()
		require.NoError(t, err)

		// Check that accounts were migrated
		assert.Nil(t, cfg.Accounts)
		assert.Equal(t, "testuser", cfg.GetAccount("github.com"))
		assert.Equal(t, "corpuser", cfg.GetAccount("enterprise.github.com"))
		assert.Equal(t, CloneMethodHTTP, cfg.GetCloneMethod("github.com"))
		assert.Equal(t, CloneMethodHTTP, cfg.GetCloneMethod("enterprise.github.com"))
	})

	t.Run("preserves existing hosts when accounts exist", func(t *testing.T) {
		cfg := &Config{
			Accounts: map[string]string{
				"github.com": "legacyuser",
			},
			Hosts: map[string]HostConfig{
				"enterprise.github.com": {
					Account:     "corpuser",
					CloneMethod: CloneMethodSSH,
				},
			},
		}

		err := cfg.migrate()
		require.NoError(t, err)

		// Legacy accounts should be migrated only if not already in hosts
		assert.Nil(t, cfg.Accounts)
		assert.Equal(t, "legacyuser", cfg.GetAccount("github.com"))
		assert.Equal(t, CloneMethodHTTP, cfg.GetCloneMethod("github.com"))
		assert.Equal(t, "corpuser", cfg.GetAccount("enterprise.github.com"))
		assert.Equal(t, CloneMethodSSH, cfg.GetCloneMethod("enterprise.github.com"))
	})

	t.Run("initializes empty hosts map when both are empty", func(t *testing.T) {
		cfg := &Config{}

		err := cfg.migrate()
		require.NoError(t, err)

		assert.NotNil(t, cfg.Hosts)
		assert.Empty(t, cfg.Hosts)
	})
}

func TestLoadConfigFromPath_WithMigration(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "settings.yaml")

	// Create a legacy config file
	legacyConfigContent := `accounts:
  github.com: legacyuser
  enterprise.github.com: corpuser`

	err := WriteConfigContent(configPath, legacyConfigContent)
	require.NoError(t, err)

	// Load the config - should migrate automatically
	cfg, err := LoadConfigFromPath(configPath)
	require.NoError(t, err)

	// Verify migration occurred
	assert.Equal(t, "legacyuser", cfg.GetAccount("github.com"))
	assert.Equal(t, "corpuser", cfg.GetAccount("enterprise.github.com"))
	assert.Equal(t, CloneMethodHTTP, cfg.GetCloneMethod("github.com"))
	assert.Equal(t, CloneMethodHTTP, cfg.GetCloneMethod("enterprise.github.com"))
}

// Helper function for tests
func WriteConfigContent(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}