/*
Copyright © 2023 Mostafa Mekky <mos.mekky@gmail.com>
*/
package main

import (
	"ymlio/cmd"
)

func main() {
	cmd.Execute()
	// executes the root command to be able to use the ymlio commands
	// err := cmd.rootCmd.Execute()
	// if err != nil {
	//	fmt.Println(err)
	//}
}
