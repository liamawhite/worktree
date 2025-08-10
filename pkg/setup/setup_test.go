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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRepoString(t *testing.T) {
	tests := []struct {
		name          string
		repo          string
		defaultBranch string
		want          *RepoConfig
		wantErr       bool
	}{
		{
			name:          "GitHub.com repo",
			repo:          "owner/repo",
			defaultBranch: "main",
			want: &RepoConfig{
				Domain:   "github.com",
				Org:      "owner",
				RepoName: "repo",
				Branch:   "main",
			},
			wantErr: false,
		},
		{
			name:          "GHE repo",
			repo:          "github.company.com/owner/repo",
			defaultBranch: "main",
			want: &RepoConfig{
				Domain:   "github.company.com",
				Org:      "owner",
				RepoName: "repo",
				Branch:   "main",
			},
			wantErr: false,
		},
		{
			name:          "invalid format - single part",
			repo:          "justname",
			defaultBranch: "main",
			want:          nil,
			wantErr:       true,
		},
		{
			name:          "invalid format - too many parts",
			repo:          "a/b/c/d",
			defaultBranch: "main",
			want:          nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepoString(tt.repo, tt.defaultBranch)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRepoConfig_IsGitHubEnterprise(t *testing.T) {
	tests := []struct {
		name   string
		config *RepoConfig
		want   bool
	}{
		{
			name: "GitHub.com",
			config: &RepoConfig{
				Domain:   "github.com",
				Org:      "owner",
				RepoName: "repo",
			},
			want: false,
		},
		{
			name: "GHE",
			config: &RepoConfig{
				Domain:   "github.company.com",
				Org:      "owner",
				RepoName: "repo",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.config.IsGitHubEnterprise())
		})
	}
}
