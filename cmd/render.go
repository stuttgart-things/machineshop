/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/getsops/sops/v3/decrypt"
	kaeffken "github.com/stuttgart-things/kaeffken/modules"
	"github.com/stuttgart-things/machineshop/internal"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"gopkg.in/yaml.v2"

	billy "github.com/go-git/go-billy/v5"
	"github.com/spf13/cobra"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

// DEFINE THE STRUCTURE TO MATCH THE YAML FILE
type InfraConfig struct {
	Apps map[string]string `yaml:"template"`
}

type TemplateBracket struct {
	begin        string `mapstructure:"begin"`
	end          string `mapstructure:"end"`
	regexPattern string `mapstructure:"regex-pattern"`
}

var (
	templateFile     string
	defaultsFile     string
	fileFormat       = "yaml"
	repo             billy.Filesystem
	defaultVariables = make(map[string]interface{})
	sopsSecrets      = make(map[string]interface{})
	flagVariables    = make(map[string]interface{})
	brackets         = map[string]TemplateBracket{
		"curly":  TemplateBracket{"{{", "}}", `\{\{(.*?)\}\}`},
		"square": TemplateBracket{"[[", "]]", `\[\[(.*?)\]\]`},
	}
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "render templates",
	Long:  `Render templates to files`,

	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		source, _ := cmd.LocalFlags().GetString("source")
		templateKind, _ := cmd.LocalFlags().GetString("kind")
		gitPath, _ := cmd.LocalFlags().GetString("path")
		templatePath, _ := cmd.LocalFlags().GetString("template")
		defaultsPath, _ := cmd.LocalFlags().GetString("defaults")
		ageKey, _ := cmd.LocalFlags().GetString("age")
		secretsPath, _ := cmd.LocalFlags().GetString("secrets")
		outputFormat, _ := cmd.LocalFlags().GetString("output")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")
		templateValues, _ := cmd.Flags().GetStringSlice("values")
		forceRenderOption, _ := cmd.LocalFlags().GetBool("force")
		b64DecodeOption, _ := cmd.LocalFlags().GetBool("b64")
		bracketFormat, _ := cmd.LocalFlags().GetString("brackets")
		templateKey, _ := cmd.LocalFlags().GetString("key")

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/RENDER")

		// Verify Bracket Format
		_, bracketFormatExists := brackets[bracketFormat]
		if !bracketFormatExists {
			log.Warn("GIVEN BRACKET FORMAT DOES NOT EXIST - GOING TO DEFAULT ", bracketFormat)
			bracketFormat = "curly"
		}

		log.Info("BRACKET FORMAT: ", bracketFormat)
		log.Info("TEMPLATE-PATH: ", templatePath)
		log.Info("DEFAULTS: ", defaultsPath)
		log.Info("OUTPUT-FORMAT: ", outputFormat)
		log.Info("DESTINATION-PATH: ", destinationPath+"\n")

		// HANDLE SOURCE: GIT FOR TEMPLATE + DEFAULTS
		if source == "git" {
			repo, _ = sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
			templateFile = sthingsCli.ReadFileContentFromGitRepo(repo, templatePath)

			// READ DEFAULTS (IF DEFINED)
			if defaultsPath != "" {
				defaultsFile = sthingsCli.ReadFileContentFromGitRepo(repo, defaultsPath)

				log.Info("LOADED DEFAULTS FILE FROM: ", defaultsPath)
				fmt.Println(defaultsFile)
				defaultVariables = sthingsCli.ReadYamlKeyValuesFromFile([]byte(defaultsFile))

			} else {
				log.Info("NO DEFAULTS FILE FROM GIT DEFINED")
			}

			// HANDLE SOURCE:LOCAL FOR TEMPLATE + DEFAULTS
		} else if source == "local" {
			templateExists, _ := sthingsBase.VerifyIfFileOrDirExists(templatePath, "file")

			if templateExists {
				log.Info("LOCAL TEMPLATE FOUND : ", templatePath)
				templateFile = sthingsBase.ReadFileToVariable(templatePath)
			} else {
				log.Error("LOCAL TEMPLATE NOT FOUND : ", templatePath)
				os.Exit(3)
			}

			// READ DEFAULTS (IF DEFINED)
			if defaultsPath != "" {

				defaultsFileExists, _ := sthingsBase.VerifyIfFileOrDirExists(defaultsPath, "file")

				// IF DEFAULTS FILE EXISTS
				if defaultsFileExists {
					defaultsFile = sthingsBase.ReadFileToVariable(defaultsPath)
					log.Info("LOADED DEFAULTS FILE FROM: ", defaultsPath)
					fmt.Println(defaultsFile)
					defaultVariables = sthingsCli.ReadYamlKeyValuesFromFile([]byte(defaultsFile))
				} else {
					log.Error("LOCAL DEFAULTS FILE NOT FOUND : ", defaultsPath)
					os.Exit(3)
				}
			} else {
				log.Info("NO DEFAULTS FILE DEFINED")
			}

			if secretsPath != "" {

				log.Info("SECRETS FILE DEFINED: ", secretsPath)

				// CHECK IF AGE KEY IS SET
				if ageKey != "" {
					os.Setenv("SOPS_AGE_KEY", ageKey)
					log.Info("USING AGE KEY: ", ageKey)
				}

				if ageKey == "" && os.Getenv("SOPS_AGE_KEY") == "" {
					log.Warn("SOPS_AGE_KEY NOT SET")
					log.Error("AGE KEY NOT SET")
				}

				// CHECK IF GIVEN SECRET FILE EXISTS
				secretFileExists, _ := sthingsBase.VerifyIfFileOrDirExists(secretsPath, "file")
				if secretFileExists {
					log.Info("SECRET FILE DOES EXIST: ", secretsPath)
				} else {
					log.Error("SECRET FILE NOT FOUND: ", secretsPath)
					os.Exit(0)
				}

				decryptedFile, err := decrypt.File(secretsPath, fileFormat)
				if err != nil {
					log.Error("FAILED TO DECRYPT: ", err)
				}

				allDecryptedSecrets := kaeffken.CreateSecretsMap(decryptedFile, nil)
				log.Info("ALL DECRYPTED SECRETS: ", allDecryptedSecrets)

				defaultVariables = sthingsBase.MergeMaps(defaultVariables, allDecryptedSecrets)

			} else {
				log.Info("NO SOPS ENCRYPTED SECRETS FILE DEFINED")
			}

		} else {
			log.Error("SOURCE CAN BE ONLY: GIT OR LOCAL", source)
			os.Exit(3)
		}

		// READ VALUES (IF DEFINED)
		if len(templateValues) > 0 {
			flagVariables = internal.VerifyReadKeyValues(templateValues, log, enableVault)
			fmt.Println("VALUES", flagVariables)
		} else {
			log.Warn("NO VALUES DEFINED")
		}

		// MERGE DEFAULTS + VALUES
		variables := sthingsBase.MergeMaps(defaultVariables, flagVariables)

		// CHECK FOR YAML LIST CONTENT
		for key, values := range variables {

			listContent, isList := sthingsBase.GetRegexSubMatch(values.(string), `\[(.*?)\]`)
			if isList {
				log.Info("FOUND MAP OR A LIST: ", key)
				fmt.Println(values)

				_, isMap := sthingsBase.GetRegexSubMatch(values.(string), `map\[(.*?)\]`)
				if isMap {
					log.Info("FOUND MAP")
					fmt.Println(key)

					// REED YAML FILE
					yamlFile, err := os.ReadFile(defaultsPath)
					if err != nil {
						fmt.Printf("yamlFile.Get err #%v ", err)
					}
					// GET TRANSFORMED VALUES FROM DICT
					mapValues := sthingsCli.GetYAMLMapValues(yamlFile, key)

					// MERGE W/ VARS
					variables = sthingsBase.MergeMaps(variables, mapValues)
					log.Info("MAP VALUES: ", variables)

				} else {
					log.Info("FOUND LIST")

					// GET LIST AS LIST
					variables[key] = strings.Split(listContent, " ")

					// GET RANDOM (SINGLE) VALUE FROM LIST
					variables["RANDOM"+key] = sthingsBase.GetRandomPickFromSlice(strings.Split(listContent, " "))
				}

			}
		}

		// GET MULTIKEY TEMPLATE
		renderOption := "missingkey=error"
		if forceRenderOption {
			renderOption = "missingkey=zero"
		}

		if templateKind == "multikey" {

			data, err := os.ReadFile(templatePath)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			var config InfraConfig
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			for key, value := range config.Apps {

				if templateKey == key {
					templateFile = value
					continue
				}
			}

		}

		renderedTemplate, err := sthingsBase.RenderTemplateInline(templateFile, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, variables)
		if err != nil {
			fmt.Println(err)
		}

		// HANDLE OUTPUT
		if len(renderedTemplate) == 0 {
			log.Error("RENDERED TEMPLATE IS EMPTY")
			os.Exit(3)
		} else {
			internal.HandleRenderOutput(outputFormat, destinationPath, string(renderedTemplate), b64DecodeOption, enableVault)
		}
	},
}

func GetYAMLMapValues(content []byte, dictName string) (transformedValues map[string]interface{}) {

	obj := make(map[string]interface{})
	transformedValues = make(map[string]interface{})

	err := yaml.Unmarshal(content, obj)
	if err != nil {
		fmt.Printf("UNMARSHALING: %v", err)
	}

	// ITERATE OVER KEYS AND VALUES
	for dictKey, subKey := range obj[dictName].(map[interface{}]interface{}) {

		for subKey, value := range subKey.(map[interface{}]interface{}) {
			transformedValues[dictName+"_"+fmt.Sprint(dictKey)+"_"+fmt.Sprint(subKey)] = fmt.Sprint(value)
		}
	}

	return transformedValues
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().String("source", "git", "source of profile: git or local")
	renderCmd.Flags().String("kind", "file", "kind of template: file or multikey")
	renderCmd.Flags().String("template", "tests/template.yaml", "path to to be rendered template")
	renderCmd.Flags().String("defaults", "", "path to defaults template file")
	renderCmd.Flags().String("brackets", "curly", "template bracket format - curly|square")
	renderCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	renderCmd.Flags().String("secrets", "", "sops encryted secrets file")
	renderCmd.Flags().String("age", "", "sops age key")
	renderCmd.Flags().String("destination", "", "path to output (if output file)")
	renderCmd.Flags().Bool("force", false, "force rendering by missing keys")
	renderCmd.Flags().StringSlice("values", []string{}, "templating values")
	renderCmd.Flags().String("key", "", "to be rendered keys - only possible if kind is multikey")
	renderCmd.Flags().Bool("b64", false, "decode base64 for output")
}
