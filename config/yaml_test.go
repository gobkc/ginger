package config

import (
	"log"
	"os"
	"testing"
)

type MutiConf struct {
	Driver Db `yaml:"db"`
}
type Db struct {
	DBinfo Mysql `yaml:"mysql"`
}
type Mysql struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Info     Info   `yaml:"info"`
}
type Info struct {
	ID      int    `yaml:"id"`
	Postion string `yaml:"postion"`
}

func TestParser(t *testing.T) {
	mutiConf := MutiConf{
		Driver: Db{
			DBinfo: Mysql{
				Name: "mysql",
				User: "admin",
				Port: 999,
				Info: Info{
					ID:      1111111,
					Postion: "china",
				},
			},
		},
	}
	// marshal data to a file
	//parser := config.NewYamlParser("test01.yaml")
	parser := NewConfParser("yaml")
	parser.Marshal(mutiConf, "test02.yaml")

	var conf MutiConf
	if err := parser.Unmarshal(&conf, "test02.yaml"); err != nil {
		log.Println("unmarshal error:", err)
	}
	log.Println("\n", conf)

	if conf != mutiConf {
		t.Error("parser deal not ok")
	} else {
		t.Log("parse deal ok")
	}
}

func TestYamlFunc(t *testing.T) {
	mutiConf := MutiConf{
		Driver: Db{
			DBinfo: Mysql{
				Name: "mysql",
				User: "admin",
				Port: 999,
				Info: Info{
					ID:      1111111,
					Postion: "china",
				},
			},
		},
	}

	// 函数方式
	YamlMarshalAndWriteFile(mutiConf, "testFun02.yaml")

	var conf MutiConf
	YamlReadFileAndUnmarshal(&conf, "testFun02.yaml")
	log.Println("\n", conf)

	if conf != mutiConf {
		t.Error("parser deal not ok")
	} else {
		t.Log("parse deal ok")
	}
}

func TestConfileNotExists(t *testing.T) {
	notExistsFile := "notExists.yml"
	os.Remove(notExistsFile)
	var conf MutiConf
	if err := JSONReadFileAndUnmarshal(&conf, notExistsFile); !os.IsNotExist(err) {
		t.Error("it should give a file not exists error")
	}
}
