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

	"github.com/liamawhite/worktree/pkg/worktree"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <worktree-name>",
	Short: "Add a new worktree",
	Long:  `Create a new worktree with a new branch based on the specified base branch.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
		base, _ := cmd.Flags().GetString("base")

		wm, err := worktree.NewWorktreeManager()
		if err != nil {
			return err
		}

		fmt.Printf("Creating worktree and branch: %s from base: %s\n", branch, base)
		if err := wm.AddWorktree(branch, base); err != nil {
			return err
		}

		worktreePath := filepath.Join(wm.GitRoot, branch)
		fmt.Printf("WT_CHDIR:%s\n", worktreePath)
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("base", "b", "main", "Base branch to create the new worktree from")
}
