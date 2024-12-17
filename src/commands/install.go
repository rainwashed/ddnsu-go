package commands

import (
	"ddnsu/v2/src/global"
	"ddnsu/v2/src/utils"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func InstallCommand(cmd *cobra.Command, args []string) {
	prompt := promptui.Select{
		Label: "Are you sure?",
		Items: []string{"Yes", "No"},
	}

	fmt.Println("Installing ddnsu as a systemd service will require that a binary be downloaded from the GitHub releases and placed in the ~/.config/ddnsu directory. A systemd service will be created in the ~/.config/systemd directory.")

	_, result, _ := prompt.Run()
	if result == "No" {
		fmt.Println(color.New(color.Italic).Sprint("Nothing was installed."))
		os.Exit(0)
	}

	// edit the systemd file
	user, err := user.Current()

	if err != nil {
		fmt.Println(color.RedString("could not detect a user. error: %v", err))
		os.Exit(1)
	}

	serviceScript := string(global.ExampleServiceEmbed[:])
	serviceScript = strings.Replace(serviceScript, "@user", user.Username, -1)

	utils.DetermineFilesExistenceAndCreateIfDoesNotExist(filepath.Join(user.HomeDir, ".config/systemd/user"), "ddnsu.service", []byte(serviceScript))

	// create the active script
	utils.DetermineFilesExistenceAndCreateIfDoesNotExist(filepath.Join(user.HomeDir, ".config/ddnsu"), "run.sh", global.ExampleShellEmbed)
	os.Chmod(filepath.Join(user.HomeDir, ".config/ddnsu/run.sh"), 0777)

	fmt.Println(color.New(color.Italic).Sprintf("Start the service with: systemctl --user start ddnsu.service"))
}
