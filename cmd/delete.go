/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/google/go-github/v68/github"
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
		user, _ := cmd.Flags().GetString("user")
		token, _ := cmd.Flags().GetString("token")
		filesToRemove, _ := cmd.Flags().GetStringSlice("files")
		commitMessage, _ := cmd.LocalFlags().GetString("message")

		// IF TOKEN IS NOT PROVIDED, TRY TO GET IT FROM ENVIRONMENT
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}

		// IF NOT DEFINED IN ENVIRONMENT OR FLAG, EXIT
		if token == "" {
			log.Error("GITHUB TOKEN NOT FOUND")
			os.Exit(3)
		}

		switch kind {

		case "branch":

			// CREATE GITHUB CLIENT
			client = github.NewClient(nil).WithAuthToken(token)

			// DELETE BRANCH
			branchPreFix := "refs/heads/"

			sthingsCli.DeleteBranch(client, repositoryName, groupName, branchPreFix+branchName)

		case "files":

			// IF FILES ARE EMPTY TO NOTHING
			if len(filesToRemove) != 0 {

				// IF TOKEN IS NOT PROVIDED, TRY TO GET IT FROM ENVIRONMENT
				if user == "" {
					user = os.Getenv("GIT_USER")
				}

				// IF MESSAGE IS EMPTY SET IT TO "DELETED FILES W/ machineshop"
				if commitMessage == "" {
					commitMessage = "Deleted files w/ machineshop"
				}

				// CREATE GIT AUTH
				auth := sthingsCli.CreateGitAuth(user, token)

				// SET GIT REPO URL
				gitRepository := "https://github.com/" + groupName + "/" + repositoryName + ".git"

				// DELETE FILES FROM BRANCH
				sthingsCli.AddCommitFileToGitRepository(gitRepository, branchName, auth, nil, filesToRemove, commitMessage)

			} else {
				log.Error("FILES TO DELETE NOT PROVIDED")
				os.Exit(3)
			}

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
	deleteCmd.Flags().String("user", "", "git user")
	deleteCmd.Flags().StringSlice("files", []string{}, "files to be deleted")
	deleteCmd.Flags().String("message", "", "commit message")
}
