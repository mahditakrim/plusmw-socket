package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		LogPath   string    `yaml:"log_path"`
		Transport transport `yaml:"transport"`
	}

	transport struct {
		Websocket     websocket     `yaml:"websocket"`
		MessageBroker messageBroker `yaml:"message_broker"`
	}

	websocket struct {
		Addr string `yaml:"addr"`
	}

	messageBroker struct {
		Nats struct {
			Addr     string            `yaml:"addr"`
			Username string            `yaml:"username"`
			Password string            `yaml:"password"`
			Subjects map[string]string `yaml:"subjects"`
		} `yaml:"nats"`
	}
)

func Init() (*Config, error) {

	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(yamlFile, conf)

	return conf, err
}
