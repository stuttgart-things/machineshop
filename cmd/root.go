/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	gitRepository string
)

const banner = `
.___  ___.      ___       ______  __    __   __  .__   __.  _______       _______. __    __    ______   .______
|   \/   |     /   \     /      ||  |  |  | |  | |  \ |  | |   ____|     /       ||  |  |  |  /  __  \  |   _  \
|  \  /  |    /  ^  \   |  ,----'|  |__|  | |  | |   \|  | |  |__       |   (---- |  |__|  | |  |  |  | |  |_)  |
|  |\/|  |   /  /_\  \  |  |     |   __   | |  | |  .    | |   __|       \   \    |   __   | |  |  |  | |   ___/
|  |  |  |  /  _____  \ |   ----.|  |  |  | |  | |  |\   | |  |____.------)   |   |  |  |  | |   --'  | |  |
|__|  |__| /__/     \__\ \______||__|  |__| |__| |__| \__| |_______|_________/    |__|  |__|  \______/  | _|

`

var rootCmd = &cobra.Command{
	Use:   "machineShop",
	Short: "machineShop - infrasturcture cli",
	Long:  `cli for managing infrastructure automation`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// fmt.Println(banner)
	// color.Blue(banner)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&gitRepository, "git", "https://github.com/stuttgart-things/stuttgart-things.git", "iac git repository")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
