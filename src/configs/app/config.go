package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	App      App
	Database Database
	Google   Google
}

type App struct {
	Port        int
	Name        string
	Environment string
	Locale      string
	Key         string
	Debug       bool
	MigrateKey  string
}

type Database struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type Google struct {
	ClientID             string
	ClientSecret         string
	PrivateStorageBucket string
	PublicStorageBucket  string
}

func ParseConfig(v *viper.Viper) *Config {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		panic(err)
	}

	return &c
}

func loadConfig(environment string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(fmt.Sprintf("../config/config-%s", environment))
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func GetConfig() *Config {
	if config == nil {
		v, err := loadConfig(os.Getenv("environment"))
		if err != nil {
			panic(err)
		}
		err = v.Unmarshal(&config)
		if err != nil {
			log.Printf("unable to decode into struct, %v", err)
			panic(err)
		}
	}
	return config
}
