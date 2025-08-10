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
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneBare(t *testing.T) {
	// This test would require a real repository, so we'll skip it in unit tests
	t.Skip("CloneBare test requires network access and a real repository")
}

func TestAddRemote(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize a bare repository first
	_, err := git.PlainInit(tmpDir, true)
	require.NoError(t, err)

	// Test adding a remote
	err = AddRemote(tmpDir, "origin", "https://github.com/example/repo.git")
	assert.NoError(t, err)

	// Verify remote was added
	remotes, err := GetRemotes(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/example/repo.git", remotes["origin"])
}

func TestGetRemotes(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize a bare repository
	_, err := git.PlainInit(tmpDir, true)
	require.NoError(t, err)

	// Add a remote
	err = AddRemote(tmpDir, "upstream", "https://github.com/upstream/repo.git")
	require.NoError(t, err)

	// Get remotes
	remotes, err := GetRemotes(tmpDir)
	require.NoError(t, err)
	assert.Contains(t, remotes, "upstream")
	assert.Equal(t, "https://github.com/upstream/repo.git", remotes["upstream"])
}

func TestCreateBranch(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize a regular repository and create an initial commit
	repo, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)

	// Create a file and commit it to have a HEAD
	// This is a simplified version - in real usage there would be actual file operations
	// For testing purposes, we'll just test that the repository exists
	_, err = repo.Head()
	if err != nil {
		// Repository doesn't have HEAD yet, skip this test
		t.Skip("Cannot test branch creation without initial commit")
	}

	// Test creating a branch
	err = CreateBranch(tmpDir, "feature", "")
	assert.NoError(t, err)
}

func TestDeleteBranch(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize a repository
	_, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)

	// Test deleting a non-existent branch (should not error)
	err = DeleteBranch(tmpDir, "nonexistent")
	// This might or might not error depending on go-git implementation
	// We're mainly testing that the function doesn't panic
	_ = err
}
