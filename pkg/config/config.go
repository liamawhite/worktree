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
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Accounts map[string]string `yaml:"accounts"`
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Accounts: map[string]string{
			"github.com": "liamawhite", // Default account for backward compatibility
		},
	}
}

// GetDefaultConfigPath returns the default config file path
func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "worktree", "settings.yaml"), nil
}

// LoadConfigFromPath loads configuration from the specified path
func LoadConfigFromPath(configPath string) (*Config, error) {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// If config file doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultCfg := DefaultConfig()
		if err := defaultCfg.SaveToPath(configPath); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		return defaultCfg, nil
	}

	// Read existing config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Initialize accounts map if nil
	if config.Accounts == nil {
		config.Accounts = make(map[string]string)
	}

	return &config, nil
}

// SaveToPath persists the configuration to the specified path
func (c *Config) SaveToPath(configPath string) error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetAccount returns the account name for the given domain
// Returns empty string if not found
func (c *Config) GetAccount(domain string) string {
	// Handle empty domain (default to github.com)
	if domain == "" {
		domain = "github.com"
	}

	return c.Accounts[domain]
}

// SetAccount sets the account name for the given domain
func (c *Config) SetAccount(domain, account string) {
	// Handle empty domain (default to github.com)
	if domain == "" {
		domain = "github.com"
	}

	if c.Accounts == nil {
		c.Accounts = make(map[string]string)
	}

	c.Accounts[domain] = account
}

// ListAccounts returns all configured domain-account pairs
func (c *Config) ListAccounts() map[string]string {
	if c.Accounts == nil {
		return make(map[string]string)
	}

	// Return a copy to prevent external modification
	result := make(map[string]string)
	for domain, account := range c.Accounts {
		result[domain] = account
	}

	return result
}
