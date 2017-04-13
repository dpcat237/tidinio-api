package app_conf

import (
	"io/ioutil"
	"github.com/olebedev/config"
)

const configFile = "app/config/config.yml"

var Data = config.Config{};

func LoadConfiguration() {
	strBytes, _ := ioutil.ReadFile(configFile)
	cfg, _ := config.ParseYaml(string(strBytes))
	Data = (*cfg)
}
