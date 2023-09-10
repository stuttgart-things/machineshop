/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

var (
	logger = sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)
)

func PrintBanner(logFilePath, gitPath, gitRepository, version, date, cmd string) {

	ptermLogo, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("machine", pterm.NewStyle(pterm.FgLightCyan)),
		putils.LettersFromStringWithStyle("Shop", pterm.NewStyle(pterm.FgLightMagenta))).
		Srender()

	pterm.DefaultCenter.Print("\n" + ptermLogo)
	pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightCyan)).WithMargin(2).Sprint(cmd))
	pterm.Info.Println(pterm.White("GIT-REPO ") + "\t\t" + pterm.LightMagenta(gitRepository) + "\n" +
		pterm.White("GIT-PATH ") + "\t\t" + pterm.LightMagenta(gitPath) + "\n" +
		pterm.White("VAULT_ADDR ") + "\t\t" + pterm.LightMagenta(os.Getenv("VAULT_ADDR")) + "\n" +
		pterm.White("VAULT_NAMESPACE ") + "\t\t" + pterm.LightMagenta(os.Getenv("VAULT_NAMESPACE")) + "\n" +
		pterm.White("VAULT_ROLE_ID ") + "\t\t" + pterm.LightMagenta(os.Getenv("VAULT_ROLE_IDta")) + "\n" +
		pterm.White("VAULT_SECRET_ID ") + "\t\t" + pterm.LightMagenta(os.Getenv("VAULT_SECRET_ID")) + "\n" +
		pterm.White("VAULT_TOKEN ") + "\t\t" + pterm.LightMagenta(os.Getenv("VAULT_TOKEN")) + "\n" +
		pterm.White("LOG-FILE ") + "\t\t" + pterm.LightMagenta(logFilePath) + "\n" +
		"\n" +
		pterm.White("VERSION ") + "\t\t\t" + pterm.LightMagenta(version+" ("+date+")"))
	pterm.Println()
}

func HandleRenderOutput(outputFormat, destinationPath, renderedTemplate string, overwrite bool) {

	switch outputFormat {
	default:
		logger.Error(outputFormat, "output format not defined")
	case "stdout":
		fmt.Println(string(renderedTemplate))
	case "file":
		logger.Info("output file written to ", destinationPath)
		sthingsBase.WriteDataToFile(destinationPath, string(renderedTemplate))
	}

}
