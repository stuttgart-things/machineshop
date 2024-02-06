/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"
	"strings"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	log = sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)
)

const regexPatternVaultSecretPath = `.+/data/.+:.+`

func VerifyReadKeyValues(templateValues []string, log *sthingsBase.Logger, enableVault bool) map[string]interface{} {

	vaultAuthType, vaultAuthFound := sthingsCli.VerifyVaultEnvVars()
	log.Info("⚡️ VAULT CREDENDITALS ⚡️", vaultAuthType)

	templateData := make(map[string]interface{})

	if len(templateValues) > 0 {

		if sthingsBase.VerifyIfStringIsBase64(templateValues[0]) {
			base64decodedValues := sthingsBase.DecodeBase64String(templateValues[0])
			templateValues = strings.Split(base64decodedValues, ",")
		}

		for _, v := range templateValues {
			values := strings.Split(v, "=")

			// IF VALUE IS BASE64 ENDCODED THIS IS NEEDED TO PATCH THE STRING BACK IF = ARE INCLUDED
			values = []string{values[0], strings.Join(values[1:], "=")}

			// CHECK FOR VAULT KV
			if len(sthingsBase.GetAllRegexMatches(values[1], regexPatternVaultSecretPath)) > 0 && enableVault {
				log.Info("VAULT PATH FOUND")
				VerifyVaultAuthType(vaultAuthType, log, vaultAuthFound)
				secretValue := sthingsCli.GetVaultSecretValue(values[1], os.Getenv("VAULT_TOKEN"))

				if secretValue == "" {
					log.Error("VAULT SECRET DOES NOT EXIST ", values[1])
				}

				values[1] = secretValue
			}

			// CHECK FOR EMPTY KEY
			if strings.TrimSpace(values[0]) == "" {
				fmt.Println("no key for value", strings.TrimSpace(values[1]), "defined. exiting")
				log.Error("no key defined. exiting")
				os.Exit(3)
			}

			templateData[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])

		}
	} else {
		log.Warn("NO VALUES DEFINED..")
	}

	return templateData
}

func VerifyVaultAuthType(vaultAuthType string, log *sthingsBase.Logger, vaultAuthFound bool) {

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
}

func ValidateGetVaultSecretValue(secretPathString string, log *sthingsBase.Logger) (secretValue string) {

	if len(sthingsBase.GetAllRegexMatches(secretPathString, regexPatternVaultSecretPath)) > 0 {
		vaultAuthType, vaultAuthFound := sthingsCli.VerifyVaultEnvVars()
		log.Info("⚡️ VAULT CREDENDITALS ⚡️", vaultAuthType)

		VerifyVaultAuthType(vaultAuthType, log, vaultAuthFound)
		secretValue = sthingsCli.GetVaultSecretValue(secretPathString, os.Getenv("VAULT_TOKEN"))
		log.Info("SECRET FOUND IN ", secretPathString)

	} else {
		secretValue = secretPathString
		log.Info("NO SECRET PATH FOUND - USING LITERAL ", secretValue)
	}

	return
}

func ValidateSourceFile(sourceFile string) {

	if sourceFile != "" {
		sourceExists, _ := sthingsBase.VerifyIfFileOrDirExists(sourceFile, "file")
		if sourceExists {
			log.Info("SOURCE FOUND : ", sourceFile)
		} else {
			log.Error("SOURCE NOT FOUND : ", sourceFile)
			os.Exit(3)
		}
	} else {
		log.Error("SOURCE UNDEFINED")
		os.Exit(3)
	}

}
