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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorktreeManager_GetWorktreeDirs(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) string
		wantCount int
		wantErr   bool
	}{
		{
			name: "finds worktree directories",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create some worktree directories
				dirs := []string{"main", "feature-branch", "review"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
						t.Fatal(err)
					}
				}

				// Create hidden directory (should be ignored)
				if err := os.MkdirAll(filepath.Join(tmpDir, ".hidden"), 0755); err != nil {
					t.Fatal(err)
				}

				// Create a file (should be ignored)
				if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644); err != nil {
					t.Fatal(err)
				}

				return tmpDir
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "handles non-existent directory",
			setup: func(t *testing.T) string {
				return "/non/existent/path"
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := tt.setup(t)
			wm := &WorktreeManager{GitRoot: gitRoot}

			got, err := wm.GetWorktreeDirs()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, got, tt.wantCount)
			}
		})
	}
}

func TestWorktreeManager_GetFilteredWorktrees(t *testing.T) {
	tmpDir := t.TempDir()

	// Create various worktree directories
	dirs := []string{"main", "master", "review", "feature-1", "feature-2", "bugfix"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	wm := &WorktreeManager{GitRoot: tmpDir}

	got, err := wm.GetFilteredWorktrees()
	require.NoError(t, err)

	// Should exclude main, master, review
	expectedCount := 3 // feature-1, feature-2, bugfix
	assert.Len(t, got, expectedCount)

	// Check that main, master, review are not in results
	for _, dir := range got {
		name := filepath.Base(dir)
		assert.NotContains(t, []string{"main", "master", "review"}, name)
	}
}

func TestWorktreeManager_CreateHooks(t *testing.T) {
	tmpDir := t.TempDir()
	wm := &WorktreeManager{GitRoot: tmpDir}

	err := wm.CreateHooks("upstream", "main")
	require.NoError(t, err)

	// Check that hooks directory was created
	hooksDir := wm.GetHooksDir()
	assert.DirExists(t, hooksDir)

	// Check that post-add hook was created and is executable
	hookFile := wm.GetPostAddHook()
	stat, err := os.Stat(hookFile)
	require.NoError(t, err)
	assert.NotZero(t, stat.Mode().Perm()&0111, "Post-add hook should be executable")

	// Check hook content
	content, err := os.ReadFile(hookFile)
	require.NoError(t, err)
	expected := "#!/bin/sh\n\n# Anything here will be ran in the root of a newly created worktree\ngit pull upstream main"
	assert.Equal(t, expected, string(content))
}

func TestWorktreeManager_GetHooksDir(t *testing.T) {
	tmpDir := t.TempDir()
	wm := &WorktreeManager{GitRoot: tmpDir}

	got := wm.GetHooksDir()
	expected := filepath.Join(tmpDir, ".hooks")

	assert.Equal(t, expected, got)
}

func TestWorktreeManager_GetPostAddHook(t *testing.T) {
	tmpDir := t.TempDir()
	wm := &WorktreeManager{GitRoot: tmpDir}

	got := wm.GetPostAddHook()
	expected := filepath.Join(tmpDir, ".hooks", "post-add.sh")

	assert.Equal(t, expected, got)
}

func TestWorktreeManager_SwitchWorktree(t *testing.T) {
	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	wm := &WorktreeManager{GitRoot: tmpDir}

	err := wm.SwitchWorktree(tmpDir)
	require.NoError(t, err)

	currentDir, _ := os.Getwd()
	// On macOS, temp dirs might be symlinked, so compare resolved paths
	expectedResolved, _ := filepath.EvalSymlinks(tmpDir)
	currentResolved, _ := filepath.EvalSymlinks(currentDir)
	assert.Equal(t, expectedResolved, currentResolved)
}
