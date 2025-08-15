package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	globalConfig Configuration
)

type Configuration struct {
	GateIO struct {
		Key    string `mapstructure:"key"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"gateio"`
	WxPusher struct {
		AppToken string `mapstructure:"app_token"`
		UserID   string `mapstructure:"user_id"`
	} `mapstructure:"wxpusher"`
	Notification struct {
		Windows string `mapstructure:"windows"`
	} `mapstructure:"notification"`
	Strategy struct {
		RMI         bool `mapstructure:"rmi"`
		Suppertrend bool `mapstructure:"suppertrend"`
		Agg         bool `mapstructure:"agg"`
	} `mapstructure:"strategy"`
}

func Init() *Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("read config file failed: ", err)
			return
		}
		if err := viper.Unmarshal(&globalConfig); err != nil {
			fmt.Println("unmarshal config file failed: ", err)
			return
		}
	})

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config file failed: ", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&globalConfig); err != nil {
		fmt.Println("unmarshal config file failed: ", err)
		os.Exit(1)
	}

	return &globalConfig
}

func GetConfig() *Configuration {
	return &globalConfig
}
