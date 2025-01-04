/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"gopkg.in/yaml.v2"
)

type RunScripts struct {
	Globals map[string]interface{} `yaml:"globals"`
	Run     []RunStep              `yaml:"run"`
}

type RunStep struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Script      string                 `yaml:"script"`
	Vars        map[string]interface{} `yaml:"vars,omitempty"`
}

func LoadRunConfig(configString string) (config RunScripts) {

	// PARSE YAML FROM THE STRING INTO THE RUNSCRIPTS STRUCT
	if err := yaml.Unmarshal([]byte(configString), &config); err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	// RETURN THE CONFIG
	return config
}

// func CreateRunSurvey(config RunScripts) {

// 	// Create a slice to hold the names and descriptions for the survey
// 	var options []string
// 	for _, step := range config.Run {
// 		option := fmt.Sprintf("%s: %s", step.Name, step.Description)
// 		options = append(options, option)
// 	}

// 	// Create a new instance of huh for the survey
// 	// Create the multi-select prompt using huh.NewMultiSelect
// 	prompt := huh.NewMultiSelect("Select the steps to install", options)

// 	// Run the prompt and allow the user to select multiple options
// 	choices, err := prompt.Run()
// 	if err != nil {
// 		log.Fatalf("Error during survey: %v", err)
// 	}

// 	// Display the selected options (step names)
// 	fmt.Println("You selected:")
// 	for _, choice := range choices {
// 		// Output the original step based on the selection
// 		for _, step := range config.Run {
// 			if fmt.Sprintf("%s: %s", step.Name, step.Description) == choice {
// 				fmt.Printf("Name: %s\nDescription: %s\nScript: %s\n", step.Name, step.Description, step.Script)
// 				break
// 			}
// 		}
// 	}
// }
