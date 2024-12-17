package services

import (
	"ddnsu/v2/src/global"
	"ddnsu/v2/src/services/cloudflare"
	"ddnsu/v2/src/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// @return true if record should be created (does not exist in remote records), return false if record should be updated
func DetermineIfRecordShouldBeCreated(toBeCreated global.DDNSURecord, cloudflareRecords []global.DDNSURecord) bool {
	toBeCreatedSerial := utils.SerializeRecord(toBeCreated)
	utils.PrintFLn("To be created Serial: %v", toBeCreatedSerial)
	for _, record := range cloudflareRecords {
		recordSerial := utils.SerializeRecord(record)
		utils.PrintFLn("Compare record serial: %v", recordSerial)
		if toBeCreatedSerial == recordSerial {
			return false
		}
	}
	return true
}

func CloudflareServiceLoop(updateFrequency time.Duration, ip string) {
	cloudflareZoneId, cloudflareZoneErr := cloudflare.ReturnZoneIdFromDomain(global.Configuration.Ddnsu.Domain, global.Token)
	if cloudflareZoneErr != nil {
		fmt.Println(color.RedString("could not retrieve zone id for domain: %v. have you ensured that your token allows read/write access to that domain?", global.Configuration.Ddnsu.Domain))
		os.Exit(1)
	}

	cloudflareRecords, cloudflareRecordsErr := cloudflare.ListDnsRecords(cloudflareZoneId, global.Token)

	if cloudflareRecordsErr != nil {
		fmt.Println(color.RedString("could not retrieve dns records for zone id: %v. have you ensured that your token allows read/write access to that domain?", cloudflareZoneId))
		os.Exit(1)
	}

	var managedRecords []global.ManagedRecord

	for _, configRecord := range global.Configuration.Ddnsu.Record {
		configRecordConverted := global.DDNSURecord{
			Name:    configRecord.Subdomain,
			Comment: global.RecordManagedPrefix + configRecord.Comment,
			Ttl:     configRecord.Ttl,
			Type:    configRecord.Rtype,
			Content: ip,
		}
		shouldBeCreated := DetermineIfRecordShouldBeCreated(configRecordConverted, cloudflareRecords)
		if shouldBeCreated {
			managedRecords = append(managedRecords, global.ManagedRecord{
				Record: configRecordConverted,
				Action: "create",
			})
		} else {
			for _, record := range cloudflareRecords {
				utils.PrintFLn("Record: %v", utils.SerializeRecord(record))
				if strings.HasPrefix(record.Comment, global.RecordManagedPrefix) {
					managedRecords = append(managedRecords, global.ManagedRecord{
						Record: record,
						Action: "update",
					})
				}
			}
		}
	}

	for _, managedRecord := range managedRecords {
		serialRecord := utils.SerializeRecord(managedRecord.Record)
		switch managedRecord.Action {
		case "create":
			_, recordAddErr := cloudflare.AddDnsRecord(cloudflareZoneId,
				managedRecord.Record.Type,
				managedRecord.Record.Name,
				strconv.Itoa(managedRecord.Record.Ttl),
				managedRecord.Record.Comment,
				ip,
				global.Token,
			)

			if recordAddErr != nil {
				fmt.Println(color.RedString("could not add %v, with error: %v", serialRecord, recordAddErr))
			} else {
				fmt.Println(color.GreenString("successfully created %v", serialRecord))
			}
		case "update":
			_, updateErr := cloudflare.UpdateDnsRecord(managedRecord.Record.Id, cloudflareZoneId, ip, global.Token)

			if updateErr != nil {
				fmt.Println(color.RedString("could not update %v, with error: %v", serialRecord, updateErr))
			} else {
				fmt.Println(color.GreenString("successfully updated %v", serialRecord))
			}
		}
	}
}

func BeginActiveLoop(updateFrequency time.Duration, provider string) {
	for {
		fmt.Println("active service has begun")
		ip := utils.MakeIpConsensus()

		if global.LastIpAddress != ip {
			fmt.Printf("ip: %v", ip)
			switch provider {
			case "cloudflare":
				CloudflareServiceLoop(updateFrequency, ip)
			default:
				fmt.Println(color.RedString("provider %v is not supported", provider))
				os.Exit(1)
			}

			global.LastIpAddress = ip
		} else {
			fmt.Println(color.YellowString("ip address has not changed; it is unnecessary to update anything."))
		}
		time.Sleep(updateFrequency)
	}
}
