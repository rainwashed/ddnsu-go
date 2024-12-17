package commands

import (
	"ddnsu/v2/src/global"
	"ddnsu/v2/src/services/cloudflare"
	"ddnsu/v2/src/utils"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func cloudflareCleanup(configuration global.DDNSUConfig, token string) {
	zoneId, zErr := cloudflare.ReturnZoneIdFromDomain(configuration.Ddnsu.Domain, token)
	remoteRecords, rErr := cloudflare.ListDnsRecords(zoneId, token)

	if zErr != nil || rErr != nil {
		fmt.Println(color.RedString("zErr or rErr had a problem. zErr: %v. rErr: %v", zErr, rErr))
	}

	for _, rRecord := range remoteRecords {
		if !strings.HasPrefix(rRecord.Comment, global.RecordManagedPrefix) {
			continue
		}

		existsInConfig := false
		rRecordSerial := utils.SerializeRecord(global.DDNSURecord{
			Name:    rRecord.Name,
			Comment: rRecord.Comment,
			Ttl:     rRecord.Ttl,
			Type:    rRecord.Type,
		})
		for _, lRecord := range configuration.Ddnsu.Record {
			lRecordSerial := utils.SerializeRecord(global.DDNSURecord{
				Name:    lRecord.Subdomain,
				Comment: lRecord.Comment,
				Ttl:     lRecord.Ttl,
				Type:    lRecord.Rtype,
			})

			if lRecordSerial == rRecordSerial {
				existsInConfig = true
			}
		}

		if !existsInConfig {
			deletedRecord, errDeletingRecord := cloudflare.DeleteDnsRecord(rRecord.Id, zoneId, token)

			if !deletedRecord {
				fmt.Println(color.RedString("failed to delete record: %v. error: %v", rRecordSerial, errDeletingRecord))
			} else {
				fmt.Println(color.GreenString("successfully deleted record: %v", rRecordSerial))
			}
		}
	}
}

func CleanupCommand(cmd *cobra.Command, args []string) {
	switch strings.ToLower(args[0]) {
	case "vercel":
		break
	case "cloudflare":
		cloudflareCleanup(global.Configuration, global.Token)
	}
}
