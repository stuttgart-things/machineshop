/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"gopkg.in/yaml.v2"
)

// QUESTION STRUCT TO HOLD THE QUESTION DATA FROM YAML
type Question struct {
	Prompt          string                 `yaml:"prompt"`
	Name            string                 `yaml:"name"`
	Default         string                 `yaml:"default,omitempty"`
	DefaultFunction string                 `yaml:"default_function,omitempty"`
	DefaultParams   map[string]interface{} `yaml:"default_params,omitempty"`
	Options         []string               `yaml:"options"`
	Kind            string                 `yaml:"kind,omitempty"` // "function" instead of "text"
	MinLength       int                    `yaml:"minLength,omitempty"`
	MaxLength       int                    `yaml:"maxLength,omitempty"`
	Type            string                 `yaml:"type,omitempty"` // Updated field to match the YAML
}

func LoadQuestionFile(filename, yamlKey string) ([]*Question, error) {
	var questions []*Question

	// READ THE YAML FILE
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// ATTEMPT TO UNMARSHAL AS A LIST DIRECTLY (FOR YAML WITHOUT `yamlKey` KEY)
	if err := yaml.Unmarshal(data, &questions); err == nil {
		return questions, nil
	}

	// IF UNMARSHALING DIRECTLY FAILS, UNMARSHAL INTO A MAP AND EXTRACT BY `yamlKey`
	var genericMap map[string]interface{}
	if err := yaml.Unmarshal(data, &genericMap); err != nil {
		return nil, err
	}

	// EXTRACT THE DATA ASSOCIATED WITH `yamlKey`
	if rawQuestions, found := genericMap[yamlKey]; found {
		rawData, err := yaml.Marshal(rawQuestions) // Marshal back into YAML for unmarshaling into []*Question
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(rawData, &questions); err != nil {
			return nil, err
		}
		return questions, nil
	}

	// RETURN AN ERROR IF `yamlKey` IS NOT FOUND
	return nil, fmt.Errorf("key '%s' not found in YAML file", yamlKey)
}

func ReadGitProfile(filename string) (HomerunDemo HomerunDemo) {

	if err := ReadYAML(filename, &HomerunDemo); err != nil {
		fmt.Printf("ERROR READING YAML FILE: %v\n", err)
	}

	return
}

// READYAML READS AND PARSES A YAML FILE INTO A PROVIDED STRUCT
func ReadYAML(filename string, out interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("could not decode YAML: %v", err)
	}

	return nil
}

// HomerunDemo REPRESENTS THE STRUCTURE OF THE YAML FILE
type HomerunDemo struct {
	Surveys   []string `yaml:"surveys"`
	Templates []string `yaml:"templates"`
	Values    []string `yaml:"values"`
	Aliases   []string `yaml:"aliases"`
}

// BUILD THE SURVEY FUNCTION WITH THE NEW RANDOM SETUP
func BuildSurvey(questions []*Question) (*huh.Form, map[string]interface{}, error) {
	var groupFields []*huh.Group
	answers := make(map[string]interface{}) // To hold question names and resolved default values

	// Create a new random source
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // New random generator

	// Iterate over each question to create the survey fields
	for _, question := range questions {
		var field huh.Field

		// Set up default values for options if applicable
		if question.Default == "" && len(question.Options) > 0 {
			question.Default = question.Options[r.Intn(len(question.Options))] // Random default selection
		}

		// Handle the different question kinds
		switch question.Kind {
		case "function": // Handle "function" kind
			if question.DefaultFunction != "" {
				if fn, ok := defaultFunctions[question.DefaultFunction]; ok {
					question.Default = fn(question.DefaultParams)
				} else {
					return nil, nil, fmt.Errorf("default function %s not found", question.DefaultFunction)
				}
			}

			field = huh.NewInput().
				Title(question.Prompt).
				Value(&question.Default).
				Validate(func(input string) error {
					if len(input) < question.MinLength {
						return fmt.Errorf("input too short, minimum length is %d", question.MinLength)
					}
					if len(input) > question.MaxLength {
						return fmt.Errorf("input too long, maximum length is %d", question.MaxLength)
					}
					return nil
				})

		case "ask": // Handle "ask" kind
			field = huh.NewInput().
				Title(question.Prompt).
				Value(&question.Default).
				Validate(func(input string) error {
					if len(input) < question.MinLength {
						return fmt.Errorf("input too short, minimum length is %d", question.MinLength)
					}
					if len(input) > question.MaxLength {
						return fmt.Errorf("input too long, maximum length is %d", question.MaxLength)
					}
					return nil
				})

			// Store a placeholder for user input
			answers[question.Name] = "" // Will be updated during survey run

		default: // Handle multiple choice select options or other fields
			options := make([]huh.Option[string], len(question.Options))
			for i, opt := range question.Options {
				options[i] = huh.NewOption(opt, opt)
			}

			field = huh.NewSelect[string]().
				Title(question.Prompt).
				Options(options...).
				Value(&question.Default)
		}

		// Determine the data type and store the value correctly in the answers map
		switch question.Type {
		case "boolean": // Store as boolean
			answers[question.Name] = question.Default == "Yes" // Convert Yes/No to true/false

		case "int": // Store as integer
			if intValue, err := strconv.Atoi(question.Default); err == nil {
				answers[question.Name] = intValue
			} else {
				return nil, nil, fmt.Errorf("invalid default value for int type: %s", question.Default)
			}

		default: // Default to string
			answers[question.Name] = question.Default
		}

		// Create a group and add the field to it
		group := huh.NewGroup(field)
		groupFields = append(groupFields, group)
	}

	// Create and return the form along with the answers map
	return huh.NewForm(groupFields...), answers, nil
}

// FUNCTION MAPPING
var defaultFunctions = map[string]func(params map[string]interface{}) string{

	"getDefaultFavoriteFood": func(params map[string]interface{}) string {
		if spiceLevel, ok := params["spiceLevel"].(string); ok && spiceLevel != "" {
			return fmt.Sprintf("spicy %s", spiceLevel)
		}
		return "steak"
	},
	"getDefaultDrink": func(params map[string]interface{}) string {
		if temp, ok := params["temperature"].(string); ok && temp != "" {
			return fmt.Sprintf("%s water", temp)
		}
		return "water"
	},
}

func RunSurvey(profilePath, surveyKey string) (surveyValues map[string]interface{}) {

	surveyValues = make(map[string]interface{})

	// READ PROFILE AND SURVEY BY KEY
	preSurvey, _ := LoadQuestionFile(profilePath, surveyKey)

	// IF SURVEY EXISTS, RUN IT
	if len(preSurvey) > 0 {
		surveyQuestions, answers, err := BuildSurvey(preSurvey)
		if err != nil {
			log.Fatalf("ERROR BUILDING SURVEY: %v", err)
		}
		log.Info("SURVEY FOUND")

		// RUN THE INTERACTIVE SURVEY
		err = surveyQuestions.Run()
		if err != nil {
			log.Fatalf("ERROR RUNNING SURVEY: %v", err)
		}
		surveyValues = answers

	} else {
		log.Info("NO SURVEY FOUND")
	}

	return surveyValues
}
