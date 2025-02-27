/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"os"
	"time"

	github "github.com/google/go-github/v68/github"
	"github.com/spf13/cobra"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var client *github.Client

// GET CURRENT TIME
var now = time.Now()
var timeString = now.Format("06-01-02-15-04")

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create things",
	Long:  `Create things on remote systems`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		kind, _ := cmd.LocalFlags().GetString("kind")
		groupName, _ := cmd.LocalFlags().GetString("group")
		repositoryName, _ := cmd.LocalFlags().GetString("repository")
		branchName, _ := cmd.LocalFlags().GetString("branch")
		baseBranch, _ := cmd.LocalFlags().GetString("base")
		token, _ := cmd.LocalFlags().GetString("token")
		files, _ := cmd.Flags().GetStringSlice("files")
		authorName, _ := cmd.LocalFlags().GetString("author")
		authorEmail, _ := cmd.LocalFlags().GetString("email")
		commitMessage, _ := cmd.LocalFlags().GetString("message")
		prTitle, _ := cmd.LocalFlags().GetString("title")
		privateRepository, _ := cmd.LocalFlags().GetBool("private")
		mergeMethod, _ := cmd.LocalFlags().GetString("merge")
		prID, _ := cmd.LocalFlags().GetString("id")
		labels, _ := cmd.Flags().GetStringSlice("labels")

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

		log.Info("GROUP: ", groupName)
		log.Info("REPOSITORY: ", repositoryName)
		log.Info("BRANCH: ", branchName)
		log.Info("BASE-BRANCH: ", baseBranch)

		switch kind {

		case "repo":

			if commitMessage == "" {
				commitMessage = repositoryName
			}

			// HARDCODED AUTO INIT FOR NOW
			autoInit := true

			log.Info("DESCRIPTION: ", commitMessage)
			log.Info("PRIVATE REPOSITORY: ", privateRepository)
			log.Info("AUTO INIT: ", autoInit)

			err, repoName := sthingsCli.CreateRepository(client, repositoryName, commitMessage, groupName, privateRepository, autoInit)
			if err != nil {
				log.Error("UNABLE TO CREATE REPOSITORY: ", err)
			} else {
				log.Info("REPOSITORY CREATED: ", repoName)
			}

		case "branch":

			// IF BRANCH IS NOT PROVIDED, CREATE ONE RANDOM NAME
			if branchName == "" {
				branchName = "machineshop-" + timeString
			}

			// IF COMMIT IS NOT PROVIDED, CREATE ONE RANDOM NAME
			if commitMessage == "" {
				commitMessage = branchName
			}

			log.Info("CREATING BRANCH")
			log.Info("FILES: ", files)
			log.Info("AUTHOR: ", authorName)
			log.Info("EMAIL: ", authorEmail)
			log.Info("MESSAGE: ", commitMessage)

			// GET GIT REFERENCE OBJECT
			ref, err := sthingsCli.GetReferenceObject(client, groupName, repositoryName, branchName, baseBranch)
			if err != nil {
				log.Fatalf("UNABLE TO GET/CREATE THE COMMIT REFERENCE: %s\n", err)
			}
			if ref == nil {
				log.Fatalf("NO ERROR WHERE RETURNED BUT THE REFERENCE IS NIL")
			}

			// CREATE A NEW GIT TREE
			gitTree, err := sthingsCli.GetGitTree(client, ref, files, groupName, repositoryName)
			if err != nil {
				log.Fatalf("UNABLE TO CREATE THE TREE BASED ON THE PROVIDED FILES: %s\n", err)
			}

			fmt.Println(gitTree)

			// PUSH COMMIT
			err = sthingsCli.PushCommit(client, ref, gitTree, groupName, repositoryName, authorName, authorEmail, commitMessage)
			if err != nil {
				log.Fatalf("UNABLE TO PUSH THE COMMIT: %s\n", err)
			}

		case "pr":
			log.Info("CREATING PULL REQUEST")
			log.Info("LABELS", labels)

			// IF KIND EQUALS PR AND TITLE IS NOT PROVIDED
			if prTitle == "" {
				log.Error("PULL REQUEST TITLE IS MISSING - EXITING")
				os.Exit(3)
			}

			// IF BRANCH IS NOT PROVIDED, CREATE ONE RANDOM NAME
			if branchName == "" {
				log.Error("BRANCH NAME IS MISSING - EXITING")
				os.Exit(3)
			}

			// CROSS REPOSITORY PULL REQUESTS ARE NOT SUPPURTED YET IN MACHINESHOP (NO USE CASE SO FAR) -
			// SO WE USE THE SAME REPOSITORY FOR SOURCE AND TARGET FOR NOW
			sourceRepo := repositoryName
			prRepo := repositoryName
			sourceOwner := groupName
			prRepoOwner := groupName
			commitBranch := branchName
			repoBranch := branchName
			prDescription := prTitle
			prSubject := prTitle

			log.Info("PULL-REQUEST TITLE: ", prTitle)

			// CREATE PULL REQUEST
			err, pullRequestID := sthingsCli.CreatePullRequest(client, prSubject, prRepoOwner, sourceOwner, commitBranch, prRepo, sourceRepo, repoBranch, baseBranch, prDescription, labels)
			if err != nil {
				log.Fatalf("UNABLE TO CREATE THE PULL REQUEST: %s\n", err)
			} else {
				log.Info("PULL-REQUEST CREATED W/ ID: ", pullRequestID)
			}

		case "merge":

			log.Info("CREATING MERGE REQUEST")
			log.Info("MERGE: ", commitMessage)
			log.Info("PR-ID: ", prID)
			log.Info("MERGE-METHOD : ", mergeMethod)

			prID, _ := cmd.LocalFlags().GetString("id")

			sthingsCli.MergePullRequest(client, repositoryName, groupName, commitMessage, mergeMethod, sthingsBase.ConvertStringToInteger(prID))

		case "labels":

			log.Info("CREATING/UPDATING LABELS OF PULL REQUEST")
			log.Info("PR-ID: ", prID)
			log.Info("TO BE UPDATED LABELS: ", labels)

			updatedLabels, err := sthingsCli.SetPullRequestLabels(client, groupName, repositoryName, sthingsBase.ConvertStringToInteger(prID), labels)
			if err != nil {
				log.Fatalf("UNABLE TO SET LABELS: %s\n", err)
			}

			log.Info("LABELS UPDATED: ", updatedLabels)
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String("kind", "branch", "kind of operation to perform")
	createCmd.Flags().String("group", "stuttgart-things", "name of group")
	createCmd.Flags().String("repository", "stuttgart-things", "name of repository")
	createCmd.Flags().String("branch", "", "(to be created) branch name")
	createCmd.Flags().String("title", "", "pull request title")
	createCmd.Flags().String("author", "machineshop", "author name")
	createCmd.Flags().String("email", "machineshop@stuttgart-things.com", "author email")
	createCmd.Flags().String("message", "", "commit message")
	createCmd.Flags().String("id", "", "pull reuqest id")
	createCmd.Flags().String("merge", "merge", "merge methhod rebase")
	createCmd.Flags().String("token", "", "github token")
	createCmd.Flags().Bool("private", false, "private repository")
	createCmd.Flags().String("base", "main", "name of (to be merged) branch")
	createCmd.Flags().StringSlice("files", []string{}, "files to be created in branch - PATH-LOCAL:PATH-TARGET")
	createCmd.Flags().StringSlice("labels", []string{}, "labels to be added to pull request")
}
