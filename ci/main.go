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

	surveys "github.com/stuttgart-things/machineshop/surveys"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"fmt"
)

// Parse the YAML content
var config struct {
	Binary []struct {
		Name    string `yaml:"name"`
		URL     string `yaml:"url"`
		Bin     string `yaml:"bin"`
		Version string `yaml:"version"`
	} `yaml:"binary"`
}

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
) (*dagger.File, error) {
	// Initialize the Go module
	goModule := dag.Go()

	// Call the Build function with the struct
	buildOutput := goModule.Binary(src, dagger.GoBinaryOpts{
		GoVersion:  "1.24.0",      // Go version
		Os:         "linux",       // OS
		Arch:       "amd64",       // Architecture
		GoMainFile: "main.go",     // Main Go file
		BinName:    "machineshop", // Binary name
	})

	// Extract the binary file from the build output directory
	binaryFile := buildOutput.File("machineshop")

	// Read the YAML file
	yamlFile := src.File("profiles/binaries.yaml")
	yamlContent, err := yamlFile.Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Parse the YAML content into a struct
	var allConfig surveys.Profile
	allConfig = sthingsCli.ReadInlineYamlToObject([]byte(yamlContent), allConfig).(surveys.Profile)

	// Extract all binaries from the YAML
	allBinaries := []string{}
	for _, binaryProfile := range allConfig.BinaryProfile {
		for key := range binaryProfile {
			allBinaries = append(allBinaries, key)
		}
	}

	// Create a new Machineshop container
	machineShopContainer := dag.Container().From("eu.gcr.io/stuttgart-things/sthings-workflow:1.30.1")

	// Copy the binary into the container at /usr/bin/
	machineShop := machineShopContainer.
		WithFile("/usr/bin/machineshop", binaryFile).
		WithExec([]string{"chmod", "+x", "/usr/bin/machineshop"})

	// Define the merged log file path
	mergedLogPath := "/var/log/machineshop_tests.log"

	// Test the binary by running it inside the container and append output to the merged log file
	machineShop = machineShop.
		WithExec([]string{"sh", "-c", "machineshop version >> " + mergedLogPath + " 2>&1"})

	// Function to handle binary installation with error handling
	installBinary := func(container *dagger.Container, binary string) *dagger.Container {
		_, err := container.
			WithExec([]string{"sh", "-c", "machineshop install --binaries " + binary + " >> " + mergedLogPath + " 2>&1"}).
			Sync(ctx)
		if err != nil {
			fmt.Printf("Failed to install binary %s: %v\n", binary, err)
			// Log the failure to the merged log file
			container = container.
				WithExec([]string{"sh", "-c", "echo 'Failed to install binary " + binary + "' >> " + mergedLogPath + " 2>&1"})
		} else {
			container = container.
				WithExec([]string{"sh", "-c", "echo 'Successfully installed binary " + binary + "' >> " + mergedLogPath + " 2>&1"})
		}
		return container
	}

	// Install each binary, skipping failures
	for _, binary := range allBinaries {
		machineShop = installBinary(machineShop, binary)
	}

	// Get the merged log file as a dagger.File object
	mergedLogFile := machineShop.File(mergedLogPath)

	// Return the container and the merged log file
	return mergedLogFile, nil
}
