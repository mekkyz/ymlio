/*
Copyright © 2023 Mostafa Mekky <mos.mekky@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ymlio",
	Short: "ymlio splits and combines yaml files",
	Long: `
ymlio is a CLI Application that allow users to easily split or combine yaml files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nHello, welcome to ymlio: The yml file tool.\n\nymlio must be run with subcommands; do --help for more\n\n")

	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	markdownFile := "ymlio-docs.md"
	if _, err := os.Stat(markdownFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s not found, \ngenerating...\n", markdownFile)
		// check if tmp_ymlio_doc folder exists.
		path := "./tmp_ymlio_doc"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, 0700)
			if err != nil {
				log.Fatal(err)
			}
		}

		// generate markdownTree for all cobra.commands.
		err := doc.GenMarkdownTree(rootCmd, "./tmp_ymlio_doc")
		if err != nil {
			log.Fatal(err)
		}
	}

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
/*
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yaml-processing.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
*/
