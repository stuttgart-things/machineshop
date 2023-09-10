/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get secret values",
	Long:  `get secret values from vault`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		authMethod, _ := cmd.LocalFlags().GetString("auth")
		secretPath, _ := cmd.LocalFlags().GetString("path")
		outputFormat, _ := cmd.LocalFlags().GetString("output")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")

		// START LOGGING
		log.Info("AUTH-METHOD: ", authMethod)
		log.Info("SECRET-PATH: ", secretPath)

		// CHECK FOR VAULT ENV VARS
		vaultAuthType, vaultAuthFound := sthingsCli.VerifyVaultEnvVars()
		log.Info("⚡️ VAULT CREDENDITALS ⚡️", vaultAuthType)

		if vaultAuthType == "approle" && vaultAuthFound {
			client, err := sthingsCli.CreateVaultClient()

			if err != nil {
				log.Error(err, "token creation (by approle) not working")
			}

			token, err := client.GetVaultTokenFromAppRole()

			if err != nil {
				log.Error(err, "token creation (by approle) not working")
			}

			os.Setenv("VAULT_TOKEN", token)
		}

		// GET SECRET VALUE
		secretValue := sthingsCli.GetVaultSecretValue(secretPath, os.Getenv("VAULT_TOKEN"))
		internal.HandleRenderOutput(outputFormat, destinationPath, secretValue, true)
		// fmt.Println(secretValue)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("auth", "approle", "vault auth method")
	getCmd.Flags().String("path", "", "path to vault secret")
	getCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	getCmd.Flags().String("destination", "", "path to output (if output file)")
}
