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