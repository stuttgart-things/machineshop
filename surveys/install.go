/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	tmpDownloadDir = "/tmp/machineShop"
	wg             sync.WaitGroup
)

func InstallBinaries(selectedInstallProfiles []string, allConfig Profile, bin string) {

	binDir := sthingsCli.AskSingleInputQuestion("BIN DIR:", bin)

	if !sthingsBase.CheckForUnixWritePermissions(binDir) {
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

						// DOWNLOAD ARCHIVE
						fmt.Println("DOWNLOADING..", name, url)
						sthingsCli.DownloadFileWithProgressBar(url, tmpDownloadDir)

						// EXTRACT (IF BINARY IS ARCHIVED)
						if strings.Contains(url, ".zip") {
							sthingsCli.UnZipArchive(tmpDownloadDir+"/"+filepath.Base(url), tmpDownloadDir+"/"+name)
						} else if strings.Contains(url, ".tar.gz") {
							sthingsCli.ExtractTarGzArchive(tmpDownloadDir+"/"+filepath.Base(url), tmpDownloadDir+"/"+name, 0700)
						}

						// MOVE BINARY
						sthingsBase.MoveRenameFileOnFS(tmpDownloadDir+"/"+name+"/"+name, binDir+"/"+name)
						fmt.Println("MOVING.." + tmpDownloadDir + "/" + name + "/" + name + " TO " + binDir + "/" + name)

						// CHANGE BINARY PERMISSION TO EXECUTE
						sthingsBase.SetUnixFilePermissions(binDir+"/"+name, 0755)

						// DELETE ARCHIVE/EXTRACTFOLDER
						sthingsBase.RemoveNestedFolder(tmpDownloadDir + "/" + name)

					}()

				}

			}

		}

		wg.Wait()

	}

}
