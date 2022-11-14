package config

import (
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	configStruct := &Config{}
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.SetConfigType("json")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./cmd/config") // path to look for the config file in
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		return nil, err
	}

	err = viper.Unmarshal(configStruct)
	if err != nil { // Handle errors reading the config file
		return nil, err
	}

	return configStruct, nil
}
