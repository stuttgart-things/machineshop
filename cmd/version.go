/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

var (
	shortened  = false
	version    = "unset"
	commit     = "unknown"
	date       = "unknown"
	output     = "yaml"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "version will output the current build information",
		Long: `Print the version information. For example:
	sthings version`,

		Run: func(_ *cobra.Command, _ []string) {
			resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
			fmt.Print(resp)
		},
	}
)

func init() {
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	versionCmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of 'yaml' or 'json'.")
	rootCmd.AddCommand(versionCmd)
}
