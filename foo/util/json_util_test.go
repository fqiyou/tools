package util

import (
	"fmt"
	"testing"
)

type User struct {
	Id 		int		`json:"id"`
	Name string		`json:"name"`
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