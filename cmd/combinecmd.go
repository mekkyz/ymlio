package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"yaml-processing/utils"

	"github.com/spf13/cobra"
)

func readFile(fileName string) string {
	fileText, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	return string(fileText)
}

// This function sets up the combine functionality
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "combine multiple YAML files",
	Run: func(cmd *cobra.Command, args []string) {

		// getting all arguments passed in from the command line after the 'combine' command
		args = args[:]

		// get the last file name passed in which the file will be combined, for example
		// if 'splityaml combine spectra.yml storage.yml hello.yml' is passed in, the lastArgsFileName below will get the hello.yml
		lastArgsFileName := args[len(args)-1]

		// open the last file passed in the command line for saving the combined file for example
		// splityaml combine a.yml b.yml c.yml d.yml
		// This will get the d.yml file and open it for read and writthin
		file, err := os.OpenFile(lastArgsFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		// if there's an error opening the file
		if err != nil {
			// exit and print the error
			log.Fatalln(err)
		}
		defer file.Close()

		// get all the file names to combine excepts the last file which contains the combine everything
		args = args[:len(args)-1]

		// iterate through all the filenames that we want to combine
		for _, fileName := range args {

			// checking if the file ends with .yml file, meaning if it's a yaml file
			if strings.Contains(fileName, ".yml") && strings.HasSuffix(fileName, ".yml") {
				// handle the anchors first
				fileNameHanchor, err := utils.HandleAnchor(fileName)
				if err != nil {
					log.Fatalln(err)
				}
				// read the file content after handling the anchors
				fileText := readFile(fileNameHanchor)
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

			fmt.Println(fileName)
		}

		fmt.Println("Done")

	},
}

func init() {
	RootCmd.AddCommand(combineCmd)
}
