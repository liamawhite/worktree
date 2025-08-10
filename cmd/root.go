package cmd

import (
	"fmt"
	"os"

	"github.com/liamawhite/worktree/pkg/config"
	"github.com/spf13/cobra"
)

// Global variable to store config path override
var globalConfigPath string

// getDefaultConfigPath returns the default config path
func getDefaultConfigPath() string {
	path, err := config.GetDefaultConfigPath()
	if err != nil {
		return ""
	}
	return path
}

// getConfigPath returns the config path to use, considering env var and flag overrides
func getConfigPath() string {
	// Priority 1: Flag override
	if globalConfigPath != "" {
		return globalConfigPath
	}
	
	// Priority 2: Environment variable override
	if envPath := os.Getenv("WORKTREE_CONFIG"); envPath != "" {
		return envPath
	}
	
	// Priority 3: Default location
	return getDefaultConfigPath()
}

var RootCmd = &cobra.Command{
	Use:     "wt",
	Aliases: []string{"worktree"},
	Short:   "Git worktree management tool",
	Long: `A CLI tool for managing Git worktrees with support for GitHub forks,
enterprise Git hosting, and interactive worktree switching.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set global config path if provided via flag
		if configPath, _ := cmd.Flags().GetString("config"); configPath != "" {
			globalConfigPath = configPath
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Add global config flag with explicit default
	defaultPath := getDefaultConfigPath()
	RootCmd.PersistentFlags().StringP("config", "c", defaultPath, "config file path")

	RootCmd.AddCommand(setupCmd)
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(rmCmd)
	RootCmd.AddCommand(clearCmd)
	RootCmd.AddCommand(switchCmd)
	RootCmd.AddCommand(configCmd)
}

// LoadConfigWithOverride loads config using the resolved config path
func LoadConfigWithOverride() (*config.Config, error) {
	return config.LoadConfigFromPath(getConfigPath())
}

// SaveConfigWithOverride saves config using the resolved config path
func SaveConfigWithOverride(cfg *config.Config) error {
	return cfg.SaveToPath(getConfigPath())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

