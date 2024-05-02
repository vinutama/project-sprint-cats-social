package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfigs *envConfigs

func init() {
	EnvConfigs = loadENV()
}

type envConfigs struct {
	DbName     string `mapstructure:"DB_NAME"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbUser     string `mapstructure:"DB_USERNAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbPort     string `mapstructure:"DB_PORT"`
	BcryptSalt string `mapstructure:"BCRYPT_SALT"`
	JwtSecret  string `mapstructure:"JWT_SECRET"`
}

func loadENV() (config *envConfigs) {
	viper.SetConfigFile("cats-social.env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error when reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return config
}
