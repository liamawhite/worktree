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

package cmd

import (
	"fmt"

	"github.com/liamawhite/worktree/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display the version, build commit, and build date of the worktree CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.GetInfo()

		if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
			jsonStr, err := info.JSON()
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(jsonStr)
		} else {
			fmt.Println(info.String())
		}
	},
}

func init() {
	versionCmd.Flags().BoolP("json", "j", false, "output version information in JSON format")
}
