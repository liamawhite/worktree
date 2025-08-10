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
	"path/filepath"

	"github.com/liamawhite/worktree/pkg/selector"
	"github.com/liamawhite/worktree/pkg/worktree"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:     "switch [worktree]",
	Aliases: []string{"sw"},
	Short:   "Switch to a different worktree",
	Long:    `Switch to a different worktree. If no worktree is specified, interactively select one.`,
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

		var selectedWorktree string

		if len(args) > 0 {
			targetWorktree := args[0]
			for _, dir := range worktreeDirs {
				if filepath.Base(dir) == targetWorktree {
					selectedWorktree = dir
					break
				}
			}
			if selectedWorktree == "" {
				return fmt.Errorf("worktree '%s' not found", targetWorktree)
			}
		} else {
			selectedWorktree, err = selector.Select("Select a worktree to switch to:", worktreeDirs)
			if err != nil {
				return err
			}

			if selectedWorktree == "" {
				fmt.Println("No worktree selected, staying where we are")
				return nil
			}
		}

		if err := wm.SwitchWorktree(selectedWorktree); err != nil {
			return fmt.Errorf("failed to switch to worktree: %w", err)
		}

		fmt.Printf("Switched to worktree: %s\n", selectedWorktree)
		fmt.Printf("WT_CHDIR:%s\n", selectedWorktree)
		return nil
	},
}
