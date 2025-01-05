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
	log              = sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)
	wg               sync.WaitGroup
	bashScriptHeader = "#!/bin/bash\n"
	allVars          map[string]interface{}
)

func RenderInstallScript(selectedScriptProfiles []string, allConfig Profile, scriptDir string) {

	// CREATE TMP DL FOLDER
	sthingsBase.CreateNestedDirectoryStructure(scriptDir, 0600)
	log.Info("CREATED TMP SCRIPT DIR ", scriptDir)

	for _, scriptProfile := range allConfig.ScriptProfile {

		for _, selectedProfile := range selectedScriptProfiles {

			// GET ALL VARS FROM TEMPLATE
			scriptTemplate := scriptProfile[selectedProfile].Script
			allTemplateVariablesAndDefaults, _, _, _ := sthingsBase.GetVariablesAndDefaultsFromTemplate(scriptTemplate, "curly")
			log.Info("ALL VARS ", allTemplateVariablesAndDefaults)

			// RENDER URL IF TEMPLATE ULR CONTAINS AT LEAST ONE VARIABLE
			script := bashScriptHeader + scriptTemplate
			if len(allTemplateVariablesAndDefaults) >= 1 {
				renderedScript, error := sthingsBase.RenderTemplateInline(string(scriptTemplate), "missingkey=error", "{{", "}}", allTemplateVariablesAndDefaults)
				if error != nil {
					log.Error("ERROR WHILE RENDERING", error)
				}
				script = bashScriptHeader + string(renderedScript)
			}

			// OUTPUT HANDLING
			log.Info("RENDERED SCRIPT ", script)
			scriptOutputPath := scriptDir + "/" + selectedProfile + ".sh"
			scriptAlreadyExists, _ := sthingsBase.VerifyIfFileOrDirExists(scriptOutputPath, "file")

			if scriptAlreadyExists {
				sthingsBase.DeleteFile(scriptOutputPath)
				log.Warn("EXISTING SCRIPT DELETED ", scriptOutputPath)
			}

			if sthingsBase.WriteDataToFile(scriptOutputPath, script) {
				log.Info("SCRIPT WRITTEN TO ", scriptDir+"/"+selectedProfile)
			} else {
				log.Info("WRITING SCRIPT TO DIR FAILED ", scriptDir+"/"+selectedProfile)
			}

			// EXECUTING SghCRIPT
			log.Info("EXECUTING SCRIPT ", scriptDir+"/"+selectedProfile)
			fmt.Println(sthingsBase.ExecuteBashScript(scriptOutputPath))

		}

	}

}

func InstallBinaries(selectedInstallProfiles []string, allConfig Profile, tmpDownloadDir, binDir string) {

	if !sthingsBase.CheckForUnixWritePermissions(binDir) {
		log.Error("NO WRITE PERMISSIONS! ", binDir)

	} else {
		// CREATE TMP DL FOLDER
		sthingsBase.CreateNestedDirectoryStructure(tmpDownloadDir, 0600)
		log.Info("CREATED TMP DOWNLOAD DIR ", tmpDownloadDir)

		// DOWNLOAD + INSTALL BINARIES
		for _, binaryProfile := range allConfig.BinaryProfile {

			for _, selectedProfile := range selectedInstallProfiles {

				if binaryProfile[selectedProfile].Url != "" {

					wg.Add(1)

					// SET VARS FOR CONCURRENCY
					name := selectedProfile
					url := binaryProfile[selectedProfile].Url
					binName := binaryProfile[selectedProfile].Bin
					version := binaryProfile[selectedProfile].Version

					// RENDER URL IF TEMPLATE ULR CONTAINS AT LEAST ONE VARIABLE
					allTemplateVariablesAndDefaults, _, _, _ := sthingsBase.GetVariablesAndDefaultsFromTemplate(url, "curly")
					allTemplateVariablesAndDefaults["version"] = version

					fmt.Println("ALL VARS", allTemplateVariablesAndDefaults)

					if len(allTemplateVariablesAndDefaults) >= 1 {
						renderedURL, _ := sthingsBase.RenderTemplateInline(string(url), "missingkey=zero", "{{", "}}", allTemplateVariablesAndDefaults)
						url = string(renderedURL)
					}
					// ADD OVERWIRTE OPTION
					// CHANGE VERSION FOR USING v PREFIX

					go func() {
						defer wg.Done()

						var tmpBinPath string

						// CHECK IF URL IS REACHABLE/VALID
						if sthingsCli.CheckUrlAvailability(url) {

							// DOWNLOAD ARCHIVE
							log.Info("DOWNLOADING! ", url)

							sthingsCli.DownloadFileWithProgressBar(url, tmpDownloadDir)

							// EXTRACT (IF BINARY IS ARCHIVED)
							if strings.Contains(url, ".zip") {
								err := sthingsCli.UnZipArchive(tmpDownloadDir+"/"+filepath.Base(url), tmpDownloadDir+"/"+name)

								if err != nil {
									tmpBinPath = tmpDownloadDir + "/" + name + "/" + binName
								}

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
						} else {
							log.Error("URL NOT REACHABLE! ", url)
						}

					}()

				}

			}

		}

		wg.Wait()

	}

}
