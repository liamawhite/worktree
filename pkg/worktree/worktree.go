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

package worktree

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/liamawhite/worktree/pkg/git"
)

//go:embed templates/post-add.sh
var postAddHookTemplate string

type WorktreeManager struct {
	GitRoot string
}

func NewWorktreeManager() (*WorktreeManager, error) {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return nil, err
	}
	return &WorktreeManager{GitRoot: gitRoot}, nil
}

func (wm *WorktreeManager) GetWorktreeDirs() ([]string, error) {
	entries, err := os.ReadDir(wm.GitRoot)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			dirs = append(dirs, filepath.Join(wm.GitRoot, entry.Name()))
		}
	}
	return dirs, nil
}

func (wm *WorktreeManager) GetFilteredWorktrees() ([]string, error) {
	dirs, err := wm.GetWorktreeDirs()
	if err != nil {
		return nil, err
	}

	var filtered []string
	for _, dir := range dirs {
		name := filepath.Base(dir)
		if name != "main" && name != "master" && name != "review" {
			filtered = append(filtered, dir)
		}
	}
	return filtered, nil
}

func (wm *WorktreeManager) GetHooksDir() string {
	return filepath.Join(wm.GitRoot, ".hooks")
}

func (wm *WorktreeManager) GetPostAddHook() string {
	return filepath.Join(wm.GetHooksDir(), "post-add.sh")
}

func (wm *WorktreeManager) CreateHooks(base, branch string) error {
	hooksDir := wm.GetHooksDir()
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return err
	}

	postAddHook := wm.GetPostAddHook()

	// Parse and execute template
	tmpl, err := template.New("post-add").Parse(postAddHookTemplate)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(postAddHook, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	data := struct {
		Base   string
		Branch string
	}{
		Base:   base,
		Branch: branch,
	}

	return tmpl.Execute(file, data)
}

func (wm *WorktreeManager) RunPostAddHook(worktreePath string) error {
	hookPath := wm.GetPostAddHook()
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return nil
	}

	return git.RunCommandInDir(worktreePath, "sh", hookPath)
}

func (wm *WorktreeManager) AddWorktree(branch, base string) error {
	if base == "" {
		base = "main"
	}

	if err := os.Chdir(wm.GitRoot); err != nil {
		return err
	}

	if err := git.RunGitCommand("worktree", "add", "-b", branch, branch, base); err != nil {
		return err
	}

	worktreePath := filepath.Join(wm.GitRoot, branch)
	if err := wm.RunPostAddHook(worktreePath); err != nil {
		return fmt.Errorf("failed to run post-add hook: %w", err)
	}

	return nil
}

func (wm *WorktreeManager) RemoveWorktree(worktreePath string) error {
	worktreeName := filepath.Base(worktreePath)

	if err := git.RunGitCommand("worktree", "remove", worktreeName, "--force"); err != nil {
		return err
	}

	if err := git.DeleteBranch(wm.GitRoot, worktreeName); err != nil {
		return err
	}

	return nil
}

func (wm *WorktreeManager) ClearWorktrees() error {
	if err := os.Chdir(wm.GitRoot); err != nil {
		return err
	}

	filteredWorktrees, err := wm.GetFilteredWorktrees()
	if err != nil {
		return err
	}

	if len(filteredWorktrees) == 0 {
		return nil
	}

	for _, worktreePath := range filteredWorktrees {
		worktree := filepath.Base(worktreePath)
		fmt.Printf("Removing worktree: %s\n", worktree)

		if err := git.RunGitCommand("worktree", "remove", worktree, "--force"); err != nil {
			fmt.Printf("Warning: failed to remove worktree %s: %v\n", worktree, err)
			continue
		}

		if err := git.DeleteBranch(wm.GitRoot, worktree); err != nil {
			fmt.Printf("Warning: failed to delete branch %s: %v\n", worktree, err)
		}
	}

	return nil
}

func (wm *WorktreeManager) SwitchWorktree(worktreePath string) error {
	return os.Chdir(worktreePath)
}
