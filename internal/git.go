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
