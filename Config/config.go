package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Name          string `yaml:"Name"`
	Version       string `yaml:"Version"`
	Port          int    `yaml:"port"`
	Node          string `yaml:"node"`
	DataDirectory string `yaml:"DataDirectory"`
	DataFileName  string `yaml:"DataFile"`
	Token         string `yaml:"Token"`
	ServiceToken  string `yaml:"ServiceToken"`
	IpHeader      string `yaml:"IpHeader"`
	Remote        string `yaml:"remote"`
}

var data, _ = os.ReadFile("./config.yml")
var Conf = ConfigStruct{}
var _ = yaml.Unmarshal(data, &Conf)
