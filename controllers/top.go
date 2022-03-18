package controllers

import (
	"github.com/beego/beego/v2/server/web/context"
	"strconv"
	"youku/models"
)

// 根据频道获取排行榜
func ChannelTop(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")
	channelId, _ := strconv.Atoi(channelIdStr)
	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}
	nums, videos, err := models.GetChannelTop(channelId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}

	ctx.JSONResp(ReturnSuccess(0, "success", videos, nums))
}

// 根据类型获取排行榜
func TypeTop(ctx *context.Context) {
	typeIdStr := ctx.Input.Query("typeId")
	typeId, _ := strconv.Atoi(typeIdStr)
	if typeId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}
	nums, videos, err := models.GetTypeTop(typeId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}

	ctx.JSONResp(ReturnSuccess(0, "success", videos, nums))
}
