package cmd

import (
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

		return setup.SetupRepository(config)
	},
}

func init() {
	setupCmd.Flags().StringP("base", "b", "main", "Base branch to use for the repository")
}

