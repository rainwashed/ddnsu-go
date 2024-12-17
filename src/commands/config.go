package commands

import (
	"ddnsu/v2/src/global"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func configRun(cmd *cobra.Command, args []string) {
	cmd.Parent().Help()
}

func viewRun(cmd *cobra.Command, args []string) {
	cmd.Parent().Help()
}

func backupRun(cmd *cobra.Command, args []string) {

	configPath := global.ConfigurationPath
	newFilePath := ""

	if len(args) == 0 {
		newFilePath = global.ConfigurationPath + ".bak"
	} else {
		newFilePath = args[0]
	}

	fmt.Printf("newFilePath:%v", newFilePath)

	data, readErr := os.ReadFile(configPath)

	if readErr != nil {
		fmt.Println("could not read configuration file.")
		return
	}

	writeErr := os.WriteFile(newFilePath, data, 0644)

	if writeErr != nil {
		fmt.Println("could not write new configuration file.")
	}
}

var ConfigCommand = &cobra.Command{
	Use:   "config <command>",
	Short: "Commands to work with the config file.",
	Args:  cobra.MinimumNArgs(1),
	Run:   configRun,
}

var ViewCommand = &cobra.Command{
	Use:   "view",
	Short: "View the configuration file.",
	Args:  cobra.ExactArgs(0),
	Run:   viewRun,
}
var BackupCommand = &cobra.Command{
	Use:   "backup [output file]",
	Short: "Backup the current configuration file.",
	Args:  cobra.MaximumNArgs(1),
	Run:   backupRun,
}
var ResetCommand = &cobra.Command{
	Use:   "reset",
	Short: "This will reset the current configuration file to the example one on GitHub.",
}
