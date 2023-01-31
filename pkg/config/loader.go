package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

func LoadConfig(reader io.Reader) (*Config, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf := Config{}
	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
