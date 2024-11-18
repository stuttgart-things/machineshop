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
