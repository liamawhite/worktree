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
	"testing"

	"github.com/stretchr/testify/require"
)

const testBinary = "wt-test"

// Framework handles the common setup logic for integration tests
type Framework struct {
	BinaryPath  string
	ProjectRoot string
	TempDir     string
	OriginalDir string
	ConfigPath  string
	t           *testing.T
}

// NewFramework creates a new test framework instance and builds the binary
func NewFramework(t *testing.T) *Framework {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	framework := &Framework{t: t}

	// Get current directory and project root
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	framework.ProjectRoot = filepath.Dir(currentDir)

	// Build the binary
	framework.buildBinary()

	// Set up temporary workspace
	framework.setupWorkspace()

	return framework
}

// buildBinary builds the worktree binary for testing
func (f *Framework) buildBinary() {
	f.t.Log("Building worktree binary...")

	f.BinaryPath = filepath.Join(f.ProjectRoot, testBinary)

	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir = f.ProjectRoot
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(f.t, err, "Failed to build worktree binary: %s", string(buildOutput))
}

// setupWorkspace creates and configures the temporary workspace
func (f *Framework) setupWorkspace() {
	// Create temporary workspace
	f.TempDir = f.t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(f.t, err)
	f.OriginalDir = originalDir

	// Change to temp directory for test
	err = os.Chdir(f.TempDir)
	require.NoError(f.t, err)

	// Set up config path
	f.ConfigPath = filepath.Join(f.TempDir, "test-config.yaml")

	f.t.Logf("Running test in directory: %s", f.TempDir)
	f.t.Logf("Using binary: %s", f.BinaryPath)
	f.t.Logf("Using config: %s", f.ConfigPath)
}

// Cleanup performs cleanup operations for the test
func (f *Framework) Cleanup() {
	// Change back to original directory
	if f.OriginalDir != "" {
		os.Chdir(f.OriginalDir)
	}

	// Remove binary
	if f.BinaryPath != "" {
		os.Remove(f.BinaryPath)
	}
}

// RunCommand executes a worktree command with the test config
func (f *Framework) RunCommand(args ...string) ([]byte, error) {
	fullArgs := append([]string{"--config", f.ConfigPath}, args...)
	cmd := exec.Command(f.BinaryPath, fullArgs...)
	cmd.Env = append(os.Environ(), "WORKTREE_CONFIG="+f.ConfigPath)
	return cmd.CombinedOutput()
}

// RunCommandInDir executes a worktree command in a specific directory
func (f *Framework) RunCommandInDir(dir string, args ...string) ([]byte, error) {
	fullArgs := append([]string{"--config", f.ConfigPath}, args...)
	cmd := exec.Command(f.BinaryPath, fullArgs...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "WORKTREE_CONFIG="+f.ConfigPath)
	return cmd.CombinedOutput()
}

// SetupAccount configures the test account
func (f *Framework) SetupAccount(domain, account string) {
	f.t.Logf("Setting up account %s for %s...", account, domain)
	output, err := f.RunCommand("config", "set-account", domain, account)
	require.NoError(f.t, err, "Failed to set account: %s", string(output))
	f.t.Logf("Set account output: %s", string(output))
}

// VerifyAccount verifies the account configuration
func (f *Framework) VerifyAccount(domain, expectedAccount string) {
	output, err := f.RunCommand("config", "get-account", domain)
	require.NoError(f.t, err, "Failed to get account: %s", string(output))
	require.Contains(f.t, string(output), expectedAccount)
	f.t.Logf("Get account output: %s", string(output))
}
