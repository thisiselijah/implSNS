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
        Region string `yaml:"region"`
        Endpoint string `yaml:"endpoint"`
        AccessKey string `yaml:"access_key"`
        SecretKey string `yaml:"secret_key"`
        SessionToken string `yaml:"session_token"`
    } `yaml:"dynamodb"`
    
}

// LoadConfig 讀取 config.yaml 並回傳 Config 結構
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
