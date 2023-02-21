/*
Copyright © 2023 Mostafa Mekky <mos.mekky@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

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
	Short: "Combine files (if files are not yaml, it will mark them with __RAW) into one yaml file.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			// If it's not return the message below to the user
			fmt.Println("Please provide at least 2 yaml files to combine and a file name to combine the content in.\nYou need at least 3 arguments after combine.\nFor example:\nymlio combine a.yml b.yml c.yml -> This will combine the content of a.yml and b.yml into a file c.yml")
			return
		}
		// Get the last file name passed in which the files will be combined
		lastArgsFileName := args[len(args)-1]
		// Open the last file passed in the command line for saving the combined file
		file, err := os.OpenFile(lastArgsFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		// Create a map to hold all the YAML documents
		allDocs := make(map[string]interface{})

		// Iterate through all the filenames that we want to combine
		fmt.Printf("Combined: ")
		for _, fileName := range args[:len(args)-1] {
			// Checking if the file ends with .yml or .yaml
			if strings.Contains(fileName, ".yml") && strings.HasSuffix(fileName, ".yml") || strings.Contains(fileName, ".yaml") && strings.HasSuffix(fileName, ".yaml") {
				// Handle the anchors first
				fileNameAnchor, err := handleAnchor(fileName)
				if err != nil {
					log.Fatalln(err)
				}
				// Read the file content after handling the anchors
				fileText := readFile(fileNameAnchor)
				// Parse the file content into a map[string]interface{}
				var doc map[string]interface{}
				if err := yaml.Unmarshal([]byte(fileText), &doc); err != nil {
					log.Fatalf("failed to parse %q: %v", fileName, err)
				}
				// Add the document to the allDocs map with the file name as the key
				allDocs[fileName] = doc
			} else {
				// If the filename isn't a yaml file, in this case it's a text file, read the content
				fileText := readFile(fileName)
				// Add the content as a string to the allDocs map with the file name as the key
				allDocs[fileName] = fileText
			}
			fmt.Printf(fileName + " ")
		}

		// Write the combined YAML documents to the output file
		encoder := yaml.NewEncoder(file)
		if err := encoder.Encode(allDocs); err != nil {
			log.Fatalf("failed to encode YAML: %v", err)
		}

		fmt.Println("\nTo: " + lastArgsFileName)
	},
}

func init() {
	rootCmd.AddCommand(combineCmd)
}
