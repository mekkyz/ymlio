/*
Copyright © 2023 Mostafa Mekky <mos.mekky@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ymlio",
	Short: "ymlio splits and combines yaml files",
	Long:  `Ymlio is a CLI Tool that allows users to easily combine or split yaml files.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// This command has the sole role to make the tool show the help page upon typing the base command
var blankCmd = &cobra.Command{
	Use:    "",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.Help()
	},
}

func Execute() {
	rootCmd.SetHelpCommand(blankCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
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
