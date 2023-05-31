/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/internal"
)

func RunTerraform(gitRepository, gitPath string) {

	fileList := internal.GetFileListFromGitRepository(gitRepository, gitPath, nil)
	fmt.Println(fileList)
}
