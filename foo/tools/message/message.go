package message

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/fqiyou/tools/foo/util"
	"io/ioutil"
	"strconv"
	"sync"
)

type Message struct {
}

type MessageDecode struct {
	MessageErrorInfo error `json:"message_error_info"`
	MessageHeadSize 	int64	`json:"message_head_size"`
	MessageInfo	struct{
		MsgHeadMap		map[string]interface{} 	`json:"msg_head_map"`
		MsgBodyMapList		[] map[string]interface{} 	`json:"msg_body_map"`

	} `json:"message_info"`
}
type MessageList [] map[string]interface{}



func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

const (
	MESSAGE_HEAD_INFO_LENGTH = 8
	ERROR_INFO_MEG_LENGTH = "消息体长度小于8"
	ERROR_INFO_MSG_HEAD_LENGTH = "消息头长度不足"
)

func NewMessage() *Message {
	return &Message{}
}


func (msg *Message) ToMessage(msg_string string) MessageDecode{
	msg_decode := new(MessageDecode)
	decode_message,err := base64.StdEncoding.DecodeString(msg_string)
	if err != nil {
		msg_decode.MessageErrorInfo = err
	}
	if len(decode_message) <= MESSAGE_HEAD_INFO_LENGTH{
		msg_decode.MessageErrorInfo = errors.New(ERROR_INFO_MEG_LENGTH);

	}
	msg_decode.MessageHeadSize,err = strconv.ParseInt(string(decode_message[:MESSAGE_HEAD_INFO_LENGTH]), 10, 64)
	if err != nil {
		msg_decode.MessageErrorInfo = err
	}
	if int64(len(decode_message)) <= (msg_decode.MessageHeadSize + MESSAGE_HEAD_INFO_LENGTH) {
		msg_decode.MessageErrorInfo = errors.New(ERROR_INFO_MSG_HEAD_LENGTH);
	}

	mes_head_string := string(decode_message[MESSAGE_HEAD_INFO_LENGTH:msg_decode.MessageHeadSize+MESSAGE_HEAD_INFO_LENGTH])
	msg_body_string := ""
	message_body_base64,err := base64.StdEncoding.DecodeString(string(decode_message[msg_decode.MessageHeadSize+MESSAGE_HEAD_INFO_LENGTH:]))
	if err != nil {
		msg_decode.MessageErrorInfo = err
	}
	message_body_base64_gz,err := GzipDecode(message_body_base64)
	if err != nil {
		util.Log.Warn("非gzip格式",err)
		msg_body_string = string(message_body_base64)
	}else {
		msg_body_string = string(message_body_base64_gz)
	}
	if err := json.Unmarshal([]byte(mes_head_string), &msg_decode.MessageInfo.MsgHeadMap); err != nil {
		msg_decode.MessageErrorInfo = err
	}
	message_body_byte_list := []byte(msg_body_string)
	//
	//if  !strings.HasPrefix(msg_body_string, "["){
	//	message_body_byte_list =[]byte("["+ msg_body_string+"]")
	//}

	if message_body_byte_list[0] == 123{
		message_body_byte_list =[]byte("["+ msg_body_string+"]")
	}

	if err := json.Unmarshal(message_body_byte_list, &msg_decode.MessageInfo.MsgBodyMapList); err != nil {
		msg_decode.MessageErrorInfo = err
	}
	return *msg_decode
}



func (msg *Message) ToMessageList(msg_string string) MessageList{
	var msg_list MessageList
	var wg    sync.WaitGroup
	defer func() {
		if err:=recover();err!=nil{
			util.Log.Error(err)

		}
	}()
	msg_decode := msg.ToMessage(msg_string)
	if msg_decode.MessageErrorInfo != nil{
		util.Log.Error(msg_decode.MessageErrorInfo)
		return msg_list
	}
	ch := make(chan map[string]interface{})

	for index , _ := range msg_decode.MessageInfo.MsgBodyMapList {
		wg.Add(1)
		go func() {
			msg_all_info := make(map[string]interface{})
			for key, value := range msg_decode.MessageInfo.MsgHeadMap {
				msg_all_info[key] = value
			}
			msg_all_info["content"] = msg_decode.MessageInfo.MsgBodyMapList[index]
			msg_all_info["$event"] = msg_decode.MessageInfo.MsgBodyMapList[index]["event"]
			msg_all_info["$type"] = msg_decode.MessageInfo.MsgBodyMapList[index]["type"]
			msg_all_info["uuid"] = util.NewUUID()
			ch <- msg_all_info
			defer wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for v := range ch {
		msg_list = append(msg_list,v)
	}
	return msg_list
}