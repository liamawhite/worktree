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

//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRepo   = "liamawhite/worktree"
	testBranch = "main"
)

func TestWorktreeFullWorkflow(t *testing.T) {
	// Step 1: Set up test environment
	framework := NewFramework(t)
	defer framework.Cleanup()

	// Step 2: Set up configuration with test account
	t.Log("Setting up configuration...")
	framework.SetupAccount("github.com", "test")
	framework.VerifyAccount("github.com", "test")

	// Step 3: Setup repository (without forking, direct clone)
	t.Log("Setting up repository...")
	output, err := framework.RunCommand("setup", testRepo, testBranch)
	require.NoError(t, err, "Failed to setup repository: %s", string(output))
	t.Logf("Setup output: %s", string(output))

	// Verify repository structure was created
	assert.DirExists(t, "worktree", "Repository directory should exist")

	// Change into repository directory
	err = os.Chdir("worktree")
	require.NoError(t, err)

	// Verify basic repository structure
	assert.DirExists(t, ".bare", "Bare repository should exist")
	assert.FileExists(t, ".git", "Git directory file should exist")
	assert.DirExists(t, ".hooks", "Hooks directory should exist")
	assert.DirExists(t, testBranch, "Main branch worktree should exist")
	assert.DirExists(t, "review", "Review worktree should exist")

	// Step 4: Add test worktrees
	t.Log("Adding test worktrees...")

	// Add test1 worktree
	output, err = framework.RunCommandInDir("worktree", "add", "test1", "--base", testBranch)
	require.NoError(t, err, "Failed to add test1 worktree: %s", string(output))
	t.Logf("Add test1 output: %s", string(output))

	// Verify test1 was created
	assert.DirExists(t, "worktree/test1", "test1 worktree should exist")

	// Add test2 worktree
	output, err = framework.RunCommandInDir("worktree", "add", "test2", "--base", testBranch)
	require.NoError(t, err, "Failed to add test2 worktree: %s", string(output))
	t.Logf("Add test2 output: %s", string(output))

	// Verify test2 was created
	assert.DirExists(t, "worktree/test2", "test2 worktree should exist")

	// Step 5: Verify all worktrees exist
	t.Log("Verifying all worktrees exist...")
	entries, err := os.ReadDir("worktree")
	require.NoError(t, err)

	var worktreeDirs []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			worktreeDirs = append(worktreeDirs, entry.Name())
		}
	}

	expectedDirs := []string{testBranch, "review", "test1", "test2"}
	for _, expected := range expectedDirs {
		assert.Contains(t, worktreeDirs, expected, "Expected worktree %s should exist", expected)
	}
	t.Logf("Found worktree directories: %v", worktreeDirs)

	// Step 6: Remove test1 worktree using direct command
	t.Log("Removing test1 worktree...")
	output, err = framework.RunCommandInDir("worktree", "rm", "test1")
	require.NoError(t, err, "Failed to remove test1 worktree: %s", string(output))
	t.Logf("Remove test1 output: %s", string(output))

	// Verify test1 was removed
	assert.NoDirExists(t, "worktree/test1", "test1 worktree should be removed")
	assert.DirExists(t, "worktree/test2", "test2 worktree should still exist")

	// Step 7: Clear all worktrees except main and review
	t.Log("Clearing all worktrees...")
	output, err = framework.RunCommandInDir("worktree", "clear")
	require.NoError(t, err, "Failed to clear worktrees: %s", string(output))
	t.Logf("Clear output: %s", string(output))

	// Step 8: Verify final state - only main and review should remain
	t.Log("Verifying final state...")
	entries, err = os.ReadDir("worktree")
	require.NoError(t, err)

	worktreeDirs = []string{}
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			worktreeDirs = append(worktreeDirs, entry.Name())
		}
	}

	expectedFinalDirs := []string{testBranch, "review"}
	assert.Len(t, worktreeDirs, len(expectedFinalDirs), "Should only have main and review worktrees")

	for _, expected := range expectedFinalDirs {
		assert.Contains(t, worktreeDirs, expected, "Expected final worktree %s should exist", expected)
	}

	// Ensure test worktrees are gone
	assert.NoDirExists(t, "worktree/test1", "test1 should be cleaned up")
	assert.NoDirExists(t, "worktree/test2", "test2 should be cleaned up")

	t.Logf("Final worktree directories: %v", worktreeDirs)
	t.Log("✅ Integration test completed successfully!")
}

func TestConfigOverrides(t *testing.T) {
	framework := NewFramework(t)
	defer framework.Cleanup()

	// Test CLI flag override
	t.Run("CLI Flag Override", func(t *testing.T) {
		output, err := framework.RunCommand("config", "set-account", "github.com", "cli-user")
		require.NoError(t, err, "CLI config failed: %s", string(output))

		output, err = framework.RunCommand("config", "get-account", "github.com")
		require.NoError(t, err)
		assert.Contains(t, string(output), "cli-user")
	})

	// Test environment variable override
	t.Run("Environment Variable Override", func(t *testing.T) {
		// Create a separate config file for env test
		envConfigPath := filepath.Join(framework.TempDir, "env-config.yaml")

		cmd := exec.Command(framework.BinaryPath, "config", "set-account", "github.com", "env-user")
		cmd.Env = append(os.Environ(), "WORKTREE_CONFIG="+envConfigPath)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Env config failed: %s", string(output))

		cmd = exec.Command(framework.BinaryPath, "config", "get-account", "github.com")
		cmd.Env = append(os.Environ(), "WORKTREE_CONFIG="+envConfigPath)
		output, err = cmd.CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(output), "env-user")
	})
}
