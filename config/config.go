package config

import (
	"github.com/joho/godotenv"
)

func LoadConfig() error {
	env := ".env"
	err := godotenv.Load(env)
	if err != nil {
		return err
	}
	return nil
}
