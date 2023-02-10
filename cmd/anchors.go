package cmd

import (
	"encoding/json"
	"log"
	"os"
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
type stringMap map[string]interface{}

func convertMapToStringMap(inputMap map[interface{}]interface{}) stringMap {
	outputMap := make(stringMap)
	for k, v := range inputMap {
		stringKey, ok := k.(string)
		if !ok {
			continue
		}

		switch v := v.(type) {
		case map[interface{}]interface{}:
			outputMap[stringKey] = convertMapToStringMap(v)
		default:
			outputMap[stringKey] = v
		}
	}
	return outputMap
}

func convertYmlToJson(inputYmlName string) (string, error) {
	input, err := os.ReadFile(inputYmlName)
	if err != nil {
		return "", err
	}

	var data map[interface{}]interface{}
	err = yaml.Unmarshal(input, &data)
	if err != nil {
		return "", err
	}

	stringData := convertMapToStringMap(data)

	output, err := json.MarshalIndent(stringData, "", "  ")
	if err != nil {
		return "", err
	}

	outputJsonName := strings.Replace(inputYmlName, ".yml", ".json", 1)
	err = os.WriteFile(outputJsonName, output, 0644)
	if err != nil {
		return "", err
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
