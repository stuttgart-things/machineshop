/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/internal"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	tmpDownloadDir = "/tmp/machineShop"
)

func InstallBin(selectedInstallProfiles []string, allConfig Profile, bin string) {

	binDir := sthingsCli.AskSingleInputQuestion("BIN DIR:", bin)

	if !internal.CheckForUnixWritePermissions(binDir) {
		fmt.Println("NO WRITE PERMISSIONS!", binDir)
	} else {

		for _, binaryProfile := range allConfig.BinaryProfile {

			for _, selectedProfile := range selectedInstallProfiles {

				if binaryProfile[selectedProfile].Url != "" {

					fmt.Println(binaryProfile[selectedProfile].Url)
					sthingsCli.DownloadFileWithProgressBar(binaryProfile[selectedProfile].Url, tmpDownloadDir)

				}

			}

		}

	}

}
