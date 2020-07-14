package config

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

type yamlParser struct {
}

func newYamlParser() ConfParser {
	return &yamlParser{}
}

// Unmarshal read file data and Unmarshal it into data
// if file is not exist, create it with initial data
func (parser *yamlParser) Unmarshal(conf interface{}, fromFile string) error {
	return YamlReadFileAndUnmarshal(conf, fromFile)
}

// Marshal marshal input data and write it into file
func (parser *yamlParser) Marshal(conf interface{}, toFile string) (err error) {
	err = YamlMarshalAndWriteFile(conf, toFile)
	return
}

func YamlReadFileAndUnmarshal(conf interface{}, fromFile string) error {
	if _, err := os.Stat(fromFile); os.IsNotExist(err) {
		YamlMarshalAndWriteFile(conf, fromFile)
		return err
	}

	data, err := ioutil.ReadFile(fromFile)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(data, conf); err != nil {
		return err
	}
	return nil
}

func YamlMarshalAndWriteFile(conf interface{}, filename string) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	// write the file
	if err = ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	return nil
}
