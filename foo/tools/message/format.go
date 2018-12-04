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


// 将message字符串封装成message对象,解base64,gz
func FormatMessage(message string)  (Message, error){
	message_struts := new(Message)
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

func FormatMessageObject(message string) (message_struts Message, message_error error){

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