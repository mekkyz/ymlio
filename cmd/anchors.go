package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// Convert JSON to YAML
func convertJsonToYml(inputJsonName string) (string, error) {
	// get the content of the JSON
	jsonData, err := os.ReadFile(inputJsonName)
	outputYmlName := strings.Replace(inputJsonName, ".json", ".yml", 1)
	if err != nil {
		log.Fatalln(err)
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		log.Fatalln(err)
	}
	// Marshal the Go map into YAML
	yamlData, err := yaml.Marshal(jsonMap)
	if err != nil {
		log.Fatalln(err)
	}
	// Write the YAML data to a file
	err = os.WriteFile(outputYmlName, yamlData, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	return outputYmlName, nil
}

// Convert YAML to JSON
func convertYmlToJson(inputYmlName string) (string, error) {
	// Making use of yq processor cli program to do the conversion
	cmd := exec.Command("yq", "-r", "-o=json", inputYmlName)
	// Get the output of the command
	output, err := cmd.CombinedOutput()
	// Get the filename to save it with
	outputJsonName := strings.Replace(inputYmlName, ".yml", ".json", 1)
	if err != nil {
		log.Fatalln(err)
	} else {
		// Save the file as .json
		os.WriteFile(outputJsonName, output, 0644)
	}
	return outputJsonName, nil
}

// Handling anchors
func handleAnchor(inputYmlFileName string) (string, error) {
	jsonFileName, err := convertYmlToJson(inputYmlFileName)
	if err != nil {
		log.Fatalln(err)
	}
	yamlName, err := convertJsonToYml(jsonFileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.Remove(jsonFileName)
	if err != nil {
		log.Fatalln(err)
	}
	return yamlName, nil
}
