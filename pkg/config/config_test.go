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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Accounts)
	assert.Equal(t, "liamawhite", cfg.Accounts["github.com"])
}

func TestGetDefaultConfigPath(t *testing.T) {
	path, err := GetDefaultConfigPath()
	require.NoError(t, err)

	assert.Contains(t, path, ".config/worktree/settings.yaml")
}

func TestConfig_GetAccount(t *testing.T) {
	cfg := &Config{
		Accounts: map[string]string{
			"github.com":            "testuser",
			"enterprise.github.com": "corpuser",
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

func TestConfig_SetAccount(t *testing.T) {
	cfg := &Config{}

	// Test setting on nil accounts map
	cfg.SetAccount("github.com", "newuser")
	assert.Equal(t, "newuser", cfg.Accounts["github.com"])

	// Test setting on existing map
	cfg.SetAccount("enterprise.github.com", "corpuser")
	assert.Equal(t, "corpuser", cfg.Accounts["enterprise.github.com"])

	// Test empty domain defaults to github.com
	cfg.SetAccount("", "defaultuser")
	assert.Equal(t, "defaultuser", cfg.Accounts["github.com"])
}

func TestConfig_ListAccounts(t *testing.T) {
	cfg := &Config{
		Accounts: map[string]string{
			"github.com":            "testuser",
			"enterprise.github.com": "corpuser",
		},
	}

	accounts := cfg.ListAccounts()

	expected := map[string]string{
		"github.com":            "testuser",
		"enterprise.github.com": "corpuser",
	}

	assert.Equal(t, expected, accounts)

	// Verify it returns a copy (modifying returned map shouldn't affect original)
	accounts["new.domain.com"] = "newuser"
	assert.NotContains(t, cfg.Accounts, "new.domain.com")
}

func TestConfig_ListAccounts_NilMap(t *testing.T) {
	cfg := &Config{}

	accounts := cfg.ListAccounts()
	assert.NotNil(t, accounts)
	assert.Empty(t, accounts)
}

func TestLoadDefaultConfig(t *testing.T) {
	// Test that LoadDefaultConfig works
	cfg, err := LoadDefaultConfig()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Should have the default account
	assert.Equal(t, "liamawhite", cfg.GetAccount("github.com"))
}

func TestLoadConfigFromPath_SaveToPath(t *testing.T) {
	// Use a temporary directory for this test
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "settings.yaml")

	// Test loading config when file doesn't exist (should create default)
	cfg, err := LoadConfigFromPath(configPath)
	require.NoError(t, err)
	assert.Equal(t, "liamawhite", cfg.GetAccount("github.com"))

	// Verify file was created
	assert.FileExists(t, configPath)

	// Modify and save config
	cfg.SetAccount("enterprise.github.com", "corpuser")
	err = cfg.SaveToPath(configPath)
	require.NoError(t, err)

	// Load config again and verify changes persisted
	cfg2, err := LoadConfigFromPath(configPath)
	require.NoError(t, err)
	assert.Equal(t, "liamawhite", cfg2.GetAccount("github.com"))
	assert.Equal(t, "corpuser", cfg2.GetAccount("enterprise.github.com"))
}
