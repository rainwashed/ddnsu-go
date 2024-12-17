package utils

import (
	"fmt"
	"os"
	"reflect"

	"ddnsu/v2/src/global"

	"github.com/manifoldco/promptui"
	"github.com/pelletier/go-toml/v2"
)

func PrintFLn(
	format string,
	a ...any,
) {
	fmt.Printf(format+"\n", a...)
}

func returnConfigFile(configFilePath string) (*global.DDNSUConfig, error) {
	fileContent, readErr := os.ReadFile(configFilePath)
	if readErr != nil {
		return &global.DDNSUConfig{}, fmt.Errorf("error while reading the configuration file at %v: %w", configFilePath, readErr)
	}

	var config global.DDNSUConfig
	unmarshalErr := toml.Unmarshal(fileContent, &config)
	if unmarshalErr != nil {
		return &global.DDNSUConfig{}, fmt.Errorf("error while parsing the configuration file at %v: %w", configFilePath, unmarshalErr)
	}

	return &config, nil
}

func ConvertByteArrayToStruct(bytearray []byte, targetType reflect.Type) error {
	return nil
}

func LoadConfigurationIntoGlobalVar(configFilePath string) (bool, error) {
	configuration, err := returnConfigFile(configFilePath)

	if err != nil {
		return false, fmt.Errorf("loading configuration file had an error %v", err)
	}

	//	fmt.Printf("configuration: %v\n", configuration)

	global.Configuration = *configuration

	return true, nil
}

func StoreActiveTokenInGlobalVar(configuration global.DDNSUConfig) {
	provider := configuration.Ddnsu.Use

	switch provider {
	case "cloudflare":
		global.Token = configuration.Services.Cloudflare.Token
	case "vercel":
		global.Token = global.Configuration.Services.Vercel.Token
	}

}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func PromptWriteConfirm(messageBefore string, messageAfter string, pathToWriteTo string) {
	fmt.Println(messageBefore + " -→ " + messageAfter)
	writePrompt := promptui.Select{
		Label: fmt.Sprintf("Confirm write to %v?", pathToWriteTo),
		Items: []string{"Yes", "No"},
	}

	_, answer, _ := writePrompt.Run()

	if answer == "Yes" {
		tomlRepresentation, tomlError := toml.Marshal(global.Configuration)
		if tomlError != nil {
			panic("error when marshalling object")
		}
		os.WriteFile(pathToWriteTo, tomlRepresentation, 0644)
		fmt.Println("Changes have been written.")
	} else {
		fmt.Println("No file was changed.")
	}
}

// IfThenElse evaluates a condition, if true returns the first parameter otherwise the second
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}
