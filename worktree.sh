#!/bin/bash

# Copyright 2025 Liam White
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Shell wrapper for worktree CLI that handles directory changes
# Usage: source this file or add the worktree function to your shell profile

worktree() {
    local cmd="$1"
    
    # Commands that should change directory
    if [[ "$cmd" == "add" || "$cmd" == "switch" || "$cmd" == "sw" ]]; then
        # Capture stderr to look for WT_CHDIR while preserving stdout and interactive TUI
        local temp_file
        temp_file=$(mktemp)
        
        # Run command with stderr redirected to temp file
        command worktree "$@" 2> >(tee "$temp_file" >&2)
        local exit_code=$?
        
        # Look for directory change indicator in stderr
        local target_dir
        target_dir=$(grep "^WT_CHDIR:" "$temp_file" 2>/dev/null | sed 's/^WT_CHDIR://')
        
        # Clean up temp file
        rm -f "$temp_file"
        
        # Change directory if target was specified
        if [[ -n "$target_dir" && -d "$target_dir" ]]; then
            cd "$target_dir" || echo "Warning: Failed to change to directory: $target_dir"
        fi
        
        return $exit_code
    else
        # For all other commands, just run normally
        command worktree "$@"
    fi
}

# Completion support (if the binary supports it)
if command -v worktree >/dev/null 2>&1; then
    complete -F _worktree worktree 2>/dev/null || true
fi
