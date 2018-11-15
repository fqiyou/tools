package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func GetYamlConfig(file_path string,conf interface{})   (interface{}, error){
	yamlFile, err := ioutil.ReadFile(file_path)

	if err != nil {
		log.Fatalf("file_path:#%v ", err)
		return conf, err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return conf, err
	}
	return conf, nil
}
