package config

import (
	"github.com/spf13/viper"
	"log"
)

var AppConfig Config

type App struct {
	Name           string `yaml:"name"`
	Version        string `yaml:"version"`
	TZ             string
	ClusterContext string
}

type Http struct {
	Port string `yaml:"port"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

type Log struct {
	Level string `yaml:"log_level"`
}

type Config struct {
	App
	Http
	Db DbConfig
	Log
	Debug       bool
	Environment string
}

func InitializeAppConfig() {
	viper.SetConfigName("../.env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	AppConfig.Http.Port = viper.GetString("PORT")
	AppConfig.Environment = viper.GetString("ENVIRONMENT")
	AppConfig.Debug = viper.GetBool("DEBUG")

	AppConfig.Db.Host = viper.GetString("DB_HOST")
	AppConfig.Db.Port = viper.GetInt("DB_PORT")
	AppConfig.Db.Database = viper.GetString("DB_DATABASE")
	AppConfig.Db.Username = viper.GetString("DB_USERNAME")
	AppConfig.Db.Password = viper.GetString("DB_PASSWORD")
	AppConfig.App.TZ = viper.GetString("TZ")
	AppConfig.App.ClusterContext = viper.GetString("CLUSTER_CONTEXT")

	log.Println("[INIT] configuration loaded")
}
