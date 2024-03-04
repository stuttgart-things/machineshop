/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/stuttgart-things/machineshop/internal"
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
		b64DecodeOption, _ := cmd.LocalFlags().GetBool("b64")

		// START LOGGING
		log.Info("AUTH-METHOD: ", authMethod)
		log.Info("SECRET-PATH: ", secretPath)

		// CHECK FOR VAULT ENV VARS
		vaultAuthType, vaultAuthFound := sthingsCli.VerifyVaultEnvVars()
		log.Info("⚡️ VAULT CREDENDITALS ⚡️", vaultAuthType)

		internal.VerifyVaultAuthType(vaultAuthType, log, vaultAuthFound)

		// GET SECRET VALUE
		secretValue := sthingsCli.GetVaultSecretValue(secretPath, os.Getenv("VAULT_TOKEN"))
		internal.HandleRenderOutput(outputFormat, destinationPath, secretValue, b64DecodeOption, true)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("auth", "approle", "vault auth method")
	getCmd.Flags().String("path", "", "path to vault secret")
	getCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	getCmd.Flags().String("destination", "", "path to output (if output file)")
	getCmd.Flags().Bool("b64", false, "decode base64 for output")
}
