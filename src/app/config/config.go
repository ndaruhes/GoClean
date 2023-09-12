package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func loadConfig(environment string) (*viper.Viper, error) {
	v := viper.New()
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fmt.Println("ANZAYYY", currentDir)

	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		if info.IsDir() {
			fmt.Println("Folder:", path)
		} else {
			fmt.Println("File:", path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	currentDir = strings.TrimSuffix(currentDir, "src")
	v.SetConfigName(fmt.Sprintf("config/config-%s", environment))
	v.AddConfigPath(currentDir)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func GetConfig() *Config {
	if config == nil {
		v, err := loadConfig("local")
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
