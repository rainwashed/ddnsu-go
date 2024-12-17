package commands

import (
	"ddnsu/v2/src/global"
	"ddnsu/v2/src/utils"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func boldRedText(text string) string {

	return color.New(color.FgRed, color.Bold).Sprint(text)
}

func TestCommand(cmd *cobra.Command, args []string) {
	ip := utils.MakeIpConsensus()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Type", "Subdomain", "Comment", "Value", "TTL"})
	errorPresent := false

	for _, entry := range global.Configuration.Ddnsu.Record {
		recordErr := recordTypeValidation(entry.Rtype)
		ttlErr := ttlValidation(strconv.Itoa(entry.Ttl))
		subdomainErr := subdomainValidation(entry.Subdomain)

		if recordErr != nil || ttlErr != nil || subdomainErr != nil {
			errorPresent = true
		}

		rtype := utils.IfThenElse(recordErr == nil, entry.Rtype, boldRedText(entry.Rtype)).(string)
		subdomain := utils.IfThenElse(subdomainErr == nil, entry.Subdomain, boldRedText(entry.Subdomain)).(string)
		comment := entry.Comment
		ttl := utils.IfThenElse(ttlErr == nil, entry.Ttl, boldRedText(strconv.Itoa(entry.Ttl)))

		t.AppendRow(table.Row{
			rtype,
			subdomain,
			comment,
			ip,
			ttl,
		})
	}

	t.Render()

	utils.PrintFLn("Checking Rate: %v", global.Configuration.Ddnsu.Rate)
	utils.PrintFLn("Public Ip: %v", ip)

	if errorPresent {
		fmt.Println("Configuration file has an invalid property somewhere. Please look through the table summary to find the entries that are red.")
	}
}
