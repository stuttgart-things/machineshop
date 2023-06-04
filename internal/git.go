/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"

	memfs "github.com/go-git/go-billy/v5/memfs"
	memory "github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5"
)

func GetFileListFromGitRepository(repository, directory string, auth *http.BasicAuth) (fileList, directoryList []string) {

	// Init memory storage and fs
	storer := memory.NewStorage()
	fs := memfs.New()

	// Clone repo into memfs
	_, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  repository,
		Auth: auth,
	})

	if err != nil {
		fmt.Println("Could not git clone repository")
	}

	files, _ := fs.ReadDir(directory)

	for _, file := range files {

		if file.IsDir() {
			directoryList = append(directoryList, file.Name())
		} else {
			fileList = append(fileList, file.Name())
		}
	}

	return
}

func GitCommitFile(repository string, auth *http.BasicAuth, fileContent []byte, filePath, commitMsg string) error {

	// filePath2 := "gitops/stage/labda-vsphere"
	// Init memory storage and fs
	storer := memory.NewStorage()
	fs := memfs.New()

	// Clone repo into memfs
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  repository,
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("Could not git clone repository %s: %w", repository, err)
	}
	fmt.Println("Repository cloned")

	// Get git default worktree
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("Could not get git worktree: %w", err)
	}

	fmt.Println(w)

	// Create new file
	newFile, err := fs.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create new file: %w", err)
	}
	newFile.Write(fileContent)
	newFile.Close()

	// DELETE file

	// err = fs.Remove(filePath2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Run git status before adding the file to the worktree
	fmt.Println(w.Status())

	// git add $filePath
	w.Add(filePath)
	// w.Remove(filePath2)
	// Run git status after the file has been added adding to the worktree
	fmt.Println("STATUUUS")
	fmt.Println(w.Status())

	// git commit -m $message
	w.Commit(commitMsg, &git.CommitOptions{})

	//Push the code to the remote
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("Could not git push: %w", err)
	}
	fmt.Println("Remote updated.", filePath)

	return nil
}

func GetGitAuth(gitUser, gitToken string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: gitUser,
		Password: gitToken,
	}
}
