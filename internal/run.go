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
