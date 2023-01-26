package utils

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

func convertJsonToYAML(inputJsonFileName string) (string, error) {
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
	return outPutYamlFileName, nil
}

func convertYamlFileToJson(inputYamlFileName string) (string, error) {
	// command := fmt.Sprintf("yq -r -o=json %s > %s", inputYamlFileName, outputJsonFileName)
	cmd := exec.Command("yq", "-r", "-o=json", inputYamlFileName)
	output, err := cmd.CombinedOutput()
	outputJsonFileName := strings.Replace(inputYamlFileName, ".yml", ".json", 1)
	if err != nil {
		fmt.Printf("Error running command: %s\n", err)
		return "", errors.New("Error running the command " + err.Error())
	} else {
		ioutil.WriteFile(outputJsonFileName, output, 0644)
		return outputJsonFileName, nil
	}
}

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
