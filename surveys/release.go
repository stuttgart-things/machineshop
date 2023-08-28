/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

func CloneRepositories(selectedReleaseProfiles []string, allConfig Profile, bin string) {

	binDir := sthingsCli.AskSingleInputQuestion("TMP DIR:", bin)

	if !internal.CheckForUnixWritePermissions(binDir) {
		fmt.Println("NO WRITE PERMISSIONS!", binDir)
	} else {

		for _, binaryProfile := range allConfig.BinaryProfile {

			for _, selectedProfile := range selectedReleaseProfiles {

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
