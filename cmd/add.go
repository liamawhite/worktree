package cmd

import (
	"fmt"

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
		return wm.AddWorktree(branch, base)
	},
}

func init() {
	addCmd.Flags().StringP("base", "b", "main", "Base branch to create the new worktree from")
}