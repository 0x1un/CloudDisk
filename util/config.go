package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	Location string `yaml:"location"`
	Port     string `yaml:"port"`
}

var (
	Conf *config
)

func getConfig() *config {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("read config file is failed.")
	}
	conf := &config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Println("Failed writting config file into the struct.")
	}
	return conf
}

func init() {
	Conf = getConfig()
}
