/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"time"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/machineShop/surveys"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "manage infrastructure in any cloud",
	Long:  `predictably provision and manage infrastructure in any cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		gitPath, _ := cmd.LocalFlags().GetString("path")

		panels := pterm.Panels{
			{{Data: pterm.White("\n/terraform")}, {Data: pterm.White("\n" + version)}, {Data: pterm.White("\n" + date)}},
			{{Data: pterm.Magenta("\nGIT-REPO: " + gitRepository)}, {Data: pterm.Magenta("\nGIT-PATH: " + gitPath)}},
			{{Data: pterm.Magenta("\nVAULT:\n" + gitRepository)}, {Data: pterm.Magenta("\nGIT-PATH:\n" + gitPath)}},
		}

		// Print panels.
		_ = pterm.DefaultPanel.WithPanels(panels).WithPadding(5).Render()

		ptermLogo, _ := pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle("machine", pterm.NewStyle(pterm.FgLightCyan)),
			putils.LettersFromStringWithStyle("Shop", pterm.NewStyle(pterm.FgLightMagenta))).
			Srender()

		pterm.DefaultCenter.Print(ptermLogo)

		pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(2).Sprint("/TERRAFORM"))

		pterm.Info.Println("GIT-REPO: " + gitRepository +
			"\nPTerm works on nearly every terminal and operating system." +
			"\nIt's super easy to use!" +
			"\nIf you want, you can customize everything :)" +
			"\nYou can see the code of this demo in the " + pterm.LightMagenta("./_examples/demo") + " directory." +
			"\n" +
			"\nThis demo was updated at: " + pterm.Green(time.Now().Format("02 Jan 2006 - 15:04:05 MST")))
		pterm.Println()

		surveys.RunTerraform(gitRepository, gitPath)
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)

	terraformCmd.Flags().String("path", "machineShop/terraform", "path to terraform automation code")

}
