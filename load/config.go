package load

import (
	"encoding/json"
	"os"
)

type ChatConfig struct {
	Model      string `json:"model"`
	Prompt     string `json:"prompt"`
	TokenSpent bool   `json:"token_spent"`
}

var configFile = "config"

func LoadConfig() (ChatConfig, error) {
	var cfg ChatConfig
	file, err := os.Open(configFile)
	if err != nil {
		return cfg, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return cfg, err
}

func SaveConfig(cfg ChatConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}
