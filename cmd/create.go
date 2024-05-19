/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v62/github"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var client *github.Client
var ctx = context.Background()

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create things",
	Long:  `Create things on remote systems`,
	Run: func(cmd *cobra.Command, args []string) {
		// FLAGS
		kind, _ := cmd.LocalFlags().GetString("kind")

		switch kind {

		case "pr":

			fmt.Println("HELLO")
			token := os.Getenv("GITHUB_TOKEN")
			if token == "" {
				log.Fatal("UNAUTHORIZED: NO TOKEN PRESENT")
			}

			client = github.NewClient(nil).WithAuthToken(token)

			// CALL GETREFERENCEOBJECT
			ref, err := sthingsCli.GetReferenceObject(client, "stuttgart-things", "machineshop", "test", "main")
			fmt.Println(ref, err)
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String("kind", "pr", "kind of operation to perform")

}
