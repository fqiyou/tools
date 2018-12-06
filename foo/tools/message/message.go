package message

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/fqiyou/tools/foo/tools/logs"
	"github.com/fqiyou/tools/foo/util"
	"strconv"
	"sync"
)

type Message struct {
	MessageErrorInfo error `json:"message_error_info"`
	MessageHeadSize 	int64	`json:"message_head_size"`
	MessageInfo	struct{
		MsgHeadMap		map[string]interface{} 	`json:"msg_head_map"`
		MsgBodyMapList		[] map[string]interface{} 	`json:"msg_body_map"`

	} `json:"message_info"`
	MessageList  [] map[string]interface{} `json:"message_list"`
}




const (
	MESSAGE_HEAD_INFO_LENGTH = 8
	ERROR_INFO_MEG_LENGTH = "消息体长度小于8"
	ERROR_INFO_MSG_HEAD_LENGTH = "消息头长度不足"
)
var (
	wg    sync.WaitGroup
	mutex sync.Mutex
)


func (msg *Message) ToMessage(msg_string string){

	decode_message,err := base64.StdEncoding.DecodeString(msg_string)
	if err != nil {
		msg.MessageErrorInfo = err
	}
	if len(decode_message) <= MESSAGE_HEAD_INFO_LENGTH{
		msg.MessageErrorInfo = errors.New(ERROR_INFO_MEG_LENGTH);

	}
	msg.MessageHeadSize,err = strconv.ParseInt(string(decode_message[:MESSAGE_HEAD_INFO_LENGTH]), 10, 64)
	if err != nil {
		msg.MessageErrorInfo = err
	}
	if int64(len(decode_message)) <= (msg.MessageHeadSize + MESSAGE_HEAD_INFO_LENGTH) {
		msg.MessageErrorInfo = errors.New(ERROR_INFO_MSG_HEAD_LENGTH);
	}

	mes_head_string := string(decode_message[MESSAGE_HEAD_INFO_LENGTH:msg.MessageHeadSize+MESSAGE_HEAD_INFO_LENGTH])
	msg_body_string := ""
	message_body_base64,err := base64.StdEncoding.DecodeString(string(decode_message[msg.MessageHeadSize+MESSAGE_HEAD_INFO_LENGTH:]))
	if err != nil {
		msg.MessageErrorInfo = err
	}
	message_body_base64_gz,err := GzipDecode(message_body_base64)
	if err != nil {
		log.Error("非gzip格式",err)
		msg_body_string = string(message_body_base64)
	}else {
		msg_body_string = string(message_body_base64_gz)
	}
	if err := json.Unmarshal([]byte(mes_head_string), &msg.MessageInfo.MsgHeadMap); err != nil {
		log.Error(err)
		msg.MessageErrorInfo = err
	}
	message_body_byte_list := []byte(msg_body_string)
	//
	//if  !strings.HasPrefix(msg_body_string, "["){
	//	message_body_byte_list =[]byte("["+ msg_body_string+"]")
	//}

	if message_body_byte_list[0] == 123{
		message_body_byte_list =[]byte("["+ msg_body_string+"]")
	}

	if err := json.Unmarshal(message_body_byte_list, &msg.MessageInfo.MsgBodyMapList); err != nil {
		log.Error(err)
		msg.MessageErrorInfo = err
	}

}



func (msg *Message) ToMessageList(msg_string string){
	defer func() {
		if err:=recover();err!=nil{
			log.Error(err)
			msg.MessageErrorInfo = errors.New("未知异常")
		}
	}()
	msg.ToMessage(msg_string)
	if msg.MessageErrorInfo != nil{
		log.Error(msg.MessageErrorInfo)
		return
	}

	for index , _ := range msg.MessageInfo.MsgBodyMapList {
		wg.Add(1)
		go func() {
			msg_all_info := make(map[string]interface{})
			for key, value := range msg.MessageInfo.MsgHeadMap {
				msg_all_info[key] = value
			}
			msg_all_info["content"] = msg.MessageInfo.MsgBodyMapList[index]
			msg_all_info["$event"] = msg.MessageInfo.MsgBodyMapList[index]["event"]
			msg_all_info["uuid"] = util.NewUUID()
			mutex.Lock()
			msg.MessageList = append(msg.MessageList,msg_all_info)
			mutex.Unlock()
			wg.Done()
		}()
		wg.Wait()
	}
}