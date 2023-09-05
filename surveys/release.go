/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

func CloneRepositories(selectedReleaseProfiles []string, allConfig ReleaseProfile, tmp string) {

	tmpDir := sthingsCli.AskSingleInputQuestion("TMP DIR:", tmp)

	if !sthingsBase.CheckForUnixWritePermissions(tmpDir) {
		fmt.Println("NO WRITE PERMISSIONS ON DIR!", tmpDir)

	} else {

		for _, repositoryProfile := range allConfig.RepositoryProfile {

			for _, selectedProfile := range selectedReleaseProfiles {

				if repositoryProfile[selectedProfile].Url != "" {

					wg.Add(1)

					name := selectedProfile
					url := repositoryProfile[selectedProfile].Url
					branch := repositoryProfile[selectedProfile].Branch
					version := repositoryProfile[selectedProfile].Version

					go func() {
						defer wg.Done()

						fmt.Println("Cloning", name, url)

						repo, _ := sthingsCli.CloneGitRepository(url, branch, version, nil)
						files, dirs := sthingsCli.GetFileListFromGitRepository(".", repo)
						fmt.Println(files, dirs)

						// hello, err := git.PlainClone("/tmp/release", true, &git.CloneOptions{
						// 	URL:               url,

						// 	Progress:          os.Stdout,
						// 	RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
						// })

					}()

				}

			}

		}

		wg.Wait()

	}

}
