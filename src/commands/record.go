package commands

import (
	"ddnsu/v2/src/global"
	"ddnsu/v2/src/utils"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func recordTypeValidation(input string) error {

	if validTypes := []string{"A", "Alias", "CAA", "CNAME", "HTTPS", "MX", "SRV", "TXT", "NS"}; slices.Contains(validTypes, strings.ToUpper(input)) {
		return nil
	}
	return fmt.Errorf("%v is not a valid record type", input)
}

func ttlValidation(input string) error {
	if input == "" {
		return nil
	}

	n, err := strconv.Atoi(input)

	if err != nil {
		return fmt.Errorf("%v is not a number or 'auto'", input)
	}

	if n < 1 {
		return fmt.Errorf("TTL must be greater than or equal to 1")
	}

	return nil
}

var subdomainValidationRegex, _ = regexp.Compile(`^\*|\@$|^[A-Za-z0-9]{1,}$`)

func subdomainValidation(input string) error {
	match := subdomainValidationRegex.MatchString(input)

	if match {
		return nil
	} else {
		return fmt.Errorf("%v is not a valid subdomain", input)
	}
}

func commentValidation(input string) error {
	return nil
}

func HandleCommand(commandType string) error {
	switch commandType {
	case "add":
		recordTypePrompt := promptui.Prompt{
			Label:    "Record Type",
			Validate: recordTypeValidation,
		}
		recordResult, recordErr := recordTypePrompt.Run()

		ttlPrompt := promptui.Prompt{
			Label:    "Time Until Live (default: 1 = auto)",
			Validate: ttlValidation,
		}
		ttlResult, ttlErr := ttlPrompt.Run()
		ttlResultI, _ := strconv.Atoi(ttlResult)

		subdomainPrompt := promptui.Prompt{
			Label:    "Subdomain/Name",
			Validate: subdomainValidation,
		}
		subdomainResult, subdomainErr := subdomainPrompt.Run()

		commentPrompt := promptui.Prompt{
			Label:    "Comment (default: '')",
			Validate: commentValidation,
		}
		commentResult, commentErr := commentPrompt.Run()

		if recordErr != nil || ttlErr != nil || subdomainErr != nil || commentErr != nil {
			return fmt.Errorf("record, ttl, subdomain, or comment had an error when running the prompt")
		}

		fmt.Printf("Will %v a record with:\nType: %v\nTTL: %v\nSubdomain: %v\nComment: %v\n", commandType, recordResult, ttlResult, subdomainResult, commentResult)

		newRecord := global.Record{
			Rtype:     recordResult,
			Ttl:       ttlResultI,
			Subdomain: subdomainResult,
			Comment:   commentResult,
		}

		global.Configuration.Ddnsu.Record = append(global.Configuration.Ddnsu.Record, newRecord)

		utils.PromptWriteConfirm("", "Add this record?", global.ConfigurationPath)
	case "update":
		{
			var records []string = make([]string, len(global.Configuration.Ddnsu.Record))
			valueColor := color.New(color.Bold, color.BgWhite).SprintfFunc()

			for i, r := range global.Configuration.Ddnsu.Record {
				iS := strconv.Itoa(i)
				records[i] = fmt.Sprintf("T:%v-S:%v-C:%v-T:%v-I:%v", valueColor(r.Rtype), valueColor(r.Subdomain), valueColor(r.Comment), valueColor(strconv.Itoa(r.Ttl)), iS)

			}

			recordUpdatePrompt := promptui.Select{
				Label: "Select Record",
				Items: records,
			}

			_, result, promptErr := recordUpdatePrompt.Run()

			if promptErr != nil {
				return fmt.Errorf("creating prompt returned an error: %v", promptErr)
			}

			stringSplitArray := strings.Split(result, ":")
			indexSelected := stringSplitArray[len(stringSplitArray)-1]
			indexSelectedI, _ := strconv.Atoi(indexSelected)

			utils.PrintFLn("Please enter the new properties of %v", result)
			recordTypePrompt := promptui.Prompt{
				Label:    "Record Type",
				Validate: recordTypeValidation,
			}
			recordResult, recordErr := recordTypePrompt.Run()

			ttlPrompt := promptui.Prompt{
				Label:    "Time Until Live (default: 0, auto)",
				Validate: ttlValidation,
			}
			ttlResult, ttlErr := ttlPrompt.Run()
			ttlResultI, _ := strconv.Atoi(ttlResult)

			subdomainPrompt := promptui.Prompt{
				Label:    "Subdomain/Name",
				Validate: subdomainValidation,
			}
			subdomainResult, subdomainErr := subdomainPrompt.Run()

			commentPrompt := promptui.Prompt{
				Label:    "Comment (default: '')",
				Validate: commentValidation,
			}
			commentResult, commentErr := commentPrompt.Run()

			if recordErr != nil || ttlErr != nil || subdomainErr != nil || commentErr != nil {
				return fmt.Errorf("record, ttl, subdomain, or comment had an error when running the prompt")
			}

			fmt.Printf("Will %v record %v to:\nType: %v\nTTL: %v\nSubdomain: %v\nComment: %v\n", commandType, result, recordResult, ttlResult, subdomainResult, commentResult)
			newRecord := global.Record{
				Rtype:     recordResult,
				Ttl:       ttlResultI,
				Subdomain: subdomainResult,
				Comment:   commentResult,
			}
			global.Configuration.Ddnsu.Record[indexSelectedI] = newRecord

			utils.PromptWriteConfirm("", "Update this record?", global.ConfigurationPath)
		}
	case "delete":
		var records []string = make([]string, len(global.Configuration.Ddnsu.Record))
		valueColor := color.New(color.Bold, color.BgWhite).SprintfFunc()

		for i, r := range global.Configuration.Ddnsu.Record {
			iS := strconv.Itoa(i)
			records[i] = fmt.Sprintf("T:%v-S:%v-C:%v-T:%v-I:%v", valueColor(r.Rtype), valueColor(r.Subdomain), valueColor(r.Comment), valueColor(strconv.Itoa(r.Ttl)), iS)
		}

		recordUpdatePrompt := promptui.Select{
			Label: "Select Record",
			Items: records,
		}

		_, result, promptErr := recordUpdatePrompt.Run()

		if promptErr != nil {
			return fmt.Errorf("creating prompt returned an error: %v", promptErr)
		}

		stringSplitArray := strings.Split(result, ":")
		indexSelected := stringSplitArray[len(stringSplitArray)-1]
		indexSelectedI, _ := strconv.Atoi(indexSelected)

		global.Configuration.Ddnsu.Record = slices.Delete(global.Configuration.Ddnsu.Record, indexSelectedI, indexSelectedI+1)

		utils.PromptWriteConfirm("", "Delete this record?", global.ConfigurationPath)
	}

	return nil
}
