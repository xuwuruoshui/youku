package controllers

import (
	context "github.com/beego/beego/v2/server/web/context"
	"regexp"
	"strconv"
	"strings"
	"youku/models"
)

func SaveRegister(ctx *context.Context) {
	mobile := ctx.Input.Query("mobile")
	password := ctx.Input.Query("password")

	if mobile == "" {
		ctx.JSONResp(ReturnError(4001, "手机号不能为空"))
		return
	}
	isValid, _ := regexp.MatchString(`^1(3|4|5|7|8)[0-9]\d{8}$`, mobile)
	if !isValid {
		ctx.JSONResp(ReturnError(4002, "手机格式不正确"))
		return
	}

	if password == "" {
		ctx.JSONResp(ReturnError(4003, "密码不能为空"))
		return
	}
	// 判断手机号是否注册
	status := models.IsUserMobile(mobile)
	if status {
		ctx.JSONResp(ReturnError(4005, "手机号已经注册"))
		return
	}
	err := models.UserSave(mobile, MD5V(password))
	if err != nil {
		ctx.JSONResp(ReturnError(5000, err))
		return

	}
	ctx.JSONResp(ReturnSuccess(0, "注册成功", nil, 0))
}

func LoginDo(ctx *context.Context) {
	mobile := ctx.Input.Query("mobile")
	password := ctx.Input.Query("password")
	if mobile == "" {
		ctx.JSONResp(ReturnError(4001, "手机号不能为空"))
		return
	}
	isValid, _ := regexp.MatchString(`^1(3|4|5|7|8)[0-9]\d{8}$`, mobile)
	if !isValid {
		ctx.JSONResp(ReturnError(4002, "手机格式不正确"))
		return
	}

	if password == "" {
		ctx.JSONResp(ReturnError(4003, "密码不能为空"))
		return
	}

	uid, name := models.IsMobileLogin(mobile, MD5V(password))
	if uid == 0 {
		ctx.JSONResp(ReturnError(4004, "手机号或密码不正确"))
		return
	}

	ctx.JSONResp(ReturnSuccess(0, "登录成功", map[string]interface{}{"uid": uid, "username": name}, 1))
}

func SendMessageDo(ctx *context.Context) {
	uids := ctx.Input.Query("uids")
	content := ctx.Input.Query("content")

	if uids == "" {
		ctx.JSONResp(ReturnError(4001, "请填写接收人"))
		return
	}

	if content == "" {
		ctx.JSONResp(ReturnError(4002, "请填写发送内容"))
		return
	}

	messageId, err := models.SendMessageDo(content)
	if err != nil {
		ctx.JSONResp(ReturnError(5000, "发送失败请联系客服"))
		return
	}
	uidConfig := strings.Split(uids, ",")
	for _, v := range uidConfig {
		userId, _ := strconv.Atoi(v)
		models.SendMessageUser(userId, messageId)
	}
	ctx.JSONResp(ReturnSuccess(0, "发送成功", "", 1))
}
