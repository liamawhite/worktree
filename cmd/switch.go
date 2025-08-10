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