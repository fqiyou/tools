package message

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"github.com/fqiyou/tools/foo/tools/logs"
	"io/ioutil"
	"strconv"
)

type MessageStruct struct {
	HeaderSize int64 	`json:"header_size"`
	HeaderStr	string 	`json:"header_str"`
	BodyStr		string 	`json:"body_str"`
}


var ERRORSLICEOUTOFRANGE = errors.New("slice bounds out of range");

// 将message字符串封装成message对象,解base64,gz
func FormatMessage(message string)  (MessageStruct, error){
	message_struts := new(MessageStruct)

	decode_message,err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return *message_struts,err
	}
	if len(decode_message) <= 8{
		err = ERRORSLICEOUTOFRANGE
		log.Error(err)
		return *message_struts,err
	}
	message_struts.HeaderSize,err = strconv.ParseInt(string(decode_message[:8]), 10, 64)
	if err != nil {
		log.Error(err)
		return *message_struts,err
	}
	if int64(len(decode_message)) <= (message_struts.HeaderSize + 8) {
		err = ERRORSLICEOUTOFRANGE
		log.Error(err)
		return *message_struts,err
	}
	message_struts.HeaderStr = string(decode_message[8:message_struts.HeaderSize+8])
	message_body_base64,err := base64.StdEncoding.DecodeString(string(decode_message[message_struts.HeaderSize+8:]))
	if err != nil {
		log.Error(err)
		return *message_struts,err
	}
	message_struts.BodyStr = string(message_body_base64)
	message_body_base64_gz,err := GzipDecode(message_body_base64)
	if err != nil {
		log.Error(err)
		return *message_struts,err
	}
	message_struts.BodyStr = string(message_body_base64_gz)
	return *message_struts,nil
}

func FormatMessageObject(message string) (message_struts MessageStruct, message_error error){

	defer func() {
		if err:=recover();err!=nil{
			log.Error(err)
			message_error = errors.New("未知异常")
		}
	}()
	message_struts,message_error = FormatMessage(message)

	return message_struts,message_error

}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}