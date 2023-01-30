package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"yaml-processing/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

// a struct to hold the filename on the config keys, and also if it has import flag, get the filename that will be imported
type ImportData struct {
	FileName         string
	FileNameToImport string
}

var only bool

// This function create a file, if it doesn't exisit already
// pass in the filepath as a string
func CreateFileIfNotExist(filePath string) {

	currentPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	fullPath := filepath.Join(currentPath, filePath)
	fmt.Println(fullPath)

	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	}
}

// This function save the YAML config file, by acception the filename to be saved and also the values inside the config to be saved
func SaveConfig(savedFileName string, yamlValues interface{}) {
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
	err = ioutil.WriteFile(savedFileName, []byte(converterBytes), 0644)
	if err != nil {
		fmt.Println(err)
		// Return if there is an error
		return
	}
}

// This function to save a text file with a given filename and content
func SaveText(fileName string, content string) {
	// Create a new file with the given filename using os.Create
	file, err := os.Create(fileName)
	if err != nil {
		// Print an error message if there is an error creating the file
		fmt.Println("Error creating result.txt file:", err)
		// Exit the program with a status code of 1
		os.Exit(1)
	}
	// Defer the closing of the file until all other function calls have completed
	defer file.Close()

	// Write the content string to the file using file.WriteString
	_, err = file.WriteString(content)
	if err != nil {
		// Print an error message if there is an error writing to the file
		fmt.Println("Error writing to result.txt file:", err)
		// Return without writing to the file
		return
	}
}

// RenameRaw function to rename multiple raw files with new filenames
func RenameRaw(rawFileNames map[string]string) {
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

// This function to get all the keys in YAML files that have the value "__IMPORT"
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
func SaveImportData(importData []ImportData) {
	// Loop through the slice of ImportData
	for _, data := range importData {
		// Read the file to be imported using ioutil.ReadFile
		file, err := ioutil.ReadFile(data.FileNameToImport)
		if err != nil {
			// Print an error message if there is an error reading the file
			fmt.Println(err)
		}
		// Write the imported data to a new file using ioutil.WriteFile
		err = ioutil.WriteFile(data.FileName, file, 0666)
		if err != nil {
			// Print an error message if there is an error writing the data to the new file
			fmt.Println(err)
		}
	}
}

func GetFileNameAndKeys(keys []string) ([]string, []ImportData) {
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
		// Check if the key contains "__import"
		if strings.Contains(key, "__import") {
			// Split the key by ".__import" and extract the file name on each config
			savedFileName := strings.Split(key, ".__import")[0]
			// Get the content of the config with the key, this will contain the fileName which is to be imported
			fileNameToImport := viper.GetString(key)
			// Create a new ImportData with the saved file name and file name to import
			newImportData := ImportData{FileName: savedFileName, FileNameToImport: fileNameToImport}
			// Check if the newImportData is not already in the importData array
			if !slices.Contains(importData, newImportData) {
				// Add the newImportData to the importData array
				importData = append(importData, newImportData)
			}
		}
	}
	// Return the extracted file names and the import data
	return fileNames, importData
}

// This function handles the splitting functionality, we pass in the keys which contains the string that has the
// the filename of each config for example if we have a YAML file with a config name like
// spectra.yml:
// production:
// port: 2000
// The function accepts a slice that contains spectra.yml
func SplitFile(keys []string) string {

	// get the fileName and the data to be imported with IMPORT tag
	fileNames, importData := GetFileNameAndKeys(keys)

	// creating a variables that will hold the configuration values
	var configValues string

	// Create a map to store the fileName with RAW tag
	var rawFileNames = make(map[string]string)

	// Iterate through each file name//
	for _, fileName := range fileNames {

		// Check if the file name contains __raw
		if strings.Contains(fileName, "__raw") {
			// Get the values of the config
			configValues := viper.GetString(fileName)
			// Replace __raw in the file name
			fileName = strings.Replace(fileName, "__raw", "", 1)
			// Trim the right side of the file name
			fileName = strings.TrimRight(fileName, ".")
			// Add .txt extension to the file name
			fileNameTxt := fileName + ".txt"
			// Save the config values to the text file
			SaveText(fileNameTxt, configValues)
			// Add the old and new names of the file to the map
			rawFileNames[fileNameTxt] = fileName
		}

		// Get the config values for each file name
		configValues = viper.GetString(fileName)
		// Check if the file name ends with .yml
		if strings.HasSuffix(fileName, ".yml") {
			// Check if the file already exists and create it if it doesn't
			CreateFileIfNotExist(fileName)
			// Save the config values to the yml file
			SaveConfig(fileName, configValues)
		}

	}

	RenameRaw(rawFileNames)
	SaveImportData(importData)
	return configValues
}

// This function accepts config keys and filenames as a slice and returns config keys that have the fileNames passed in, in them
func GetFilteredKeys(keys []string, fileNames []string) []string {
	// create an empty slice that contains the keys with the fileName passed along with --only flag
	var filteredKeys []string

	// iterate through the onlyFileNames
	for _, fileName := range fileNames {
		// iterate through the keys
		for _, key := range keys {
			// if the key contains the fileNames
			if strings.Contains(key, fileName) {
				// if the key is not already in the filteredKeys slice
				if !slices.Contains(filteredKeys, key) {
					// append or add it to the filteredKey slice
					filteredKeys = append(filteredKeys, key)
				}
			}
		}
	}
	return filteredKeys
}

var splitYmlCmd = &cobra.Command{
	Use:   "split",
	Short: "split command will split the file",
	Run: func(cmd *cobra.Command, args []string) {

		// check to see if the user pass in the yaml file to split
		if len(args) < 1 {
			// if it's not return the message below to the user
			fmt.Println("Please provide a YAML file to split")
			return
		}

		// if the --only flag is passed in, we want to get all the filenames
		var onlyFileNames []string
		// if there is --only flag alongside the split command for example
		// splityaml split bigyml.yml --only database.yml spectra.yml
		if only {
			// get the fileNames like database.yml spectra.yml
			onlyFileNames = args[1:]
		}

		// get the Yaml file to split
		fileLocation := args[0]
		fmt.Println("Splitting this file =>", fileLocation)
		// yq -r -o=json sampleanchor.yml > config.json
		// run it through the HandleAnchor function then return the fileName
		fileLocation, err := utils.HandleAnchor(fileLocation)
		// If there's an error while handling Anchor
		if err != nil {
			// exit the program with the error
			log.Fatalln(err)
		}

		// read the YAML file we want to split
		viper.SetConfigFile(fileLocation)
		// set the config type to yaml
		viper.SetConfigType("yaml")
		// read the configuration
		err = viper.ReadInConfig()
		// if an error occured while reading the file
		if err != nil {
			fmt.Println(err)
			return
		}

		// get all the keys in the YAML file
		keys := viper.AllKeys()

		// if the only one filename is passed in along with --only flag
		if len(onlyFileNames) == 1 {
			fmt.Println(onlyFileNames)

			// get the config keys that have onlyFileName in them
			filteredKeys := GetFilteredKeys(keys, onlyFileNames)

			// get the values
			configValues := SplitFile(filteredKeys)
			// Print it to the terminal as stated in the requirements
			fmt.Println(configValues)
			fmt.Println("Done")
			return
		}

		// if there is more than one fileNames passed in along the --only flag
		if len(onlyFileNames) > 1 {
			fmt.Println(onlyFileNames)

			// get the config keys that have onlyFileNames in them
			filteredKeys := GetFilteredKeys(keys, onlyFileNames)

			SplitFile(filteredKeys)
			fmt.Println("Done")
			return

		} else {
			// if there is no --only flag split the file with all the keys

			SplitFile(keys)
			fmt.Println("Done")
		}

	},
}

func init() {
	// add 'split' command to the 'splityaml' base command
	RootCmd.AddCommand(splitYmlCmd)
	// add the --only flag to the split command
	splitYmlCmd.Flags().BoolVarP(&only, "only", "o", false, "split only some config")
	// splitYmlCmd.PersistentFlags().String("only", "", "Specify some YAML file you want to extract")
}
