package services

import (
	"ddnsu/v2/src/global"
	"ddnsu/v2/src/services/cloudflare"
	"ddnsu/v2/src/utils"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func OnLoginCommandRun(cmd *cobra.Command, args []string) {
	provider := strings.ToLower(args[0])
	token := args[1]

	switch provider {
	case "vercel":
		fmt.Print("Using Vercel provider")
	case "cloudflare":
		validToken, err := cloudflare.TestToken(token)

		if err != nil || !validToken {
			utils.PrintFLn("token does not appear to be valid or correctly-set. error: %v", err)
			panic("invalid token")
		}

		fmt.Println("token passes checking. adding to configuration file")
		originalToken := global.Configuration.Services.Cloudflare.Token
		global.Configuration.Services.Cloudflare.Token = token

		utils.PromptWriteConfirm(fmt.Sprintf("cloudflare.token='%v'", originalToken), fmt.Sprintf("cloudflare.token='%v'", global.Configuration.Services.Cloudflare.Token), global.ConfigurationPath)
	}
}
