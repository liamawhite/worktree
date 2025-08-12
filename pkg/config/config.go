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
	"strings"

	"gopkg.in/yaml.v3"
)

// CloneMethod represents the method used for cloning repositories
type CloneMethod string

const (
	// CloneMethodHTTP uses HTTPS URLs for cloning
	CloneMethodHTTP CloneMethod = "http"
	// CloneMethodSSH uses SSH URLs for cloning  
	CloneMethodSSH CloneMethod = "ssh"
)

// String returns the string representation of the clone method
func (c CloneMethod) String() string {
	return string(c)
}

// IsValid checks if the clone method is valid
func (c CloneMethod) IsValid() bool {
	return c == CloneMethodHTTP || c == CloneMethodSSH
}

// ParseCloneMethod parses a string into a CloneMethod
func ParseCloneMethod(s string) (CloneMethod, error) {
	method := CloneMethod(strings.ToLower(s))
	if !method.IsValid() {
		return "", fmt.Errorf("invalid clone method: %s (valid options: http, ssh)", s)
	}
	return method, nil
}

// HostConfig represents configuration for a specific host/domain
type HostConfig struct {
	Account     string      `yaml:"account"`
	CloneMethod CloneMethod `yaml:"clone_method,omitempty"`
}

// Config represents the application configuration
type Config struct {
	// Legacy field for backward compatibility - will be migrated to Hosts
	Accounts map[string]string    `yaml:"accounts,omitempty"`
	// New field for host configurations
	Hosts    map[string]HostConfig `yaml:"hosts,omitempty"`
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Hosts: map[string]HostConfig{},
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

	// Migrate legacy config format to new format
	if err := config.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate config: %w", err)
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

// migrate migrates legacy config format to new format
func (c *Config) migrate() error {
	// If we have legacy accounts, migrate them to hosts
	if len(c.Accounts) > 0 {
		if c.Hosts == nil {
			c.Hosts = make(map[string]HostConfig)
		}
		
		// Only migrate accounts that don't already exist in hosts
		for domain, account := range c.Accounts {
			if _, exists := c.Hosts[domain]; !exists {
				c.Hosts[domain] = HostConfig{
					Account:     account,
					CloneMethod: CloneMethodHTTP, // Default to HTTP for existing configs
				}
			}
		}
		// Clear legacy accounts after migration
		c.Accounts = nil
	}

	// Initialize hosts map if nil
	if c.Hosts == nil {
		c.Hosts = make(map[string]HostConfig)
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

	if config, exists := c.Hosts[domain]; exists {
		return config.Account
	}

	return ""
}

// SetAccount sets the account name for the given domain
func (c *Config) SetAccount(domain, account string) {
	// Handle empty domain (default to github.com)
	if domain == "" {
		domain = "github.com"
	}

	if c.Hosts == nil {
		c.Hosts = make(map[string]HostConfig)
	}

	// Preserve existing clone method if it exists, otherwise default to HTTP
	existing := c.Hosts[domain]
	cloneMethod := existing.CloneMethod
	if cloneMethod == "" {
		cloneMethod = CloneMethodHTTP
	}

	c.Hosts[domain] = HostConfig{
		Account:     account,
		CloneMethod: cloneMethod,
	}
}

// ListAccounts returns all configured domain-account pairs
func (c *Config) ListAccounts() map[string]string {
	if c.Hosts == nil {
		return make(map[string]string)
	}

	// Return a copy to prevent external modification
	result := make(map[string]string)
	for domain, config := range c.Hosts {
		result[domain] = config.Account
	}

	return result
}

// GetHostConfig returns the full host configuration for the given domain
func (c *Config) GetHostConfig(domain string) HostConfig {
	// Handle empty domain (default to github.com)
	if domain == "" {
		domain = "github.com"
	}

	if config, exists := c.Hosts[domain]; exists {
		return config
	}

	// Return default config if not found
	return HostConfig{
		Account:     "",
		CloneMethod: CloneMethodHTTP,
	}
}

// GetCloneMethod returns the clone method for the given domain
func (c *Config) GetCloneMethod(domain string) CloneMethod {
	config := c.GetHostConfig(domain)
	if config.CloneMethod == "" {
		return CloneMethodHTTP // Default to HTTP
	}
	return config.CloneMethod
}

// SetCloneMethod sets the clone method for the given domain
func (c *Config) SetCloneMethod(domain string, method CloneMethod) {
	// Handle empty domain (default to github.com)
	if domain == "" {
		domain = "github.com"
	}

	if c.Hosts == nil {
		c.Hosts = make(map[string]HostConfig)
	}

	// Preserve existing account if it exists
	existing := c.Hosts[domain]
	account := existing.Account

	c.Hosts[domain] = HostConfig{
		Account:     account,
		CloneMethod: method,
	}
}

// ListHosts returns all configured hosts with their full configuration
func (c *Config) ListHosts() map[string]HostConfig {
	if c.Hosts == nil {
		return make(map[string]HostConfig)
	}

	// Return a copy to prevent external modification
	result := make(map[string]HostConfig)
	for domain, config := range c.Hosts {
		result[domain] = config
	}

	return result
}

// GenerateRepositoryURL generates the appropriate repository URL based on the clone method
func (c *Config) GenerateRepositoryURL(domain, org, repo string) string {
	cloneMethod := c.GetCloneMethod(domain)
	
	switch cloneMethod {
	case CloneMethodSSH:
		return fmt.Sprintf("git@%s:%s/%s.git", domain, org, repo)
	case CloneMethodHTTP:
		fallthrough
	default:
		return fmt.Sprintf("https://%s/%s/%s.git", domain, org, repo)
	}
}

// GenerateUserRepositoryURL generates a repository URL for a user's fork
func (c *Config) GenerateUserRepositoryURL(domain, repo string) string {
	account := c.GetAccount(domain)
	if account == "" {
		// Fallback to HTTP with empty account - this will likely fail but preserves existing behavior
		return fmt.Sprintf("https://%s//%s.git", domain, repo)
	}
	
	cloneMethod := c.GetCloneMethod(domain)
	
	switch cloneMethod {
	case CloneMethodSSH:
		return fmt.Sprintf("git@%s:%s/%s.git", domain, account, repo)
	case CloneMethodHTTP:
		fallthrough
	default:
		return fmt.Sprintf("https://%s/%s/%s.git", domain, account, repo)
	}
}
