/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"
	"strings"
	"time"

	"github.com/stuttgart-things/machineshop/internal"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	currentTime = time.Now()
	timestamp   = currentTime.Format("2006-01-02_15-04-05")
)

func RenderInstallScriptNew(runConfig internal.RunScripts, scripts []string, variables map[string]interface{}) {

	var allScripts []string

	// DEBUGGING OUTPUT
	// fmt.Printf("PARSED GLOBALS: %+v\n", runConfig.Globals)
	// for i, entry := range runConfig.Run {
	// 	fmt.Printf("RUN %d: %+v\n", i+1, entry)
	// }

	// fmt.Println("Scripts to be run:", scripts)

	// RUN SURVEY (IF SCRIPTS NOT DEFINED)
	// internal.CreateRunSurvey(runConfig)

	// LOAD GLOBAL VARS
	allVars = runConfig.Globals
	log.Info("GLOBAL VARS: ", allVars)

	// VARIABLES FROM PARAMETER
	log.Info("FLAG VARIABLES: ", variables)

	// SELECT TO BE RUN SCRIPTS
	for _, entry := range runConfig.Run {

		if sthingsBase.CheckForStringInSlice(scripts, entry.Name) {
			renderVars := make(map[string]interface{})
			renderVars = sthingsBase.MergeMaps(allVars, entry.Vars)
			log.Info("COMMAND VARS MERGED w/ GLOBALS: ", renderVars)

			renderVars = sthingsBase.MergeMaps(renderVars, variables)
			log.Info("VARS MERGED w/ FLAG VARS (USED FOR RENDERING): ", renderVars)

			// RENDER SCRIPT
			toBeRendered := bashScriptHeader + entry.Script
			renderedScript, error := sthingsBase.RenderTemplateInline(toBeRendered, "missingkey=error", "{{", "}}", renderVars)
			if error != nil {
				log.Error("ERROR WHILE RENDERING", error)
			}
			fmt.Println("RENDERED SCRIPT", string(renderedScript))

			// WRITE (RENDRED) SCRIPT TO TMP
			sthingsBase.WriteDataToFile("/tmp/"+entry.Name+".sh", string(renderedScript))
			log.Info("SCRIPTS WRITTEN TO TMP: ", "/tmp/"+entry.Name+".sh")

			// ADD TO LIST OF SCRIPTS TO BE RUN (ONE BIG STRING)
			allScripts = append(allScripts, string(renderedScript))

		}

		// LOAD VARS FROM SCRIPT
		// renderVars = sthingsBase.MergeMaps(allVars, entry.Vars)

		// RENDER SCRIPT

		// WRITE (RENDRED) SCRIPT TO TMP

		// ADD TO LIST OF SCRIPTS TO BE RUN (ONE BIG STRING)

		// fmt.Printf("Step %d:\n", i+1)
		// fmt.Printf("  Description: %s\n", entry.Description)
		// fmt.Printf("  Script: %s\n", entry.Script)
		// if entry.Vars != nil {
		// 	fmt.Println("  Vars:")
		// 	for key, value := range entry.Vars {
		// 		fmt.Printf("    %s: %v\n", key, value)
		// 	}
		// }
	}

	// WRITE ALL SCRIPTS TO TMP
	mergedScripts := string(strings.Join(allScripts, "\n"))
	filename := fmt.Sprintf("script_%s", timestamp)
	sthingsBase.WriteDataToFile("/tmp/"+filename+".sh", mergedScripts)
	log.Info("ALL SCRIPTS WRITTEN TO TMP: ", "/tmp/"+filename+".sh")

	// EXECUTE ALL SCRIPTS - LOOP OVER ALL SCRIPTS
	internal.ExecuteShellScript("all", mergedScripts, "/tmp")
}
