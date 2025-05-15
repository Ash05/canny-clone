package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	ApiUrl            string `json:"apiUrl"`
	DatabaseURL       string `json:"databaseUrl"`
	Port              string `json:"port"`
	GoogleClientID    string `json:"googleClientId"`
	GoogleClientSecret string `json:"googleClientSecret"`
	GoogleRedirectURL string `json:"googleRedirectUrl"`
	JWTSecret         string `json:"jwtSecret"`
}

var config Configuration

func LoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	file, err := ioutil.ReadFile("config." + env + ".json")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	if err := json.Unmarshal(file, &config); err != nil {
		panic("Failed to parse config: " + err.Error())
	}
}

func GetConfig() Configuration {
	return config
}
