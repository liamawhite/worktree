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

# Shell wrapper for wt CLI that handles directory changes
# Usage: source this file or add the wt function to your shell profile

wt() {
    local cmd="$1"
    
    # Commands that should change directory
    if [[ "$cmd" == "add" || "$cmd" == "switch" || "$cmd" == "sw" ]]; then
        # Run the CLI and capture both output and potential directory change
        local output
        output=$(command wt "$@" 2>&1)
        local exit_code=$?
        
        # Print the output
        echo "$output"
        
        # Look for directory change indicator in the output
        local target_dir
        target_dir=$(echo "$output" | grep "^WT_CHDIR:" | sed 's/^WT_CHDIR://')
        
        # Change directory if target was specified
        if [[ -n "$target_dir" && -d "$target_dir" ]]; then
            cd "$target_dir" || echo "Warning: Failed to change to directory: $target_dir"
        fi
        
        return $exit_code
    else
        # For all other commands, just run normally
        command wt "$@"
    fi
}

# Completion support (if the binary supports it)
if command -v wt >/dev/null 2>&1; then
    complete -F _wt wt 2>/dev/null || true
fi
