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

package cmd

import (
	"fmt"
	"sort"

	"github.com/liamawhite/worktree/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage worktree configuration",
	Long:  `Configure domain-to-account mappings for different Git hosting providers.`,
}

var setAccountCmd = &cobra.Command{
	Use:   "set-account <domain> <account>",
	Short: "Set account name for a domain",
	Long: `Set the account name to use for a specific domain.

Examples:
  wt config set-account github.com myusername
  wt config set-account enterprise.github.com john.doe
  wt config set-account gitlab.company.com jdoe`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		account := args[1]

		cfg, err := LoadConfigWithOverride()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfg.SetAccount(domain, account)

		if err := SaveConfigWithOverride(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Set account for %s to %s\n", domain, account)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured hosts with their full configuration",
	Long:  `List all configured hosts showing both account and clone method.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfigWithOverride()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		hosts := cfg.ListHosts()

		if len(hosts) == 0 {
			fmt.Println("No hosts configured")
			return nil
		}

		fmt.Println("Configured hosts:")

		// Sort domains for consistent output
		var domains []string
		for domain := range hosts {
			domains = append(domains, domain)
		}
		sort.Strings(domains)

		for _, domain := range domains {
			hostConfig := hosts[domain]
			cloneMethod := hostConfig.CloneMethod
			if cloneMethod == "" {
				cloneMethod = config.CloneMethodHTTP // Default display
			}
			fmt.Printf("  %s: %s (clone: %s)\n", domain, hostConfig.Account, cloneMethod)
		}

		return nil
	},
}

var setCloneMethodCmd = &cobra.Command{
	Use:   "set-clone-method <domain> <method>",
	Short: "Set clone method for a domain",
	Long: `Set the clone method (http or ssh) to use for a specific domain.

Examples:
  wt config set-clone-method github.com ssh
  wt config set-clone-method enterprise.github.com http
  wt config set-clone-method gitlab.company.com ssh`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		methodStr := args[1]

		cfg, err := LoadConfigWithOverride()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		method, err := config.ParseCloneMethod(methodStr)
		if err != nil {
			return err
		}

		cfg.SetCloneMethod(domain, method)

		if err := SaveConfigWithOverride(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Set clone method for %s to %s\n", domain, method)
		return nil
	},
}

func init() {
	configCmd.AddCommand(setAccountCmd)
	configCmd.AddCommand(listCmd)
	configCmd.AddCommand(setCloneMethodCmd)
}
