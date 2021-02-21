/*
Copyright © 2020 Alex Watt <alex@alexcwatt.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"hymnal/data"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all hymns",
	Run: func(cmd *cobra.Command, args []string) {
		for _, hymn := range data.Hymns {
			fmt.Printf("%s\n", hymn)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
