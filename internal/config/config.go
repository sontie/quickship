package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string             `yaml:"version"`
	Hosts    map[string]Host    `yaml:"hosts"`
	Projects []Project          `yaml:"projects"`
}

type Host struct {
	Addr string `yaml:"addr"`
	User string `yaml:"user"`
}

type Project struct {
	Name    string              `yaml:"name"`
	Repo    string              `yaml:"repo"`
	Path    string              `yaml:"path"`
	Envs    map[string][]string `yaml:"envs"`
	Scripts map[string]string   `yaml:"scripts"`
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

	return &cfg, nil
}
