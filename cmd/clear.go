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

	"github.com/liamawhite/worktree/pkg/worktree"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all worktrees except main, master and review",
	Long:  `Remove all worktrees except main, master and review branches.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm, err := worktree.NewWorktreeManager()
		if err != nil {
			return fmt.Errorf("must be in a git repository to clear worktrees: %w", err)
		}

		filteredWorktrees, err := wm.GetFilteredWorktrees()
		if err != nil {
			return err
		}

		if len(filteredWorktrees) == 0 {
			fmt.Println("No worktrees to clear")
			return nil
		}

		fmt.Println("Removing all worktrees except main, master and review")
		needsChdir, err := wm.ClearWorktrees()
		if err != nil {
			return err
		}

		if needsChdir {
			fmt.Fprintf(os.Stderr, "WT_CHDIR:%s\n", wm.GitRoot)
		}

		return nil
	},
}
