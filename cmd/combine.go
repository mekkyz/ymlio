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

// setting up the combine functionality
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "combine multiple YAML files",
	Run: func(cmd *cobra.Command, args []string) {

		// getting all arguments passed in from the command line after the 'combine' command
		args = args[:]

		// get the last file name passed in which the files will be combinen
		lastArgsFileName := args[len(args)-1]

		// open the last file passed in the command line for saving the combined file
		file, err := os.OpenFile(lastArgsFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		// error handling
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		// get all the file names to combine excepts the last file which contains the combine everything
		args = args[:len(args)-1]

		// iterate through all the filenames that we want to combine
		fmt.Printf("Combined: ")
		for _, fileName := range args {
			// checking if the file ends with .yml file, meaning if it's a yaml file
			if strings.Contains(fileName, ".yml") && strings.HasSuffix(fileName, ".yml") || strings.Contains(fileName, ".yaml") && strings.HasSuffix(fileName, ".yaml") {
				// handle the anchors first
				fileNameAnchor, err := handleAnchor(fileName)
				if err != nil {
					log.Fatalln(err)
				}
				// read the file content after handling the anchors
				fileText := readFile(fileNameAnchor)
				// create the config text by adding filename as keys and their content
				config := fmt.Sprintf("%s:\n  %s\n", fileName, fileText)
				// write the content to the last file we passed in
				_, err = file.WriteString(config)
			} else {

				// if the filename isn't a yaml file, in this case it's a text file, read the content
				fileText := readFile(fileName)

				// create the config file and adding the __RAW tag
				config := fmt.Sprintf("%s:\n  __RAW: |\n    %s\n", fileName, strings.Replace(fileText, "\n", "\n    ", -1))

				// write the content to the last file we passed in
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
