package util

import (
	"github.com/satori/go.uuid"
)

func NewUUID()  (uuid_string string){
	return uuid.NewV4().String()
}