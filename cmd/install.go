/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"github.com/stuttgart-things/machineShop/internal"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install binaries",
	Long:  `install binaries`,
	Run: func(cmd *cobra.Command, args []string) {
		gitPath, _ := cmd.LocalFlags().GetString("path")
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/INSTALL")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

}
