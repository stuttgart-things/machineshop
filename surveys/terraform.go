/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/functions"
)

func RunTerraform(gitRepository, gitPath string) {

	fmt.Println(gitRepository)
	fmt.Println(gitPath)

	functions.GetFileListFromGitRepository(gitRepository, gitPath, nil)
}
