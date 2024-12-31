/*
Copyright © 202 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	kaeffken "github.com/stuttgart-things/kaeffken/modules"

	"github.com/getsops/sops/v3/decrypt"
	ipservice "github.com/stuttgart-things/clusterbook/ipservice"
	"github.com/stuttgart-things/machineshop/internal"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
	"google.golang.org/grpc"

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
		path, _ := cmd.LocalFlags().GetString("path")
		output, _ := cmd.LocalFlags().GetString("output")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")
		b64DecodeOption, _ := cmd.LocalFlags().GetBool("b64")
		system, _ := cmd.LocalFlags().GetString("system")
		ageKey, _ := cmd.LocalFlags().GetString("key")
		fileFormat, _ := cmd.LocalFlags().GetString("format")

		// START LOGGING
		log.Info("PATH: ", path)
		log.Info("SYSTEM: ", system)

		switch system {

		case "ips":

			log.Info("⚡️ CONNECTING TO CLUSTERBOOK ⚡️")
			log.Info("CLUSTERBOOK SERVER: ", destinationPath)
			log.Info("COUNT IPS: ", sthingsBase.ConvertStringToInteger(output))
			log.Info("NETWORK-KEY: ", path)

			clusterBookServer := destinationPath //"clusterbook.rke2.sthings-vsphere.labul.sva.de:443"
			secureConnection := "true"
			countIps := int32(sthingsBase.ConvertStringToInteger(output))
			networkKey := path

			// SELECT CREDENTIALS BASED ON SECURECONNECTION
			conn, err := grpc.NewClient(clusterBookServer, internal.GetCredentials(secureConnection))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()

			c := ipservice.NewIpServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// Testen der GetIpAddressRange-Methode
			ipReq := &ipservice.IpRequest{
				CountIpAddresses: countIps,
				NetworkKey:       networkKey,
			}

			ipRes, err := c.GetIpAddressRange(ctx, ipReq)
			if err != nil {
				log.Fatalf("could not get IP address range: %v", err)
			}

			log.Printf("Available IPs: %s", ipRes.IpAddressRange)

		case "vault":

			log.Info("AUTH-METHOD: ", authMethod)

			// CHECK FOR VAULT ENV VARS
			vaultAuthType, vaultAuthFound := sthingsCli.VerifyVaultEnvVars()
			log.Info("⚡️ VAULT CREDENDITALS ⚡️", vaultAuthType)

			internal.VerifyVaultAuthType(vaultAuthType, log, vaultAuthFound)

			// GET SECRET VALUE
			secretValue := sthingsCli.GetVaultSecretValue(path, os.Getenv("VAULT_TOKEN"))
			internal.HandleRenderOutput(output, destinationPath, secretValue, b64DecodeOption, true)

		case "sops":
			fmt.Println("SOPS")

			// GET SECRET PARAMETERS
			secretParameters := strings.Split(path, ":")
			decryptedFilePath := secretParameters[0]
			secretKey := secretParameters[1]
			log.Info("SECRET KEY: ", secretKey)

			// CHECK IF GIVEN SECRET FILE EXISTS
			secretFileExists, _ := sthingsBase.VerifyIfFileOrDirExists(decryptedFilePath, "file")
			if secretFileExists {
				log.Info("SECRET FILE DOES NOT EXIST: ", decryptedFilePath)
			} else {
				log.Error("SECRET FILE NOT FOUND: ", decryptedFilePath)
				os.Exit(0)
			}

			// CHECK FOR SOPS ENV VARS
			// CHECK IF AGE KEY IS SET
			if ageKey != "" {
				os.Setenv("SOPS_AGE_KEY", ageKey)
				log.Info("USING AGE KEY: ", ageKey)
			}

			if ageKey == "" && os.Getenv("SOPS_AGE_KEY") == "" {
				log.Warn("SOPS_AGE_KEY NOT SET")
				log.Error("AGE KEY NOT SET")
			}

			decryptedFile, err := decrypt.File(decryptedFilePath, fileFormat)
			if err != nil {
				log.Error("FAILED TO DECRYPT: ", err)
			}

			fmt.Println("DECRYPTED FILE: ", string(decryptedFile))

			allDecryptedSecrets := kaeffken.CreateSecretsMap(decryptedFile, nil)

			// LOOP OVER ALL DECRYPTED SECRETS
			var secretValueForGivenKey string
			for key, value := range allDecryptedSecrets {

				if key == secretKey {
					secretValueForGivenKey = value.(string)
				}

				log.Info("FOUND KEY: ", key)
			}

			if secretValueForGivenKey == "" {
				log.Error("SECRET KEY NOT FOUND: ", secretKey)
				os.Exit(0)
			} else {
				log.Info("SECRET VALUE FOR KEY "+secretKey+": ", secretValueForGivenKey)
			}

			// if sopsKeyPresent || sopsKeyFilePresent {

			// 	// DECRYPT SOPS FILE
			// 	err, plain := sthingsCli.DecryptSopsFile(decryptedFile, "yaml")
			// 	if err != nil {
			// 		fmt.Println("FAILED TO DECRYPT: %w", err)
			// 	} else {

			// 		// READ IN DECRYPTED FILE
			// 		defaultVariables = sthingsCli.ReadYamlKeyValuesFromFile([]byte(plain))
			// 		// PRINT DECRYPTED KEY
			// 		fmt.Println(defaultVariables[secretKey])
			// 	}

			// } else {
			// 	log.Error("NO SOPS-KEY FOUND IN ENV")
			// 	os.Exit(0)
			// }
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("auth", "approle", "vault auth method")
	getCmd.Flags().String("path", "", "path to vault secret")
	getCmd.Flags().String("key", "", "sops age key")
	getCmd.Flags().String("system", "vault", "secret system: vault|sops")
	getCmd.Flags().String("output", "stdout", "output stdout|file")
	getCmd.Flags().String("destination", "", "path to output (if output file)")
	getCmd.Flags().Bool("b64", false, "decode base64 for output")
	getCmd.Flags().String("format", "yaml", "sops file format/extension")
}
