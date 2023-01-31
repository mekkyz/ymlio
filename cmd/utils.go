package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// This function as the name says converts the json file to YAML
func convertJsonToYAML(inputJsonFileName string) (string, error) {
	// get the content of the jsonData
	jsonData, err := ioutil.ReadFile(inputJsonFileName)
	outPutYamlFileName := strings.Replace(inputJsonFileName, ".json", ".yml", 1)
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
	err = ioutil.WriteFile(outPutYamlFileName, yamlData, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	// return the file name
	return outPutYamlFileName, nil
}

// This function as the name says converts the YAML file to JSON which will handle the anchors
func convertYamlFileToJson(inputYamlFileName string) (string, error) {
	// making use of yq processor command line program to do the conversion
	cmd := exec.Command("yq", "-r", "-o=json", inputYamlFileName)
	// get the output of the command
	output, err := cmd.CombinedOutput()
	// get the filename to save it with
	outputJsonFileName := strings.Replace(inputYamlFileName, ".yml", ".json", 1)
	if err != nil {
		fmt.Printf("Error running command: %s\n", err)
		return "", errors.New("Error running the command " + err.Error())
	} else {
		// Save the file as .json
		ioutil.WriteFile(outputJsonFileName, output, 0644)
		// return the filename with .json file extension
		return outputJsonFileName, nil
	}
}

// This function handle the anchors in YAML file
func HandleAnchor(inputYmlFileName string) (string, error) {
	jsonFileName, err := convertYamlFileToJson(inputYmlFileName)
	if err != nil {
		log.Fatalln(err)
		return "", errors.New("Unable to convert YAML file to json " + err.Error())
	}

	yamlFileName, err := convertJsonToYAML(jsonFileName)
	if err != nil {
		log.Fatalln(err)
		return "", errors.New("Unable to convert json file to YAML " + err.Error())
	}
	err = os.Remove(jsonFileName)
	if err != nil {
		fmt.Println("Unable to remove ", jsonFileName)
	}
	return yamlFileName, nil
}
