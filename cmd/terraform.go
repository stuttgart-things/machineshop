/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"io"
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/machineShop/internal"
	"github.com/stuttgart-things/machineShop/surveys"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "manage infrastructure in any cloud",
	Long:  `predictably provision and manage infrastructure in any cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		gitPath, _ := cmd.LocalFlags().GetString("path")

		// panels := pterm.Panels{
		// 	{{Data: pterm.White("\n/terraform")}, {Data: pterm.White("\n" + version)}, {Data: pterm.White("\n" + date)}},
		// 	{{Data: pterm.Magenta("\nGIT-REPO: " + gitRepository)}, {Data: pterm.Magenta("\nGIT-PATH: " + gitPath)}},
		// 	{{Data: pterm.Magenta("\nVAULT:\n" + gitRepository)}, {Data: pterm.Magenta("\nGIT-PATH:\n" + gitPath)}},
		// }

		// Print panels.
		// _ = pterm.DefaultPanel.WithPanels(panels).WithPadding(5).Render()

		ptermLogo, _ := pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle("machine", pterm.NewStyle(pterm.FgLightCyan)),
			putils.LettersFromStringWithStyle("Shop", pterm.NewStyle(pterm.FgLightMagenta))).
			Srender()

		pterm.DefaultCenter.Print("\n" + ptermLogo)

		pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightCyan)).WithMargin(2).Sprint("/TERRAFORM"))

		pterm.Info.Println(pterm.White("GIT-REPO: ") + "\t" + pterm.LightMagenta(gitRepository) + "\n" +
			pterm.White("GIT-PATH: ") + "\t" + pterm.LightMagenta(gitPath) + "\n" +
			pterm.White("VAULT_ADDR: ") + "\t" + pterm.LightMagenta(os.Getenv("VAULT_ADDR")) + "\n" +
			pterm.White("VAULT_NAMESPACE: ") + pterm.LightMagenta(os.Getenv("VAULT_NAMESPACE")) + "\n" +
			pterm.White("VAULT_ROLE_ID: ") + "\t" + pterm.LightMagenta(os.Getenv("VAULT_ROLE_IDta")) + "\n" +
			pterm.White("VAULT_SECRET_ID: ") + pterm.LightMagenta(os.Getenv("VAULT_SECRET_ID")) + "\n" +
			pterm.White("VAULT_TOKEN: ") + "\t" + pterm.LightMagenta(os.Getenv("VAULT_TOKEN")) + "\n" +
			pterm.White("LOG-FILE: ") + "\t" + pterm.LightMagenta(logFilePath) + "\n" +
			"\n" +
			pterm.White("VERSION: ") + "\t" + pterm.LightMagenta(version+" ("+date+")"))
		pterm.Println()

		surveys.RunTerraform(gitRepository, gitPath)

		fileWriter := internal.CreateFileLogger(logFilePath)

		multiWriter := io.MultiWriter(os.Stdout, fileWriter)

		logger := pterm.DefaultLogger.
			WithLevel(pterm.LogLevelTrace).
			WithWriter(multiWriter). // Only show logs with a level of Trace or higher.
			WithCaller()             // ! Show the caller of the log function.

		logger.Trace("Doing not so important stuff", logger.Args("priority", "super low"))
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)

	terraformCmd.Flags().String("path", "machineShop/terraform", "path to terraform automation code")

}
