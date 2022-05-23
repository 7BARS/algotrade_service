package config

import (
	"fmt"
	"io/ioutil"
	
	"algotrade_service/internal/utils"

	"gopkg.in/yaml.v3"
)

func GetConfig(pathConfig string) (*Config, error) {
	ok, err := utils.Exists(pathConfig)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("config file: %s, does not exists", pathConfig)
	}

	data, err := ioutil.ReadFile(pathConfig)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
