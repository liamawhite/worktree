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

package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func FindGitRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check for bare repository setup
	if _, err := os.Stat(filepath.Join(cwd, ".bare")); err == nil {
		return cwd, nil
	}

	// Use go-git to find repository root
	_, err = git.PlainOpenWithOptions(cwd, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}

	// Walk up directories to find the root
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return filepath.Dir(dir), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("git repository root not found")
}

func RunGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunGitCommandInDir(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunGitCommandOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CloneBare clones a repository as a bare repository using go-git
func CloneBare(url, path string) error {
	cloneOptions := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	// If this is an SSH URL, configure SSH authentication
	if strings.HasPrefix(url, "git@") {
		auth, err := getSSHAuth()
		if err != nil {
			return fmt.Errorf("failed to configure SSH authentication: %w", err)
		}
		cloneOptions.Auth = auth
	}

	_, err := git.PlainClone(path, true, cloneOptions)
	return err
}

// getSSHAuth configures SSH authentication using the SSH agent or default key locations
func getSSHAuth() (transport.AuthMethod, error) {
	// Try to use SSH agent first
	auth, err := ssh.NewSSHAgentAuth("git")
	if err == nil {
		return auth, nil
	}

	// Fallback to default SSH key locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Try common SSH key locations
	keyPaths := []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
	}

	for _, keyPath := range keyPaths {
		if _, err := os.Stat(keyPath); err == nil {
			auth, err := ssh.NewPublicKeysFromFile("git", keyPath, "")
			if err == nil {
				return auth, nil
			}
		}
	}

	return nil, fmt.Errorf("no SSH authentication method available (tried SSH agent and common key locations)")
}

// OpenRepository opens a git repository using go-git
func OpenRepository(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

// AddRemote adds a remote to a repository using go-git
func AddRemote(repoPath, remoteName, url string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})
	return err
}

// GetRemotes returns all remotes for a repository
func GetRemotes(repoPath string) (map[string]string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, remote := range remotes {
		cfg := remote.Config()
		if len(cfg.URLs) > 0 {
			result[cfg.Name] = cfg.URLs[0]
		}
	}
	return result, nil
}

// CreateBranch creates a new branch using go-git
func CreateBranch(repoPath, branchName, baseBranch string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	// Get the base branch reference
	var baseRef *plumbing.Reference
	if baseBranch != "" {
		baseRef, err = repo.Reference(plumbing.ReferenceName("refs/heads/"+baseBranch), true)
		if err != nil {
			return fmt.Errorf("base branch %s not found: %w", baseBranch, err)
		}
	} else {
		// Get HEAD
		headRef, err := repo.Head()
		if err != nil {
			return err
		}
		baseRef = headRef
	}

	// Create new branch reference
	newRef := plumbing.NewHashReference(
		plumbing.ReferenceName("refs/heads/"+branchName),
		baseRef.Hash(),
	)

	return repo.Storer.SetReference(newRef)
}

// DeleteBranch deletes a branch using go-git
func DeleteBranch(repoPath, branchName string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	return repo.Storer.RemoveReference(plumbing.ReferenceName("refs/heads/" + branchName))
}
