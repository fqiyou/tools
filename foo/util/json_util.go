package util

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

func JsonPrint(model interface{})  {
	ba,_ := json.Marshal(model)
	fmt.Println(string(ba))
}

func ModelToString(model interface{}) string {
	ba,_ := json.Marshal(model)
	return string(ba)
}


func MapToStruts(map_object interface{},struct_object interface{}) error {
	decConfig := mapstructure.DecoderConfig{TagName:"json",Result:&struct_object,WeaklyTypedInput:true,}

	dec,err := mapstructure.NewDecoder(&decConfig)
	if err != nil{
		return err
	}
	if err := dec.Decode(map_object); err != nil {

		return err
	}
	//
	//if err := mapstructure.Decode(map_object, &struct_object); err != nil {
	//	return err
	//}
	return nil
}
