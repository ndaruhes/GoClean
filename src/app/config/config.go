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

	printFolderStructure(currentDir)

	rootFolderIndex := strings.LastIndex(currentDir, "GoClean")
	if rootFolderIndex == -1 {
		return nil, fmt.Errorf("root folder GoClean tidak ditemukan dalam path")
	}

	configPath := filepath.Join(currentDir[:rootFolderIndex+len("GoClean")], "config")
	v.SetConfigName(fmt.Sprintf("config-%s", environment))
	v.AddConfigPath(configPath)
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
