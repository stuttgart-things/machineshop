/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"
	"sync"

	"github.com/stuttgart-things/machineShop/internal"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	tmpDownloadDir = "/tmp/machineShop"
	wg             sync.WaitGroup
)

func InstallBin(selectedInstallProfiles []string, allConfig Profile, bin string) {

	binDir := sthingsCli.AskSingleInputQuestion("BIN DIR:", bin)

	if !internal.CheckForUnixWritePermissions(binDir) {
		fmt.Println("NO WRITE PERMISSIONS!", binDir)
	} else {

		for _, binaryProfile := range allConfig.BinaryProfile {

			for _, selectedProfile := range selectedInstallProfiles {

				if binaryProfile[selectedProfile].Url != "" {

					wg.Add(1)

					name := selectedProfile
					url := binaryProfile[selectedProfile].Url

					go func() {
						defer wg.Done()

						fmt.Println("Downloading", name, url)
						sthingsCli.DownloadFileWithProgressBar(url, tmpDownloadDir)

					}()

				}

			}

		}

		wg.Wait()

	}

}