package utils

import (
	"context"
	"ddnsu/v2/src/global"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

func TestNetworking() bool {
	netreq, netok := http.Get("https://example.com")
	if netok != nil {
		fmt.Println(netok.Error())
		return false
	}

	if netreq.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

var client = &http.Client{
	Transport: &http.Transport{
		IdleConnTimeout: 1 * time.Second,
	},
}

func MakeIpConsensus() string {
	ipProviders := global.Configuration.Ddnsu.IpProviders
	var ips map[string]int = make(map[string]int)

	// NOTE: i will become _
	for _, provider := range ipProviders {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, "GET", provider, nil)
		resp, errorResp := client.Do(req)

		// TODO: make this a context goroutine some shit

		if errorResp != nil {
			fmt.Printf("Provider %v returned an error when requested. Skipping, will not be added to consensus.\n", provider)
			continue
		}

		defer resp.Body.Close()

		type NeededResponse struct {
			Ip string
		}

		var responseObject = NeededResponse{}
		err := json.NewDecoder(resp.Body).Decode(&responseObject)

		if err != nil {
			continue
		}

		ips[responseObject.Ip]++

	}

	key := ""
	value := 0
	for k, v := range ips {
		if v > value {
			key = k
			value = v
		}
	}

	return key

}

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

func DetermineIfNeedConfigCreationAndCreateIfDoesNotExist(path string, fileName string, embedContent []byte) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.FileMode(0777))
		file, err := os.Create(filepath.Join(path, fileName))
		if err != nil {
			fmt.Println(color.RedString("could not create configuration file. please do so manually (stage 1)."))
			os.Exit(1)
		}

		_, errWrite := file.Write(embedContent)

		if errWrite != nil {
			fmt.Println(color.RedString("could not create configuration file. please do so manually (stage 2)."))
		}

		fmt.Println(color.New(color.FgGreen, color.Italic).Sprintf("configuration file could not be located; created an example configuration file at %v\n", filepath.Join(path, fileName)))

		defer file.Close()

		return true
	} else {
		return false
	}
}
