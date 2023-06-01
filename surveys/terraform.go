/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/stuttgart-things/machineShop/internal"
)

var (
	selectedProfile string
	configSubFolder = "state"

	// LOGGING
	// fileWriter  = internal.CreateFileLogger("/tmp/machineShop.log")
	// multiWriter = io.MultiWriter(os.Stdout, fileWriter)
	// logger      = pterm.DefaultLogger.
	// 		WithLevel(pterm.LogLevelTrace).
	// 		WithWriter(multiWriter). // Only show logs with a level of Trace or higher.
	// 		WithCaller()             // ! Show the caller of the log function.

)

func RunTerraform(gitRepository, gitPath string) {

	if sthingsCli.AskSingleSelectQuestion("OPERATION:", []string{"apply", "destroy"}) == "apply" {

		fileList, _ := internal.GetFileListFromGitRepository(gitRepository, gitPath, nil)
		selectedProfile = sthingsCli.AskSingleSelectQuestion("SELECT PROFILE:", fileList)
		fmt.Println(selectedProfile)

	} else {
		_, folderList := internal.GetFileListFromGitRepository(gitRepository, gitPath+"/"+configSubFolder, nil)
		selectedProfile = sthingsCli.AskSingleSelectQuestion("SELECT PROFILE:", folderList)
		fmt.Println(selectedProfile)

	}
}
