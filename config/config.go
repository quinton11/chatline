package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/quinton11/chatline/internal/utils"
)

func LoadConfig() error {
	env := ".env"
	err := godotenv.Load(env)
	if err != nil {
		return err
	}
	return nil
}

const (
	configFile = "chatline.json"
	configDir  = ".chatline"
)

func ConfigFileCheck() {
}

// returns config file path and dir path
func GetConfigFileDetails() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	dirPath := path.Join(home, configDir)
	filePath := path.Join(dirPath, configFile)

	return dirPath, filePath, nil
}

func WriteConfig(room utils.RoomConfig) error {
	//check if file exists
	b, err := json.Marshal(room)
	if err != nil {
		fmt.Println("Error in marshalling")
		return err
	}

	dirPath, filePath, err := GetConfigFileDetails()
	if err != nil {
		fmt.Println("Get config error")
		return err
	}

	//dir
	err = PathCheck(dirPath, os.ModeDir)
	if err != nil {
		fmt.Println("Path check error")

		return err
	}

	//file
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("In not exist error")
		_, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Create error")
			return err
		}
	}
	err = os.WriteFile(filePath, b, 0600)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfig(room *utils.RoomConfig) error {
	dirPath, filePath, err := GetConfigFileDetails()
	if err != nil {
		return err
	}

	err = PathCheck(dirPath, os.ModeDir)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, room)
	if err != nil {
		return err
	}

	return nil
}

func PathCheck(dirpath string, mode fs.FileMode) error {
	if _, err := os.Stat(dirpath); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(dirpath, mode); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}
