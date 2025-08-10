package cmd

import (
	"fmt"

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
		return wm.ClearWorktrees()
	},
}

