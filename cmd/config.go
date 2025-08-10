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

var getAccountCmd = &cobra.Command{
	Use:   "get-account <domain>",
	Short: "Get account name for a domain",
	Long: `Get the configured account name for a specific domain.

Examples:
  wt config get-account github.com
  wt config get-account enterprise.github.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		cfg, err := LoadConfigWithOverride()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		account := cfg.GetAccount(domain)
		if account == "" {
			fmt.Printf("No account configured for domain: %s\n", domain)
			return nil
		}

		fmt.Printf("%s: %s\n", domain, account)
		return nil
	},
}

var listAccountsCmd = &cobra.Command{
	Use:   "list-accounts",
	Short: "List all configured domain-account mappings",
	Long:  `List all configured domain-to-account mappings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfigWithOverride()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		accounts := cfg.ListAccounts()

		if len(accounts) == 0 {
			fmt.Println("No accounts configured")
			return nil
		}

		fmt.Println("Configured accounts:")

		// Sort domains for consistent output
		var domains []string
		for domain := range accounts {
			domains = append(domains, domain)
		}
		sort.Strings(domains)

		for _, domain := range domains {
			fmt.Printf("  %s: %s\n", domain, accounts[domain])
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(setAccountCmd)
	configCmd.AddCommand(getAccountCmd)
	configCmd.AddCommand(listAccountsCmd)
}
