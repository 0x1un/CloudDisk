package util

import (
	"encoding/json"
	"log"
)

type RespJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewRespJson(code int, msg string, data interface{}) *RespJson {
	return &RespJson{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func (rp *RespJson) JsonBytes() []byte {
	if data, err := json.Marshal(rp); err != nil {
		log.Println(err.Error())
		return nil
	} else {
		return data
	}
}

func (rp *RespJson) JsonString() string {
	return string(rp.JsonBytes())
}
