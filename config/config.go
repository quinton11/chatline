package config

import (
	"encoding/json"
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
		return err
	}

	dirPath, filePath, err := GetConfigFileDetails()
	if err != nil {
		return err
	}

	//dir
	err = PathCheck(dirPath, os.ModeDir)
	if err != nil {
		return err
	}

	//file
	var file *os.File
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(filePath, 0, os.ModePerm)
		if err != nil {
			return err
		}
	}

	defer file.Close()

	_, err = file.Write(b)
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
