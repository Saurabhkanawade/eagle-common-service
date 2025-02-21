package twilio_config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type TwilioConfig struct {
	AccountSID  string
	AuthToken   string
	PhoneNumber string
}

func LoadTwilioConfig() TwilioConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return TwilioConfig{
		AccountSID:  os.Getenv("TWILIO_SID"),
		AuthToken:   os.Getenv("TWILIO_AUTH_TOKEN"),
		PhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
	}
}
