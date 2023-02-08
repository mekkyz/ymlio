package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func readFile(fileName string) string {
	fileText, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	return string(fileText)
}

// Setting up the combine functionality
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "combine multiple YAML files",
	Run: func(cmd *cobra.Command, args []string) {
		// Getting all arguments passed in from the command line after the 'combine' command
		args = args[:]
		// Get the last file name passed in which the files will be combined
		lastArgsFileName := args[len(args)-1]
		// Open the last file passed in the command line for saving the combined file
		file, err := os.OpenFile(lastArgsFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		// Get all the file names to combine excepts the last file which contains the combine everything
		args = args[:len(args)-1]
		// Iterate through all the filenames that we want to combine
		fmt.Printf("Combined: ")
		for _, fileName := range args {
			// Checking if the file ends with .yml or .yaml
			if strings.Contains(fileName, ".yml") && strings.HasSuffix(fileName, ".yml") || strings.Contains(fileName, ".yaml") && strings.HasSuffix(fileName, ".yaml") {
				// Handle the anchors first
				fileNameAnchor, err := handleAnchor(fileName)
				if err != nil {
					log.Fatalln(err)
				}
				// Read the file content after handling the anchors
				fileText := readFile(fileNameAnchor)
				// Create the config text by adding filename as keys and their content
				config := fmt.Sprintf("%s:\n  %s\n", fileName, fileText)
				// Write the content to the last file we passed in
				_, err = file.WriteString(config)
			} else {
				// If the filename isn't a yaml file, in this case it's a text file, read the content
				fileText := readFile(fileName)
				// Create the config file and adding the __RAW tag
				config := fmt.Sprintf("%s:\n  __RAW: |\n    %s\n", fileName, strings.Replace(fileText, "\n", "\n    ", -1))
				// Write the content to the last file we passed in
				_, err = file.WriteString(config)
			}
			fmt.Printf(fileName + " ")
		}
		fmt.Printf("\nTo: " + lastArgsFileName + "\n-------------------\n")
	},
}

func init() {
	rootCmd.AddCommand(combineCmd)
}
