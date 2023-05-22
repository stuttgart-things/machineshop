/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/machineShop/surveys"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "manage infrastructure in any cloud",
	Long:  `predictably provision and manage infrastructure in any cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("terraform called")
		surveys.RunTerraform()
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)

}
