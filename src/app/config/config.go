package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	App           App
	Cors          Cors
	Database      Database
	Elasticsearch Elasticsearch
	Google        Google
	OutboundURL   OutboundURL
}

type App struct {
	Port          int
	Name          string
	AliasName     string
	Environment   string
	DefaultLocale string
	Key           string
	Debug         bool
	MigrateKey    string
}

type Cors struct {
	AllowOrigins     string
	AllowHeaders     string
	AllowMethods     string
	AllowCredentials bool
	ExposeHeaders    string
}

type Database struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type Elasticsearch struct {
	Username string
	Password string
}

type Google struct {
	ClientID             string
	ClientSecret         string
	PrivateStorageBucket string
	PublicStorageBucket  string
}

type OutboundURL struct {
	ElasticsearchOutbound string
}

func loadConfig(appEnvironment string, appRootFolder string) (*viper.Viper, error) {
	v := viper.New()
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	printFolderStructure(currentDir)
	fmt.Println("ENVIRONMENT = ", appEnvironment)
	fmt.Println("APP ROOT FOLDER = ", appRootFolder)

	rootFolderIndex := strings.LastIndex(currentDir, appRootFolder)
	if rootFolderIndex == -1 {
		return nil, fmt.Errorf("root folder " + appRootFolder + "tidak ditemukan dalam path")
	}

	configPath := filepath.Join(currentDir[:rootFolderIndex+len(appRootFolder)], "config")
	v.SetConfigName(fmt.Sprintf("config-%s", appEnvironment))
	v.AddConfigPath(configPath)
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New(fmt.Sprintf("config-%s file not found", appEnvironment))
		}
		return nil, err
	}

	return v, nil
}

func GetConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if config == nil {
		v, err := loadConfig(os.Getenv("APP_ENVIRONMENT"), os.Getenv("APP_ROOT_FOLDER"))
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

func printFolderStructure(currentDir string) {
	fmt.Println("")
	fmt.Println("FOLDER STRUCTURE")
	fmt.Println("------------------------------------------------------------------------")

	// Daftar folder yang ingin dikecualikan
	excludedFolders := map[string]bool{
		".idea":   true,
		".git":    true,
		".vscode": true,
		"tmp":     true,
	} // Ganti dengan nama folder yang ingin Anda kecualikan

	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Periksa apakah folder ada dalam daftar yang ingin dikecualikan
		if _, ok := excludedFolders[info.Name()]; ok && info.IsDir() {
			return filepath.SkipDir // Lewati folder ini dan isinya
		}

		// Hitung berapa banyak karakter indent yang diperlukan
		relativePath, _ := filepath.Rel(currentDir, path)
		depth := strings.Count(relativePath, string(filepath.Separator))
		indent := strings.Repeat("  ", depth)

		// Cetak folder atau file dengan indent
		if info.IsDir() {
			fmt.Printf("%sFolder: %s\n", indent, path)
		} else {
			fmt.Printf("%sFile: %s\n", indent, path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("------------------------------------------------------------------------")
}
