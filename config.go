package main

import (
	"encoding/json"
	"os"
)

const (
	configFileName = "config.json"
)

var (
	config Config
)

// Config contains current settings of program
type Config struct {
	ServiceKey string `json:"servicekey"`
}

// Save saves config to default configFileName
func (c Config) Save() error {
	w, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer w.Close()

	// 현재 설정으로 기본 config 파일 생성
	prettyConfig, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(prettyConfig)
	return err
}

func loadConfig() error {
	if !isConfigValid() {
		config.ServiceKey = getServiceKey()
		return config.Save()
	}

	confR, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	defer confR.Close()
	jDec := json.NewDecoder(confR)
	err = jDec.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}

func isConfigValid() bool {
	if !isExist(configFileName) {
		return false
	}

	confR, err := os.Open(configFileName)
	if err != nil {
		panic(err)
	}
	defer confR.Close()
	jDec := json.NewDecoder(confR)
	err = jDec.Decode(&config)
	if err != nil {
		panic(err)
	}

	return true
}

func getServiceKey() string {
	serviceKey := os.Getenv("SERVICEKEY")
	if serviceKey != "" {
		return serviceKey
	}

	if config.ServiceKey != "" {
		return config.ServiceKey
	}

	panic("no servicekey")
}
