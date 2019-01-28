package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func GetYamlConfig(file_path string,conf interface{})   (interface{}, error){
	yamlFile, err := ioutil.ReadFile(file_path)

	if err != nil {
		Log.Error(err)
		return conf, err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		Log.Error(err)
		return conf, err
	}
	return conf, nil
}
