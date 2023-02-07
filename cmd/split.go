package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

// hold the filename on the config keys, and also if it has import flag, get the filename that will be imported
type ImportData struct {
	FileName         string
	FileNameToImport string
}

var only bool

func createFileIfNotExist(filePath string) {

	currentPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	fullPath := filepath.Join(currentPath, filePath)

	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	}
}

// Save YAML config file, by acception the filename to be saved and also the values inside the config to be saved
func saveConfig(savedFileName string, yamlValues interface{}) {
	// Initialize an empty slice of byte to hold the converted values
	var converterBytes []byte
	var err error

	// Check if the type of the yamlValues input is a map
	if reflect.TypeOf(yamlValues).Kind() == reflect.Map {
		// Convert the yamlValues to a slice of bytes using yaml.Marshal
		converterBytes, err = yaml.Marshal(yamlValues)
		if err != nil {
			fmt.Println("Cannot convert converter to []byte: ", err)
			// Return if there is an error
			return
		}
	}
	// Write the converted values to the file with the given filename using ioutil.WriteFile
	err = os.WriteFile(savedFileName, []byte(converterBytes), 0644)
	if err != nil {
		fmt.Println(err)
		// Return if there is an error
		return
	}
}

// save a text file with a given filename and content
func saveText(fileName string, content string) {
	// Create a new file with the given filename using os.Create
	file, err := os.Create(fileName)
	// error handling
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Defer the closing of the file until all other function calls have completed
	defer file.Close()

	// Write the content string to the file using file.WriteString
	_, err = file.WriteString(content)
	if err != nil {
		// Print an error message if there is an error writing to the file
		fmt.Println(err)
		return
	}
}

// rename multiple raw files with new filenames
func renameRaw(rawFileNames map[string]string) {
	// Loop through the map of old filenames and new filenames
	for oldFileName, newFileName := range rawFileNames {
		// Rename the file with the old filename to the new filename using os.Rename
		err := os.Rename(oldFileName, newFileName)
		if err != nil {
			// Print an error message if there is an error renaming the file
			fmt.Println(err)
			// Return without renaming the rest of the files
			return
		}
	}
}

// get all the keys in YAML files that have the value "__IMPORT"
func GetImportFileNames() []string {
	// Initialize an empty slice of strings to store the import keys
	var importKeys []string
	// Loop through all the keys in the Viper configuration
	for _, key := range viper.AllKeys() {
		// Check if the value for the key is equal to "__IMPORT"
		if viper.GetString(key) == "__IMPORT" {
			// If the value is equal to "__IMPORT", append the key to the importKeys slice
			importKeys = append(importKeys, key)
		}
	}
	// Return the slice of import keys
	return importKeys
}

// SaveImportData function to save the imported data to new files
func saveImportData(importData []ImportData) {
	// Loop through the slice of ImportData
	for _, data := range importData {
		// Read the file to be imported
		file, err := os.ReadFile(data.FileNameToImport)
		if err != nil {
			// Print an error message if there is an error reading the file
			fmt.Println(err)
		}
		// Write the imported data to a new file using ioutil.WriteFile
		err = os.WriteFile(data.FileName, file, 0666)
		if err != nil {
			// Print an error message if there is an error writing the data to the new file
			fmt.Println(err)
		}
	}
}

func getFileNameAndKeys(keys []string) ([]string, []ImportData) {
	// create a splice of ImportData struct
	var importData []ImportData
	// create a slice of strings that contains the filenames
	var fileNames []string

	// Loop through all the keys
	for _, key := range keys {
		// Check if the key contains ".yml"
		if strings.Contains(key, ".yml") {
			// Split the key by ".yml" and extract the file name
			splitYml := strings.Split(key, ".yml")
			fileName := splitYml[0] + ".yml"
			// Check if the file name is not already in the fileNames array
			if !slices.Contains(fileNames, key) {
				// Add the file name to the fileNames array
				fileNames = append(fileNames, fileName)
			}
		}

		// Check if the key starts with "."
		if strings.HasPrefix(key, ".") {
			// Check if the key is not already in the fileNames array
			if !slices.Contains(fileNames, key) {
				// Add the key to the fileNames array
				fileNames = append(fileNames, key)
			}
		}
		// check if the key contains "__import"
		if strings.Contains(key, "__import") {
			// split the key by ".__import" and extract the file name on each config
			savedFileName := strings.Split(key, ".__import")[0]
			fileNameToImport := viper.GetString(key)
			newImportData := ImportData{FileName: savedFileName, FileNameToImport: fileNameToImport}
			// Check if the newImportData is not already in the importData array
			if !slices.Contains(importData, newImportData) {
				importData = append(importData, newImportData)
			}
		}
	}
	return fileNames, importData
}

// handle splitting, pass in  keys with the file name string of each config
func splitFile(keys []string) interface{} {

	// get the fileName and the data to be imported with IMPORT tag
	fileNames, importData := getFileNameAndKeys(keys)

	// create a map to store the fileName with RAW tag
	var rawFileNames = make(map[string]string)

	// iterate through each file name
	for _, fileName := range fileNames {
		// check if the file name contains __raw
		if strings.Contains(fileName, "__raw") {
			configValues := viper.GetString(fileName)
			fileName = strings.Replace(fileName, "__raw", "", 1)
			fileName = strings.TrimRight(fileName, ".")
			fileNameTxt := fileName + ".txt"
			saveText(fileNameTxt, configValues)
			rawFileNames[fileNameTxt] = fileName
		}
		// get the config values for each file name
		configValues := viper.Get(fileName)
		if strings.HasSuffix(fileName, ".yml") || strings.HasSuffix(fileName, ".yaml") {
			createFileIfNotExist(fileName)
			saveConfig(fileName, configValues)
		}
		//fmt.Println()

	}
	renameRaw(rawFileNames)
	saveImportData(importData)
	output := viper.AllSettings()
	return output
}

// accepts config keys and filenames as a slice and returns config keys that have the fileNames passed in, in them
func getFilteredKeys(keys []string, fileNames []string) []string {
	// create an empty slice that contains the keys with the fileName passed along with --only flag
	var filteredKeys []string

	// iterate through the onlyFileNames
	for _, fileName := range fileNames {
		// iterate through the keys
		for _, key := range keys {
			if strings.Contains(key, fileName) {
				if !slices.Contains(filteredKeys, key) {
					filteredKeys = append(filteredKeys, key)
				}
			}
		}
	}
	return filteredKeys
}

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split command will split the file",
	Run: func(cmd *cobra.Command, args []string) {

		// check to see if the user pass in the yaml file to split
		if len(args) < 1 {
			// if it's not return the message below to the user
			fmt.Println("Please provide a yaml file to split")
			return
		}

		// if the --only flag is passed in, we want to get all the filenames
		var onlyFileNames []string
		// if there is --only flag alongside the split command for example
		if only {
			// get the fileNames like database.yml spectra.yml
			onlyFileNames = args[1:]
		}

		// get the Yaml file to split
		fileLocation := args[0]

		onlyFiles := args[1:]
		if len(onlyFileNames) == 0 {
			fmt.Println("Splitted: [all]", "\nFrom:", fileLocation, "\n-------------------")
		} else {
			fmt.Println("Splitted:", onlyFiles, "\nFrom:", fileLocation, "\n-------------------")
		}

		// yq -r -o=json sampleanchor.yml > config.json
		// run it through the HandleAnchor function then return the fileName
		fileLocationTemp, err := handleAnchor(args[0])
		// If there's an error while handling Anchor
		if err != nil {
			// exit the program with the error
			log.Fatalln(err)
		}

		// read the YAML file we want to split
		viper.SetConfigFile(fileLocationTemp)
		// set the config type to yaml
		viper.SetConfigType("yaml")
		// read the configuration
		err = viper.ReadInConfig()
		// if an error occured while reading the file
		if err != nil {
			fmt.Println(err)
			return
		}

		/* err = os.Remove(fileLocationTemp)
		if err != nil {
			// exit the program with the error
			log.Fatalln(err)
		}
		*/
		// get all the keys in the YAML file
		keys := viper.AllKeys()

		// if the only one filename is passed in along with --only flag
		if len(onlyFileNames) == 1 {
			// get the config keys that have onlyFileName in them
			filteredKeys := getFilteredKeys(keys, onlyFileNames)

			splitFile(filteredKeys)
			// Print it to the terminal as stated in the requirements
			// fmt.Println(configValues)

			cmd := exec.Command("yq", "eval-all", onlyFileNames[0])
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(out.String())

			return
		} else if len(onlyFileNames) > 1 {
			// if there is more than one fileNames passed in along the --only flag

			// get the config keys that have onlyFileNames in them
			filteredKeys := getFilteredKeys(keys, onlyFileNames)

			splitFile(filteredKeys)

			return

		} else {
			// if there is no --only flag split the file with all the keys
			splitFile(keys)
		}
	},
}

func init() {
	// add 'split' command to the 'splityaml' base command
	rootCmd.AddCommand(splitCmd)
	// add the --only flag to the split command
	splitCmd.Flags().BoolVarP(&only, "only", "o", false, "split only some config")
	// splitYmlCmd.PersistentFlags().String("only", "", "Specify some YAML file to extract")
}
