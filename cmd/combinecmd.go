package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"yaml-processing/utils"

	"github.com/spf13/cobra"
)

var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "combine multiple YAML files",
	Run: func(cmd *cobra.Command, args []string) {

		args = args[:]

		// get the last file passed in
		lastArgs := args[len(args)-1]

		file, err := os.OpenFile(lastArgs, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		args = args[:len(args)-1]

		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		for _, fileName := range args {
			fileNameHanchor, err := utils.HandleAnchor(fileName)
			if err != nil {
				log.Fatalln(err)
			}
			fileText, err := ioutil.ReadFile(fileNameHanchor)
			if err != nil {
				fmt.Println(err)
			}

			config := fmt.Sprintf("%s:\n\t%s\n\n", fileName, fileText)
			_, err = file.WriteString(config)
			fmt.Println(fileName)
		}

		fmt.Println("Done")
	},
}

func init() {
	RootCmd.AddCommand(combineCmd)
}
