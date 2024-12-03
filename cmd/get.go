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
			secretKeyFile := secretParameters[0]
			secretKey := secretParameters[1]
			log.Info("SECRET KEY: ", secretKey)

			// CHECK IF GIVEN SECRET FILE EXISTS
			secretFileExists, _ := sthingsBase.VerifyIfFileOrDirExists(secretKeyFile, "file")
			if secretFileExists {
				log.Info("SECRET FILE DOES NOT EXIST: ", secretKeyFile)
			} else {
				log.Error("SECRET FILE NOT FOUND: ", secretKeyFile)
				os.Exit(0)
			}

			// CHECK FOR SOPS ENV VARS
			sopsKeyFilePresent := sthingsCli.VerifyEnvVars([]string{"SOPS_AGE_KEY_FILE"})
			sopsKeyPresent := sthingsCli.VerifyEnvVars([]string{"SOPS_AGE_KEY"})
			log.Info("SOPS_AGE_KEY_FILE set on env: ", sopsKeyFilePresent)
			log.Info("SOPS_AGE_KEY      set on env: ", sopsKeyPresent)

			if sopsKeyPresent || sopsKeyFilePresent {

				// DECRYPT SOPS FILE
				err, plain := sthingsCli.DecryptSopsFile(secretKeyFile, "yaml")
				if err != nil {
					fmt.Println("FAILED TO DECRYPT: %w", err)
				} else {

					// READ IN DECRYPTED FILE
					defaultVariables = sthingsCli.ReadYamlKeyValuesFromFile([]byte(plain))
					// PRINT DECRYPTED KEY
					fmt.Println(defaultVariables[secretKey])
				}

			} else {
				log.Error("NO SOPS-KEY FOUND IN ENV")
				os.Exit(0)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("auth", "approle", "vault auth method")
	getCmd.Flags().String("path", "", "path to vault secret")
	getCmd.Flags().String("system", "vault", "secret system: vault|sops")
	getCmd.Flags().String("output", "stdout", "output stdout|file")
	getCmd.Flags().String("destination", "", "path to output (if output file)")
	getCmd.Flags().Bool("b64", false, "decode base64 for output")
}
