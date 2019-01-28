package util

import (
	"testing"
)

type MysqlDbConf struct {
	Name 	string		`yaml:"name" json:"name"`
	Type	string		`yaml:"type" json:"type"`
	User	string		`yaml:"user" json:"user"`
	PassWord	string	`yaml:"password" json:"pass_word"`
	Host	string	`yaml:"host" json:"host"`
	Port 	string	`yaml:"port" json:"port"`
	Db		string   `yaml:"db" json:"db"`
}

type MysqlDbListConf struct {
	DbList[] MysqlDbConf `yaml:"db_list" json:"db_list"`
}

type YamlDbConf struct {
	Mysql MysqlDbListConf `yaml:"mysql" json:"mysql"`
}


func TestGetYamlConfig(t *testing.T)  {
	conf := new(YamlDbConf)
	_, err := GetYamlConfig("/Users/chaoyang/GoProject/src/github.com/fqiyou/tools/foo/conf/db_conf.yaml",&conf)
	if err != nil {
		Log.Fatal(err)
		return
	}
	JsonPrint(conf)
}