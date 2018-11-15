package util

import (
	"encoding/json"
	"fmt"
)

func JsonPrint(model interface{})  {
	ba,_ := json.Marshal(model)
	fmt.Println(string(ba))
}

func ModelToString(model interface{}) string {
	ba,_ := json.Marshal(model)
	return string(ba)
}