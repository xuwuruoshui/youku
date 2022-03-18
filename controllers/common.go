package controllers

import (
	"crypto/md5"
	"encoding/hex"
	beego "github.com/beego/beego/v2/server/web"
	"time"
)

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(code int, msg interface{}, item interface{}, count int64) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg, Items: item, Count: count}
	return
}

func ReturnError(code int, msg interface{}) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg}
	return
}

// 加密
func MD5V(password string) string {
	hash := md5.New()
	key, _ := beego.AppConfig.String("md5code")
	hash.Write([]byte(password + key))
	return hex.EncodeToString(hash.Sum(nil))
}

// 格式化时间
func DateFormat(t int64) string {
	video_time := time.Unix(t, 0)
	return video_time.Format("2006-01-02")
}
