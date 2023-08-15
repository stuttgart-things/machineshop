/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"
	// sthingsCli "github.com/stuttgart-things/sthingsCli"
)

func InstallBin(selectedInstallProfiles []string, allConfig Profile) {

	// downloadTarget := "/tmp"

	for _, binaryProfile := range allConfig.BinaryProfile {

		// for key := range binaryProfile {

		for _, selectedProfile := range selectedInstallProfiles {
			fmt.Println(binaryProfile[selectedProfile].Url)
		}

		// }

	}

	// sthingsCli.DownloadFileWithProgressBar(downloadURL, downloadTarget)

}
