/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push artifacts",
	Long:  `push artifacts to external systems`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		to, _ := cmd.LocalFlags().GetString("to")
		sourceFile, _ := cmd.LocalFlags().GetString("source")

		sourceExists, _ := sthingsBase.VerifyIfFileOrDirExists(sourceFile, "file")
		if sourceExists {
			log.Info("SOURCE FOUND : ", templatePath)
			templateFile = sthingsBase.ReadFileToVariable(sourceFile)
		} else {
			log.Error("SOURCE NOT FOUND : ", sourceFile)
			os.Exit(3)
		}

		if to == "git" {
			fmt.Println("push called")

		}

	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("source", "", "source file path")
	pushCmd.Flags().String("auth", "", "auth data for external system")
	pushCmd.Flags().String("to", "git", "push to")
}
