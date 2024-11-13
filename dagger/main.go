/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package main

import (
	"context"
	"dagger/machineshop/internal/dagger"
	"fmt"
)

type Machineshop struct {
	MachineShopContainer *dagger.Container
}

// GetMachineShopContainer return the default image for golang
func (m *Machineshop) GetMachineShopContainer() *dagger.Container {
	return dag.Container().
		From("ghcr.io/stuttgart-things/github.com/stuttgart-things/machineshop:v2.6.3")
}

func New(
	// machineShop container
	// It need contain machineShop
	// +optional
	machineShopContainer *dagger.Container,

) *Machineshop {
	machineShop := &Machineshop{}

	if machineShopContainer != nil {
		machineShop.MachineShopContainer = machineShopContainer
	} else {
		machineShop.MachineShopContainer = machineShop.GetMachineShopContainer()
	}

	return machineShop
}

// Test machineShop version command
func (m *Machineshop) TestVersion(ctx context.Context) (versionOutput string) {

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
