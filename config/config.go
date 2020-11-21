package config

import (
	"errors"

	"github.com/BurntSushi/toml"
)

func NewWithPath(path string, check byte) (*ConfigData, error) {
	var config ConfigData
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, errors.New("Config File Load Failed: " + err.Error())
	}

	if checkDatabase&check == checkDatabase {
		if err = checkWithDefConfig(&config.Database); err != nil {
			return nil, err
		}
	}
	if checkWebsite&check == checkWebsite {
		if err = checkWithDefConfig(&config.Website); err != nil {
			return nil, err
		}
	}
	if checkAuthInfo&check == checkAuthInfo {
		if err = checkWithDefConfig(&config.AuthInfo); err != nil {
			return nil, err
		}
	}
	if checkServer&check == checkServer {
		if err = checkWithDefConfig(&config.Server); err != nil {
			return nil, err
		}
	}
	if checkRedis&check == checkRedis {
		if err = checkWithDefConfig(&config.Redis); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func checkWithDefConfig(c BaseConfig) error {
	if c == nil {
		return nil
	}
	return c.check()
}
