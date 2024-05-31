/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/google/go-github/v62/github"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete things",
	Long:  `Delete things on remote systems`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		kind, _ := cmd.LocalFlags().GetString("kind")
		repositoryName, _ := cmd.LocalFlags().GetString("repository")
		groupName, _ := cmd.LocalFlags().GetString("group")
		branchName, _ := cmd.LocalFlags().GetString("branch")
		token, _ := cmd.LocalFlags().GetString("token")

		switch kind {

		case "branch":

			// IF TOKEN IS NOT PROVIDED, TRY TO GET IT FROM ENVIRONMENT
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}

			// IF NOT DEFINED IN ENVIRONMENT OR FLAG, EXIT
			if token == "" {
				log.Error("GITHUB TOKEN NOT FOUND")
				os.Exit(3)
			}

			// CREATE GITHUB CLIENT
			client = github.NewClient(nil).WithAuthToken(token)

			// DELETE BRANCH
			branchPreFix := "refs/heads/"

			sthingsCli.DeleteBranch(client, repositoryName, groupName, branchPreFix+branchName)
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().String("kind", "branch", "kind of operation to perform")
	deleteCmd.Flags().String("group", "stuttgart-things", "name of group")
	deleteCmd.Flags().String("repository", "stuttgart-things", "name of repository")
	deleteCmd.Flags().String("branch", "", "(to be deleted) branch name")
	deleteCmd.Flags().String("token", "", "github token")
}
