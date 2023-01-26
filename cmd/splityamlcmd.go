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

type ImportData struct {
	FileName         string
	FileNameToImport string
}

var only bool

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

func SaveConfig(savedFileName string, yamlValues interface{}) {
	var converterBytes []byte
	var err error
	if reflect.TypeOf(yamlValues).Kind() == reflect.Map {
		converterBytes, err = yaml.Marshal(yamlValues)
		if err != nil {
			fmt.Println("Cannot convert converter to []byte: ", err)
			return
		}
	}
	err = ioutil.WriteFile(savedFileName, []byte(converterBytes), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SaveText(fileName string, content string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating result.txt file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write the value of __RAW to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to result.txt file:", err)
		return
	}
}

func RenameRaw(rawFileNames map[string]string) {
	for oldFileName, newFileName := range rawFileNames {
		err := os.Rename(oldFileName, newFileName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func GetImportFileNames() []string {
	var importKeys []string
	for _, key := range viper.AllKeys() {
		if viper.GetString(key) == "__IMPORT" {
			importKeys = append(importKeys, key)
		}
	}
	return importKeys
}

func SaveImportData(importData []ImportData) {
	for _, data := range importData {

		file, err := ioutil.ReadFile(data.FileNameToImport)
		if err != nil {
			fmt.Println(err)
		}

		err = ioutil.WriteFile(data.FileName, file, 0666)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func GetFileNameAndKeys(keys []string) ([]string, []ImportData) {
	var importData []ImportData
	var fileNames []string

	for _, key := range keys {
		if strings.Contains(key, ".yml") {
			splitYml := strings.Split(key, ".yml")
			fileName := splitYml[0] + ".yml"
			if !slices.Contains(fileNames, key) {
				fileNames = append(fileNames, fileName)
			}
		}

		if strings.HasPrefix(key, ".") {

			if !slices.Contains(fileNames, key) {
				fileNames = append(fileNames, key)
				// fmt.Println(key)
			}
		}
		// handling import file
		if strings.Contains(key, "__import") {

			savedFileName := strings.Split(key, ".__import")[0]

			fileNameToImport := viper.GetString(key)
			newImportData := ImportData{FileName: savedFileName, FileNameToImport: fileNameToImport}
			if !slices.Contains(importData, newImportData) {
				importData = append(importData, newImportData)
			}

		}
	}
	return fileNames, importData
}

func SplitFile(keys []string) {

	fileNames, importData := GetFileNameAndKeys(keys)

	// fmt.Printf("%+v\n", importData)
	// fmt.Println()
	// fmt.Printf("%+v\n", fileNames)
	// os.Exit(1)

	var rawFileNames = make(map[string]string)

	for _, fileName := range fileNames {

		if strings.Contains(fileName, "__raw") {
			configValues := viper.GetString(fileName)
			fileName = strings.Replace(fileName, "__raw", "", 1)
			fileName = strings.TrimRight(fileName, ".")
			fileNameTxt := fileName + ".txt"
			SaveText(fileNameTxt, configValues)

			rawFileNames[fileNameTxt] = fileName
		}

		configValues := viper.Get(fileName)
		if strings.HasSuffix(fileName, ".yml") {
			CreateFileIfNotExist(fileName)
			SaveConfig(fileName, configValues)
		}
		fmt.Println()

	}

	RenameRaw(rawFileNames)
	SaveImportData(importData)
}

var splitYmlCmd = &cobra.Command{
	Use:   "split",
	Short: "split command will split the file",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("Please provide a YAML file to split")
			return
		}

		var onlyFileNames []string
		if only {
			onlyFileNames = args[1:]
		}

		fileLocation := args[0]
		fmt.Println("Splitting this file =>", fileLocation)
		// yq -r -o=json sampleanchor.yml > config.json
		fileLocation, err := utils.HandleAnchor(fileLocation)
		if err != nil {
			log.Fatalln(err)
		}

		viper.SetConfigFile(fileLocation)
		viper.SetConfigType("yaml")
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		keys := viper.AllKeys()

		if len(onlyFileNames) > 1 {
			fmt.Println(onlyFileNames)

			var filteredKeys []string

			for _, oName := range onlyFileNames {
				for _, key := range keys {
					if strings.Contains(key, oName) {
						if !slices.Contains(filteredKeys, key) {
							filteredKeys = append(filteredKeys, key)
						}
					}
				}
			}

			SplitFile(filteredKeys)
			fmt.Println("Done")
			return

		} else {

			SplitFile(keys)
			fmt.Println("Done")
		}

	},
}

func init() {
	RootCmd.AddCommand(splitYmlCmd)
	splitYmlCmd.Flags().BoolVarP(&only, "only", "o", false, "split only some config")
	// splitYmlCmd.PersistentFlags().String("only", "", "Specify some YAML file you want to extract")
}
