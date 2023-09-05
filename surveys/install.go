/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"path/filepath"
	"strings"
	"sync"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	log            = sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)
	tmpDownloadDir = "/tmp/machineShop"
	wg             sync.WaitGroup
)

func InstallBinaries(selectedInstallProfiles []string, allConfig Profile, bin string) {

	binDir := sthingsCli.AskSingleInputQuestion("BIN DIR:", bin)

	if !sthingsBase.CheckForUnixWritePermissions(binDir) {
		log.Error("NO WRITE PERMISSIONS! ", binDir)
	} else {

		for _, binaryProfile := range allConfig.BinaryProfile {

			for _, selectedProfile := range selectedInstallProfiles {

				if binaryProfile[selectedProfile].Url != "" {

					wg.Add(1)

					name := selectedProfile
					url := binaryProfile[selectedProfile].Url
					binName := binaryProfile[selectedProfile].Bin

					go func() {
						defer wg.Done()

						var tmpBinPath string

						// DOWNLOAD ARCHIVE
						log.Info("DOWNLOADING! ", url)
						sthingsCli.DownloadFileWithProgressBar(url, tmpDownloadDir)

						// EXTRACT (IF BINARY IS ARCHIVED)
						if strings.Contains(url, ".zip") {
							sthingsCli.UnZipArchive(tmpDownloadDir+"/"+filepath.Base(url), tmpDownloadDir+"/"+name)
							tmpBinPath = tmpDownloadDir + "/" + name + "/" + binName
						} else if strings.Contains(url, ".tar.gz") {
							sthingsCli.ExtractTarGzArchive(tmpDownloadDir+"/"+filepath.Base(url), tmpDownloadDir+"/"+name, 0700)
							tmpBinPath = tmpDownloadDir + "/" + name + "/" + binName
						} else {
							tmpBinPath = tmpDownloadDir + "/" + binName
						}

						destinationBinPath := binDir + "/" + name

						// VERIFY IF BIN EXISTS ALREADY
						binExists, _ := sthingsBase.VerifyIfFileOrDirExists(destinationBinPath, "file")

						if binExists { // ADD OVERWRITE OPTION
							sthingsBase.DeleteFile(destinationBinPath)
							log.Warn("EXISTING BINARY DELETED ", destinationBinPath)
						}

						// MOVE BINARY
						sthingsBase.MoveRenameFileOnFS(tmpBinPath, destinationBinPath)
						log.Info("MOVING ", tmpBinPath+" TO "+destinationBinPath)

						// CHANGE BINARY PERMISSION TO EXECUTE
						sthingsBase.SetUnixFilePermissions(destinationBinPath, 0755)
						log.Info("CHANGING PERMISSIONS TO EXECUTABLE OF ", destinationBinPath)

						// DELETE ARCHIVE/EXTRACTFOLDER
						sthingsBase.RemoveNestedFolder(tmpDownloadDir + "/" + name)
						log.Info("REMOVING ", tmpDownloadDir+"/"+name)

					}()

				}

			}

		}

		wg.Wait()

	}

}
