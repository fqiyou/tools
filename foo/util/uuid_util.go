package util

import (
	"github.com/satori/go.uuid"
)

func NewUUID()  (uuid_string string){
	u1,err := uuid.NewV4()
	if err != nil {
		return "00000000-0000-0000-0000-000000000000"
	}
	return u1.String()
}