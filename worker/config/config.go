package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

type ConfigStruct struct {
	Id              string
	GameDir         string
	RequestInterval int
	Server          string
}

var Cfg *ConfigStruct

func ReadFile(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.UnmarshalStrict(file, &Cfg)
	if err != nil {
		return err
	}

	return nil
}
