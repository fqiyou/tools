package message

import "errors"

type Message struct {
	HeaderSize int64 	`json:"header_size"`
	HeaderStr	string 	`json:"header_str"`
	BodyStr		string 	`json:"body_str"`
}

var ERRORSLICEOUTOFRANGE = errors.New("slice bounds out of range");

