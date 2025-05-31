package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	DynamoDB struct {
		Region       string `yaml:"region"`
		Endpoint     string `yaml:"endpoint"`
		AccessKey    string `yaml:"access_key"`
		SecretKey    string `yaml:"secret_key"`
		SessionToken string `yaml:"session_token"`
	} `yaml:"dynamodb"`
	JWT struct { // 新增 JWT 設定
		SecretKey     string `yaml:"secret_key"`
		ExpiryMinutes int    `yaml:"expiry_minutes"`
	} `yaml:"jwt"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    // 設定預設值 (如果需要)
    if cfg.JWT.SecretKey == "" {
        cfg.JWT.SecretKey = "default-secret-please-change" // 最好是必要欄位
    }
    if cfg.JWT.ExpiryMinutes == 0 {
        cfg.JWT.ExpiryMinutes = 60
    }
    return &cfg, nil
}