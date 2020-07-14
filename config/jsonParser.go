package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type jsonParser struct {
}

func newJSONParser() ConfParser {
	return &jsonParser{}
}

// Load read file data and Unmarshal it into data
// if file is not exist, create it with initial data
func (parser *jsonParser) Unmarshal(conf interface{}, fromFile string) error {
	return YamlReadFileAndUnmarshal(conf, fromFile)
}

// UnLoad marshal input data and write it into file
func (parser *jsonParser) Marshal(conf interface{}, toFile string) (err error) {
	err = YamlMarshalAndWriteFile(conf, toFile)
	return
}

// JSONReadFileAndUnmarshal 集成读取文件和反序列化操作
func JSONReadFileAndUnmarshal(conf interface{}, fromFile string) error {
	if _, err := os.Stat(fromFile); os.IsNotExist(err) {
		YamlMarshalAndWriteFile(conf, fromFile)
		return err
	}

	data, err := ioutil.ReadFile(fromFile)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, conf); err != nil {
		return err
	}
	return nil
}

// JSONMarshalAndWriteFile 集成序列化和写文件操作
func JSONMarshalAndWriteFile(conf interface{}, filename string) error {
	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	// write the file
	if err = ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	return nil
}
