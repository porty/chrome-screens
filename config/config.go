package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const filename = "config.json"

type Config struct {
	SlackAPIToken string `envconfig:"slack_api_token" required:"true" json:"-"`
	LastURL       string `json:"lastUrl,omitempty"`
}

func Load() *Config {
	config := Config{}
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file: " + err.Error())
	}
	if err := envconfig.Process("chromescreens", &config); err != nil {
		log.Print("Error reading config from environment: " + err.Error())
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Print(filename + " doesn't exist")
		} else {
			log.Printf("There is a problem reading %s: %s", filename, err.Error())
		}

	} else {
		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Printf("%s is corrupted: %s", filename, err.Error())
		}
	}

	return &config
}

func (c *Config) Save() error {
	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}
