/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/huh"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"gopkg.in/yaml.v2"
)

type TemplateBracket struct {
	begin        string `mapstructure:"begin"`
	end          string `mapstructure:"end"`
	regexPattern string `mapstructure:"regex-pattern"`
}

var (
	renderOption = "missingkey=error"
	brackets     = map[string]TemplateBracket{
		"curly":  {"{{", "}}", `\{\{(.*?)\}\}`},
		"square": {"[[", "]]", `\[\[(.*?)\]\]`},
	}
	bracketFormat = "curly"
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

func ReadProfileFile(filename string, target interface{}) error {
	// Read and unmarshal the YAML file into the provided target
	if err := ReadYAML(filename, target); err != nil {
		fmt.Printf("ERROR READING YAML FILE: %v\n", err)
		return err
	}

	return nil
}

// func ReadGitProfile(filename string) (HomerunDemo HomerunDemo) {

// 	if err := ReadYAML(filename, &HomerunDemo); err != nil {
// 		fmt.Printf("ERROR READING YAML FILE: %v\n", err)
// 	}

// 	return
// }

// READYAML READS AND PARSES A YAML FILE INTO A PROVIDED STRUCT
func ReadYAML(filename string, out interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("could not decode YAML: %w", err)
	}

	return nil
}

// HomerunDemo REPRESENTS THE STRUCTURE OF THE YAML FILE
type HomerunDemo struct {
	Surveys          []string                     `yaml:"surveys"`
	Templates        []string                     `yaml:"templates"`
	Values           []string                     `yaml:"values"`
	Aliases          []string                     `yaml:"aliases"`
	BodyTemplate     string                       `yaml:"bodyTemplate"`
	Authors          []string                     `yaml:"authors"`
	Usecases         map[string][]string          `yaml:"usecases"`
	MessageTemplates map[string]map[string]string `yaml:"messageTemplates"`
	Objects          map[string][]string          `yaml:"objects"`
	Artifacts        map[string][]string          `yaml:"artifacts"`
	Urls             map[string][]string          `yaml:"urls"`
}

// BUILD THE SURVEY FUNCTION WITH THE NEW RANDOM SETUP
func BuildSurvey(
	questions []*Question) (
	*huh.Form, map[string]interface{},
	error) {
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
		surveyQuestions, _, err := BuildSurvey(preSurvey)

		if err != nil {
			log.Fatalf("ERROR BUILDING SURVEY: %v", err)
		}
		log.Info("SURVEY FOUND")

		// RUN THE INTERACTIVE SURVEY
		err = surveyQuestions.Run()
		if err != nil {
			log.Fatalf("ERROR RUNNING SURVEY: %v", err)
		}

		// SET ANWERS TO ALL VALUES
		for _, question := range preSurvey {
			surveyValues[question.Name] = question.Default
		}

	} else {
		log.Info("NO SURVEY FOUND")
	}

	return surveyValues
}

func RenderAliases(aliases []string, allValues map[string]interface{}) map[string]interface{} {

	fmt.Println("ALL VALUES: ", allValues)

	for _, alias := range aliases {

		// SPLIT ALIAS KEY/VALUE BY :
		aliasValues := strings.Split(alias, ":")

		// RENDER KEY
		aliasKey, err := sthingsBase.RenderTemplateInline(aliasValues[0], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
		if err != nil {
			fmt.Println(err)
		}

		// RENDER VALUE
		aliasValue, err := sthingsBase.RenderTemplateInline(aliasValues[1], renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, allValues)
		if err != nil {
			fmt.Println(err)
		}

		// ASSIGN ALIAS TO ALL VALUES
		key := string(strings.TrimSpace(string(aliasKey)))
		value := string(strings.TrimSpace(string(aliasValue)))

		allValues[string(key)] = string(value)
		log.Info("ALIAS ADDED: ", key, ":", string(value))
	}

	return allValues
}

func RunSurveyFiles(surveys []string, values map[string]interface{}) map[string]interface{} {

	var allQuestions []*Question

	// LOAD ALL QUESTION FILES
	for _, questionFile := range surveys {

		// RENDER QUESTION FILE
		renderedQuestionFilePath, err := sthingsBase.RenderTemplateInline(questionFile, renderOption, brackets[bracketFormat].begin, brackets[bracketFormat].end, values)
		if err != nil {
			log.Error("ERROR RENDERING QUESTION FILE: ", err)
		}

		log.Info("LOADING QUESTION FILE: ", string(renderedQuestionFilePath))

		questions, _ := LoadQuestionFile(string(renderedQuestionFilePath), "")

		if len(questions) > 0 {
			log.Info("LOADED QUESTIONS FROM FILE: ", len(questions))
		} else {
			log.Warn("NO QUESTIONS FOUND IN FILE: ", string(renderedQuestionFilePath))
		}

		allQuestions = append(allQuestions, questions...)
	}

	survey, defaults, err := BuildSurvey(allQuestions)
	if err != nil {
		log.Fatalf("ERROR BUILDING SURVEY: %v", err)
	}

	err = survey.Run()
	if err != nil {
		log.Fatalf("ERROR RUNNING SURVEY: %v", err)
	}

	log.Info("DEFAULTS: ", defaults)

	return defaults
}

// ReadYAMLToMap reads a YAML file and exports it as a map[string]interface{}
func ReadYAMLToMap(filename string) (map[string]interface{}, error) {
	// Read the YAML file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Parse the YAML into a generic map
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	return result, nil
}

// RandomFromSlice selects a random element from a slice of strings
func RandomFromSlice(inputSlice []string) string {
	// Check if the slice is empty
	if len(inputSlice) == 0 {
		return "" // Return empty string if the slice is empty
	}
	// Generate a random index and return the element
	return inputSlice[rand.Intn(len(inputSlice))]
}

func RenderTemplateInlineWithFunctions(templateFunctions template.FuncMap, templateData string, templateVariables map[string]interface{}) ([]byte, error) {

	tmpl, err := template.New("template").Funcs(templateFunctions).Parse(templateData)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateVariables)
	if err != nil {
		panic(err)
	}

	return buf.Bytes(), nil

}
