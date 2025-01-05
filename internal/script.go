/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"io"
	"os"
	"os/exec"
)

func ExecuteShellScript(scriptName, script, logDir string) {

	// OPEN THE LOG FILE FOR WRITING
	logFile, err := os.OpenFile(logDir+"/"+scriptName+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("FAILED TO OPEN LOG FILE: %v", err)
	}
	defer logFile.Close()

	// CREATE A MULTIWRITER TO WRITE TO BOTH STDOUT AND THE LOG FILE
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// CREATE THE COMMAND TO EXECUTE THE SCRIPT
	cmd := exec.Command("sh", "-c", script) // Use "sh -c" for shell commands

	// SET THE COMMAND'S STDOUT AND STDERR TO THE MULTIWRITER
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	// RUN THE COMMAND AND CHECK FOR ERRORS
	if err := cmd.Run(); err != nil {
		log.Fatalf("SCRIPT EXECUTION FAILED: %v", err)
	}

	log.Println("SCRIPT RUN SUCCESSFUL")
}
