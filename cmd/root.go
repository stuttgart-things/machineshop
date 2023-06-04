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
	gitUser       string
	gitToken      string
	logFilePath   string
)

var rootCmd = &cobra.Command{
	Use:   "machineShop",
	Short: "machineShop - infrasturcture cli",
	Long:  `cli for managing infrastructure automation`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&gitRepository, "git", "https://github.com/stuttgart-things/stuttgart-things.git", "iac git repository")
	rootCmd.PersistentFlags().StringVar(&logFilePath, "log", "/tmp/machineshop.log", "log file path")
	rootCmd.PersistentFlags().StringVar(&gitUser, "gitUser", "git/data/github:username", "git user")
	rootCmd.PersistentFlags().StringVar(&gitToken, "gitToken", "git/data/github:token", "git token")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
