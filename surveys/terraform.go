/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/internal"
)

func RunTerraform(gitRepository, gitPath string) {

	fmt.Println(gitRepository)
	fmt.Println(gitPath)

	internal.GetFileListFromGitRepository(gitRepository, gitPath, nil)
}
