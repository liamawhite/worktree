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
	"os"

	"github.com/liamawhite/worktree/pkg/setup"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [domain/]org/repo",
	Short: "Setup a new worktree repository",
	Long: `Setup a new repository with worktrees. Supports both GitHub.com and GitHub Enterprise.
Clones the repository and configures upstream/origin remotes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]
		branch, _ := cmd.Flags().GetString("base")

		config, err := setup.ParseRepoString(repo, branch)
		if err != nil {
			return err
		}

		err = setup.SetupRepository(config, getConfigPath())
		if err != nil {
			return err
		}

		// Output directory change signal for auto-cd functionality
		fmt.Fprintf(os.Stderr, "WT_CHDIR:%s\n", config.RepoName)
		return nil
	},
}

func init() {
	setupCmd.Flags().StringP("base", "b", "main", "Base branch to use for the repository")
}
