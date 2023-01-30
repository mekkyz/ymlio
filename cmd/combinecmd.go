package cmd

import (
	"fmt"
	"log"
	"yaml-processing/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

		fileNameHanchorFileName, err := utils.HandleAnchor(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		viper.Set(fileNameHanchorFileName, "")
		viper.SetConfigType("yml")
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatalln(err)
		}

		for _, fileName := range args[:len(args)-1] {
			fileNameHanchorFileName, err = utils.HandleAnchor(fileName)

			if err != nil {
				log.Fatalln(err)
			}
			viper.Set(fileNameHanchorFileName, viper.AllSettings())
			viper.SetConfigFile(fileNameHanchorFileName)
			viper.SetConfigType("yml")
			err = viper.MergeInConfig()
			if err != nil {
				log.Fatalln(err)
			}

		}

		viper.WriteConfigAs(lastArgsFileName)
		fmt.Println("Done")
	},
}

func init() {
	RootCmd.AddCommand(combineCmd)
}
