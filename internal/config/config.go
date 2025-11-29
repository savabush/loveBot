package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Cfg *Config = NewConfig()

type Config struct {
	TelegramBotToken string
	UserID1          int64
	UserID2          int64
}

func NewConfig() *Config {

	configName := ".env"
	viper.AddConfigPath(".")
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("config file %v not found \ncreating config from user input", configName)

		log.Println("Please set TELEGRAM_BOT_TOKEN:")
		var telegramBotToken string
		_, err := fmt.Scanln(&telegramBotToken)
		if err != nil {
			log.Fatalf("failed to read TELEGRAM_BOT_TOKEN: %v", err)
		}
		viper.Set("TELEGRAM_BOT_TOKEN", telegramBotToken)

		log.Println("Please set USER_ID_1:")
		var userID1 int64
		_, err = fmt.Scanln(&userID1)
		if err != nil {
			log.Fatalf("failed to read USER_ID_1: %v", err)
		}
		viper.Set("USER_ID_1", userID1)

		log.Println("Please set USER_ID_2:")
		var userID2 int64
		_, err = fmt.Scanln(&userID2)
		if err != nil {
			log.Fatalf("failed to read USER_ID_2: %v", err)
		}
		viper.Set("USER_ID_2", userID2)

		err = viper.WriteConfigAs(".env")
		if err != nil {
			log.Fatalf("failed to write config file: %v", err)
		}
		log.Println("Config file created successfully.")
	}
	return &Config{
		TelegramBotToken: viper.GetString("TELEGRAM_BOT_TOKEN"),
		UserID1:          viper.GetInt64("USER_ID_1"),
		UserID2:          viper.GetInt64("USER_ID_2"),
	}
}
