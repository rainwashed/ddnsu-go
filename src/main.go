package main

import (
	"ddnsu/v2/src/commands"
	"ddnsu/v2/src/global"
	"ddnsu/v2/src/services"
	"ddnsu/v2/src/utils"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "ddnsu",
	Short: "Dynamic domain name server updater is a service that dynamically updates nameserver records such as Cloudflare or Vercel.",
	Run:   runCommand,
}

func runCommand(cmd *cobra.Command, args []string) {
	cmd.Help()
	fmt.Println("Made with ❤️ by @rainwashed.")
}

func main() {
	if runtime.GOOS != "linux" {
		panic(fmt.Sprintf("%v is not a valid operating system for ddnsu. Expected: linux", runtime.GOOS))
	}

	homeDir, _ := os.UserHomeDir()
	configurationFilePath := filepath.Join(homeDir, ".config", "ddnsu", "config.toml")
	utils.DetermineFilesExistenceAndCreateIfDoesNotExist(filepath.Join(homeDir, ".config/ddnsu"), "config.toml", ExampleConfigEmbed)
	global.ConfigurationPath = configurationFilePath
	loaded, loadErr := utils.LoadConfigurationIntoGlobalVar(global.ConfigurationPath)
	utils.StoreActiveTokenInGlobalVar(global.Configuration)

	global.ExampleServiceEmbed = ExampleServiceEmbed
	global.ExampleShellEmbed = ExampleShellEmbed

	if loadErr != nil || !loaded {
		fmt.Println(color.RedString("Configuration file could not be loaded. There seems to be syntax issue. The location can be found at %v. The error is: %v", configurationFilePath, loadErr))
	}

	var versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Return the current version",
		Long:  "Return the current version of DDNSU command line utility and check if it is up-to-date.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _args []string) {
			utils.PrintFLn("DDNSU is currently on version %v", global.Configuration.Ddnsu.Version)
		},
	}
	rootCommand.AddCommand(versionCommand)

	var loginCommand = &cobra.Command{
		Use:   "login <vercel/cloudflare> <token>",
		Short: "Login to either Vercel or Cloudflare DNS server",
		Args:  cobra.ExactArgs(2),
		Run:   services.OnLoginCommandRun,
	}
	rootCommand.AddCommand(loginCommand)

	var recordCommand = &cobra.Command{
		Use:   "record <add/delete/update>",
		Short: "Add, delete, and modify the records that ddnsu will use.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			recordCommand := args[0]
			switch strings.ToLower(recordCommand) {
			case "add":
				commands.HandleCommand("add")
			case "delete":
				commands.HandleCommand("delete")
			case "update":
				commands.HandleCommand("update")
			default:
				fmt.Println("Not a valid command")
			}

		},
	}
	rootCommand.AddCommand(recordCommand)

	var cleanCommand = &cobra.Command{
		Use:   "clean <cloudflare/vercel>",
		Short: "Cleanup remote records to match configuration file.",
		Args:  cobra.ExactArgs(1),
		Run:   commands.CleanupCommand,
	}
	rootCommand.AddCommand(cleanCommand)

	var testCommand = &cobra.Command{
		Use:   "test",
		Short: "Emulates what the DNS records would look like based on the current configuration.",
		Args:  cobra.ExactArgs(0),
		Run:   commands.TestCommand,
	}
	rootCommand.AddCommand(testCommand)

	var startCommand = &cobra.Command{
		Use:   "start",
		Short: "Start the DDNSU service",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			services.BeginActiveLoop(time.Duration(global.Configuration.Ddnsu.Rate)*time.Minute, global.Configuration.Ddnsu.Use)
		},
	}
	rootCommand.AddCommand(startCommand)

	var installCommand = &cobra.Command{
		Use:   "install",
		Short: "Install ddnsu as a systemd service",
		Args:  cobra.ExactArgs(0),
		Run:   commands.InstallCommand,
	}
	rootCommand.AddCommand(installCommand)

	// config subcommands
	commands.ConfigCommand.AddCommand(commands.ViewCommand)
	commands.ConfigCommand.AddCommand(commands.BackupCommand)
	commands.ConfigCommand.AddCommand(commands.ResetCommand)
	rootCommand.AddCommand(commands.ConfigCommand)

	rootCommand.Execute()
}
