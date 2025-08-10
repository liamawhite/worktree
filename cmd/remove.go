package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/liamawhite/worktree/pkg/selector"
	"github.com/liamawhite/worktree/pkg/worktree"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm [worktree-name]",
	Aliases: []string{"d", "delete", "del", "remove"},
	Short:   "Remove a worktree",
	Long:    `Remove a worktree by name, or interactively select one if no name is provided (excluding main/master and review).`,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wm, err := worktree.NewWorktreeManager()
		if err != nil {
			return fmt.Errorf("must be in a git repository to remove worktrees: %w", err)
		}

		var selectedWorktree string

		if len(args) > 0 {
			worktreeName := args[0]
			allWorktrees, err := wm.GetWorktreeDirs()
			if err != nil {
				return err
			}

			for _, wt := range allWorktrees {
				if filepath.Base(wt) == worktreeName {
					selectedWorktree = wt
					break
				}
			}

			if selectedWorktree == "" {
				return fmt.Errorf("worktree '%s' not found", worktreeName)
			}
		} else {
			filteredWorktrees, err := wm.GetFilteredWorktrees()
			if err != nil {
				return err
			}

			if len(filteredWorktrees) == 0 {
				fmt.Println("No worktrees available to remove")
				return nil
			}

			selectedWorktree, err = selector.Select("Select a worktree to remove:", filteredWorktrees)
			if err != nil {
				return err
			}

			if selectedWorktree == "" {
				fmt.Println("No worktree selected, no action taken")
				return nil
			}
		}

		worktreeName := filepath.Base(selectedWorktree)
		fmt.Printf("Removing worktree: %s\n", worktreeName)

		return wm.RemoveWorktree(selectedWorktree)
	},
}