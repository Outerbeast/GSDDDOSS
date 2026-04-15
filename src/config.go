package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ListenerAddr       string `json:"listener_addr"`
	ListenerPort       int    `json:"listener_port"`
	WindowsRuleGrouped bool   `json:"windows_rule_ip_grouped"`
	CommandAddBlock    string `json:"command_add_block"`
}

func configLoad() (*Config, error) {
	file, err := os.Open("config.json")

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Warning: Configuration file not found. (A default one will be created.)")
			defaultConfig := Config{
				ListenerAddr:       "127.0.0.1",
				ListenerPort:       8008,
				WindowsRuleGrouped: true,
				CommandAddBlock:    "",
			}

			if err := configSave(defaultConfig); err != nil {
				return nil, err
			}

			return &defaultConfig, nil
		}

		return nil, err
	}

	defer file.Close()

	var data Config
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func configSave(data Config) error {
	file, err := os.Create("config.json")

	if err != nil {
		return err
	}

	defer file.Close()
	json.NewEncoder(file).Encode(data)

	return nil
}
