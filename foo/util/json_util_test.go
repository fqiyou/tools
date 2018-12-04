package util

import (
	"fmt"
	"testing"
)

type User struct {
	Id 		int		`json:"id"`
	Name string		`json:"name"`
	LastName string 	`json:"last_name"`
}


func TestJsonPrint(t *testing.T) {
	user := User{Id:1,Name:"测试"}
	JsonPrint(user)
}

func TestModelToString(t *testing.T) {
	user := User{Id:1,Name:"测试"}
	a := ModelToString(user)
	fmt.Println(a)
}




func TestMapToStruts(t *testing.T){

	map_object := make(map[string]interface{})
	map_object["name"] = "yc"
	map_object["Id"] = 21
	map_object["last_name"] = "yang chao"

	user := new(User)
	err := MapToStruts(map_object,&user)
	fmt.Println(err)
	fmt.Println(user.LastName)
	JsonPrint(map_object)

}