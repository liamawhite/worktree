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

	"github.com/liamawhite/worktree/pkg/selector"
	"github.com/liamawhite/worktree/pkg/worktree"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:     "switch",
	Aliases: []string{"sw"},
	Short:   "Switch to a different worktree",
	Long:    `Interactively select and switch to a different worktree.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm, err := worktree.NewWorktreeManager()
		if err != nil {
			return fmt.Errorf("must be in a git repository to switch worktrees: %w", err)
		}

		worktreeDirs, err := wm.GetWorktreeDirs()
		if err != nil {
			return err
		}

		if len(worktreeDirs) == 0 {
			fmt.Println("No worktrees available")
			return nil
		}

		selectedWorktree, err := selector.Select("Select a worktree to switch to:", worktreeDirs)
		if err != nil {
			return err
		}

		if selectedWorktree == "" {
			fmt.Println("No worktree selected, staying where we are")
			return nil
		}

		if err := wm.SwitchWorktree(selectedWorktree); err != nil {
			return fmt.Errorf("failed to switch to worktree: %w", err)
		}

		fmt.Printf("Switched to worktree: %s\n", selectedWorktree)
		return nil
	},
}
