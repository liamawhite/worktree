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

package setup

import (
	"fmt"
	"os"
	"strings"

	"github.com/liamawhite/worktree/pkg/config"
	"github.com/liamawhite/worktree/pkg/git"
	"github.com/liamawhite/worktree/pkg/worktree"
)

type RepoConfig struct {
	Domain   string
	Org      string
	RepoName string
	Branch   string
}

func ParseRepoString(repo, defaultBranch string) (*RepoConfig, error) {
	parts := strings.Split(repo, "/")

	config := &RepoConfig{
		Branch: defaultBranch,
	}

	switch len(parts) {
	case 3: // domain/org/repo (GHE or explicit github.com)
		config.Domain = parts[0]
		config.Org = parts[1]
		config.RepoName = parts[2]
	case 2: // org/repo (assume GitHub.com)
		config.Domain = "github.com"
		config.Org = parts[0]
		config.RepoName = parts[1]
	default:
		return nil, fmt.Errorf("invalid repository format. Expected [domain/]org/repo")
	}

	return config, nil
}

func (rc *RepoConfig) IsGitHubEnterprise() bool {
	return rc.Domain != "github.com"
}

func SetupRepository(config *RepoConfig) error {
	if config.IsGitHubEnterprise() {
		return setupGHERepo(config)
	}
	return setupGitHubRepo(config)
}

func setupGHERepo(repoConfig *RepoConfig) error {
	// Load configuration to get account name for this domain
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	account := cfg.GetAccount(repoConfig.Domain)
	if account == "" {
		// For GHE, we'll clone directly from the original repository if no account is configured
		fmt.Printf("No account configured for %s, cloning directly from %s/%s\n", repoConfig.Domain, repoConfig.Org, repoConfig.RepoName)
		return setupDirectCloneGHE(repoConfig)
	}

	fmt.Printf("Cloning forked %s repository from %s and hiding .git internals\n", repoConfig.RepoName, repoConfig.Domain)

	if err := os.MkdirAll(repoConfig.RepoName, 0755); err != nil {
		return err
	}

	if err := os.Chdir(repoConfig.RepoName); err != nil {
		return err
	}

	// Use the configured account for the fork
	repoURL := fmt.Sprintf("https://%s/%s/%s.git", repoConfig.Domain, account, repoConfig.RepoName)
	if err := git.CloneBare(repoURL, ".bare"); err != nil {
		return err
	}

	if err := createGitDirFile(); err != nil {
		return err
	}

	// Add upstream remote pointing to the original repository
	fmt.Println("Adding upstream remote")
	upstreamURL := fmt.Sprintf("https://%s/%s/%s.git", repoConfig.Domain, repoConfig.Org, repoConfig.RepoName)
	if err := git.AddRemote(".bare", "upstream", upstreamURL); err != nil {
		return err
	}

	return finishSetup("upstream", repoConfig.Branch)
}

// setupDirectCloneGHE clones directly from the original GHE repository
func setupDirectCloneGHE(repoConfig *RepoConfig) error {
	fmt.Printf("Cloning %s/%s/%s repository directly and hiding .git internals\n", repoConfig.Domain, repoConfig.Org, repoConfig.RepoName)

	if err := os.MkdirAll(repoConfig.RepoName, 0755); err != nil {
		return err
	}

	if err := os.Chdir(repoConfig.RepoName); err != nil {
		return err
	}

	repoURL := fmt.Sprintf("https://%s/%s/%s.git", repoConfig.Domain, repoConfig.Org, repoConfig.RepoName)
	if err := git.CloneBare(repoURL, ".bare"); err != nil {
		return err
	}

	if err := createGitDirFile(); err != nil {
		return err
	}

	return finishSetup("origin", repoConfig.Branch)
}

func setupGitHubRepo(repoConfig *RepoConfig) error {
	// Load configuration to get account name
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	account := cfg.GetAccount(repoConfig.Domain)
	if account == "" {
		return fmt.Errorf("no account configured for %s. Use 'wt config set-account %s <username>' to configure", repoConfig.Domain, repoConfig.Domain)
	}

	fmt.Printf("Cloning %s repository and configuring remotes\n", repoConfig.RepoName)

	if err := os.MkdirAll(repoConfig.RepoName, 0755); err != nil {
		return err
	}

	if err := os.Chdir(repoConfig.RepoName); err != nil {
		return err
	}

	// Clone from the original repository
	originURL := fmt.Sprintf("https://%s/%s/%s.git", repoConfig.Domain, repoConfig.Org, repoConfig.RepoName)
	if err := git.CloneBare(originURL, ".bare"); err != nil {
		return err
	}

	if err := createGitDirFile(); err != nil {
		return err
	}

	// If the account is different from the original org, add the fork as a remote
	if account != repoConfig.Org {
		fmt.Printf("Adding %s remote for your fork\n", account)
		forkURL := fmt.Sprintf("https://%s/%s/%s.git", repoConfig.Domain, account, repoConfig.RepoName)
		if err := git.AddRemote(".bare", account, forkURL); err != nil {
			return err
		}
	}

	return finishSetup("origin", repoConfig.Branch)
}


func createGitDirFile() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitdirContent := fmt.Sprintf("gitdir: %s/.bare", cwd)
	return os.WriteFile(".git", []byte(gitdirContent), 0644)
}

func finishSetup(base, branch string) error {
	fmt.Println("Creating worktree hooks")
	wm := &worktree.WorktreeManager{GitRoot: "."}
	if err := wm.CreateHooks(base, branch); err != nil {
		return err
	}

	fmt.Printf("Creating worktree for base branch %s\n", branch)
	if err := git.RunGitCommand("worktree", "add", branch); err != nil {
		return err
	}

	fmt.Println("Creating worktree for a review branch")
	if err := git.RunGitCommand("worktree", "add", "review", "--force"); err != nil {
		return err
	}

	return nil
}
