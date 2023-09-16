/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push artifacts",
	Long:  `push artifacts to external systems`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("kind", "git", "push to - git|s3")
	pushCmd.Flags().String("source", "", "source file")
	pushCmd.Flags().String("auth", "", "auth data for external system")
	pushCmd.Flags().String("destination", "", "destination path")
}
