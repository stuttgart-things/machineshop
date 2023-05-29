/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stuttgart-things/machineShop/surveys"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "manage infrastructure in any cloud",
	Long:  `predictably provision and manage infrastructure in any cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		gitPath, _ := cmd.LocalFlags().GetString("path")

		surveys.RunTerraform(gitRepository, gitPath)
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)

	terraformCmd.Flags().String("path", "machineShop/terraform", "path to terraform automation code")

}
