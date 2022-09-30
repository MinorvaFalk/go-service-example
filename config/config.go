package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type config struct {
	Dsn  string `mapstructure:"dsn"`
	Port string `mapstructure:"port"`
}

// Read config from named file (config.yaml)
func InitConfig() *config {
	var c config

	viper.AddConfigPath(filepath.Join(rootDir(), "config"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found\n%v", err))
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("failed to read configuration\n%v", err))
	}

	fmt.Printf("config : %v\n", c)

	return &c
}

// Read config from .env (unnamed file)
func ReadEnv() *config {
	var c config

	viper.AddConfigPath(filepath.Join(rootDir(), "config"))

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found\n%v", err))
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("failed to read configuration\n%v", err))
	}

	fmt.Printf("config : %v\n", c)

	return &c
}

// Fetch root directory path
func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
