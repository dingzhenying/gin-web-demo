package utils

import (
	"encoding/json"
	"reflect"
)

/**
 * 接口消息处理
 */
type Message struct {
	Success bool
	Code    int
	Msg     string
	Data    interface{}
}

func GetSuccessMsg(obj interface{}) map[string]interface{} {
	var message = Message{true, 0, "成功!", obj}
	t := reflect.TypeOf(message)
	v := reflect.ValueOf(message)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
	//info, _ := json.Marshal(message)
	//return string(info)

}
func getErrMsg(msg string, data string) string {
	var message = Message{false, -1, msg, data}
	info, _ := json.Marshal(message)
	return string(info)
}
