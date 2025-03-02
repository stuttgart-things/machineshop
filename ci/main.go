/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/ci/internal/dagger"
	"fmt"
)

type Ci struct {
	MachineShopContainer *dagger.Container
}

// GetMachineShopContainer return the default image for golang
func (m *Ci) GetMachineShopContainer() *dagger.Container {
	return dag.Container().
		From("ghcr.io/stuttgart-things/github.com/stuttgart-things/machineshop:v2.6.3")
}

func New(
	// machineShop container
	// It need contain machineShop
	// +optional
	machineShopContainer *dagger.Container,

) *Ci {
	machineShop := &Ci{}

	if machineShopContainer != nil {
		machineShop.MachineShopContainer = machineShopContainer
	} else {
		machineShop.MachineShopContainer = machineShop.GetMachineShopContainer()
	}

	return machineShop
}

// Test machineShop version command
func (m *Ci) TestVersion(ctx context.Context) (versionOutput string) {

	fmt.Println("RUNNING VERSION COMMAND...")

	versionOutput, err := m.MachineShopContainer.
		WithExec(
			[]string{"machineshop", "version"}).
		Stdout(ctx)

	if err != nil {
		fmt.Println("ERROR RUNNING VERSION COMMAND: ", err)
	}

	fmt.Println(versionOutput)

	return versionOutput
}

// Test machineShop version command
func (m *Ci) TestInstall(ctx context.Context) (versionOutput string) {

	fmt.Println("RUNNING INSTALL COMMAND...")

	installCmdOutput, err := m.MachineShopContainer.
		WithExec(
			[]string{"machineshop", "install", "--profile", "machineShop/binaries.yaml", "--binaries", "sops,kubectl,flux"}).
		Stdout(ctx)

	if err != nil {
		fmt.Println("ERROR RUNNING INSTALL COMMAND: ", err)
	}

	fmt.Println(installCmdOutput)

	return installCmdOutput
}

func (m *Ci) BuildAndUse(
	ctx context.Context,
	src *dagger.Directory,
) (*dagger.Container, error) {
	// Initialize the Go module
	goModule := dag.Go()

	// Call the Build function with the struct
	buildOutput := goModule.Binary(src, dagger.GoBinaryOpts{
		GoVersion:  "1.24.0",      // Go verssion
		Os:         "linux",       // OS
		Arch:       "amd64",       // Architecture
		GoMainFile: "main.go",     // Main Go file
		BinName:    "machineshop", // Binary name
	})

	// Extract the binary file from the build output directory
	binaryFile := buildOutput.File("machineshop")

	// Create a new Machineshop container
	machineShopContainer := dag.Container().From("eu.gcr.io/stuttgart-things/sthings-workflow:1.30.1")

	// Copy the binary into the container at /usr/bin/
	machineShop := machineShopContainer.
		WithFile("/usr/bin/machineshop", binaryFile).
		WithExec([]string{"chmod", "+x", "/usr/bin/machineshop"})
	// Debug: List files in /usr/bin/ to verify the binary is copied
	// debugOutput, err := machineShop.WithExec([]string{"ls", "-l", "/usr/bin/"}).Stdout(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to list /usr/bin/: %w", err)
	// }
	// fmt.Println("Contents of /usr/bin/:", debugOutput)

	// Optionally, test the binary by running it inside the container
	testVersion, err := machineShop.WithExec([]string{"machineshop", "version"}).Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to test version: %w", err)
	}

	testInstall, err := machineShop.WithExec([]string{"machineshop", "install", "--profile", "machineShop/binaries.yaml", "--binaries", "sops,kubectl,flux"}).Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to test install: %w", err)
	}

	// Print the test output (optional)
	fmt.Println("Binary test output:", testVersion)
	fmt.Println("Binary test output:", testInstall)

	// Return the container with the binary for further use
	return machineShop, nil
}
