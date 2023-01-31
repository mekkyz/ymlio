/*
Copyright © 2023 Mostafa Mekky <mos.mekky@gmail.com>
*/
package main

import (
	"fmt"
	"yaml-processing/cmd"
)

func main() {
	// executes the root command so we can be able to use the splityaml commands
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
