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

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	require.NoError(t, err)
	
	assert.Contains(t, path, ".config/worktree/settings.yaml")
}

func TestConfig_GetAccount(t *testing.T) {
	cfg := &Config{
		Accounts: map[string]string{
			"github.com":           "testuser",
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
			"github.com":           "testuser",
			"enterprise.github.com": "corpuser",
		},
	}
	
	accounts := cfg.ListAccounts()
	
	expected := map[string]string{
		"github.com":           "testuser",
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

func TestGetConfigPathWithOverride(t *testing.T) {
	tests := []struct {
		name     string
		override string
		wantPath string
	}{
		{
			name:     "with override",
			override: "/custom/path/config.yaml",
			wantPath: "/custom/path/config.yaml",
		},
		{
			name:     "empty override uses default",
			override: "",
			wantPath: "", // Will be default path, we just check it's not empty
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetConfigPathWithOverride(tt.override)
			require.NoError(t, err)
			
			if tt.override != "" {
				assert.Equal(t, tt.wantPath, result)
			} else {
				assert.NotEmpty(t, result)
				assert.Contains(t, result, "settings.yaml")
			}
		})
	}
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