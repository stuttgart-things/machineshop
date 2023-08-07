/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
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
