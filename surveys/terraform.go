/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pterm/pterm"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/stuttgart-things/machineShop/internal"
)

var (

	// LOGGING
	fileWriter  = internal.CreateFileLogger("/tmp/machineShop.log")
	multiWriter = io.MultiWriter(os.Stdout, fileWriter)
	logger      = pterm.DefaultLogger.
			WithLevel(pterm.LogLevelTrace).
			WithWriter(multiWriter). // Only show logs with a level of Trace or higher.
			WithCaller()             // ! Show the caller of the log function.

)

func RunTerraform(gitRepository, gitPath string) {
	fileList := internal.GetFileListFromGitRepository(gitRepository, gitPath, nil)
	fmt.Println(fileList)
	logger.Info("Found files in ", logger.Args("", strings.Join(fileList, " ")))

	sthingsCli.AskSingleSelectQuestion("Select profile:", fileList)

}
