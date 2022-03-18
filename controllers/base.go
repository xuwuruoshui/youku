package controllers

import (
	"github.com/beego/beego/v2/server/web/context"
	"strconv"
	"youku/models"
)

// 获取频道地区列表
func ChannelRegion(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")

	channelId, _ := strconv.Atoi(channelIdStr)
	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}

	nums, regions, err := models.GetChannelRegion(channelId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", regions, nums))
}

// 获取评到类型列表
func ChannelType(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")

	channelId, _ := strconv.Atoi(channelIdStr)
	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定类型"))
		return
	}

	nums, regions, err := models.GetChannelType(channelId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", regions, nums))
}
