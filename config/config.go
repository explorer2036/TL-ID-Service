package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// dbStruct defines fields for databases
type dbStruct struct {
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Name   string `yaml:"name"`
}

// serverStruct defines fields for main logic
type serverStruct struct {
	ListenAddr string `yaml:"listen_addr"`
}

// logStruct defines fields for log
type logStruct struct {
	OutputLevel        string `yaml:"output_level"`
	OutputPath         string `yaml:"output_path"`
	RotationPath       string `yaml:"rotation_path"`
	RotationMaxSize    int    `yaml:"rotation_max_size"`
	RotationMaxAge     int    `yaml:"rotation_max_age"`
	RotationMaxBackups int    `yaml:"rotation_max_backups"`
	JSONEncoding       bool   `yaml:"json_encoding"`
}

// Config structure for server
type Config struct {
	Server serverStruct `yaml:"server"`
	DB     dbStruct     `yaml:"db"`
	Log    logStruct    `yaml:"log"`
}

// ParseYamlFile the config file
func ParseYamlFile(filename string, c *Config) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}
