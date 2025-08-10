package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindGitRoot(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (string, func())
		wantErr bool
	}{
		{
			name: "finds bare repo",
			setup: func(t *testing.T) (string, func()) {
				tmpDir := t.TempDir()
				bareDir := filepath.Join(tmpDir, ".bare")
				if err := os.MkdirAll(bareDir, 0755); err != nil {
					t.Fatal(err)
				}
				
				oldCwd, _ := os.Getwd()
				os.Chdir(tmpDir)
				
				return tmpDir, func() { os.Chdir(oldCwd) }
			},
			wantErr: false,
		},
		{
			name: "fails outside git repo",
			setup: func(t *testing.T) (string, func()) {
				tmpDir := t.TempDir()
				oldCwd, _ := os.Getwd()
				os.Chdir(tmpDir)
				
				return "", func() { os.Chdir(oldCwd) }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, cleanup := tt.setup(t)
			defer cleanup()

			got, err := FindGitRoot()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				// On macOS, temp dirs might be symlinked, so compare resolved paths
				expectedResolved, _ := filepath.EvalSymlinks(expected)
				gotResolved, _ := filepath.EvalSymlinks(got)
				assert.Equal(t, expectedResolved, gotResolved)
			}
		})
	}
}

func TestRunGitCommandOutput(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid command",
			args:    []string{"--version"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"--invalid-flag-that-does-not-exist"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunGitCommandOutput(tt.args...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}