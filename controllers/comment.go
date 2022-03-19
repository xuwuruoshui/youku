package controllers

import (
	"github.com/beego/beego/v2/server/web/context"
	"strconv"
	"youku/models"
)

type CommentInfo struct {
	Id           int             `json:"id"`
	Content      string          `json:"content"`
	AddTime      int64           `json:"addTime"`
	AddTimeTitle string          `json:"addTimeTitle"`
	UserId       int             `json:"userId"`
	Stamp        int             `json:"stamp"`
	PraiseCount  int             `json:"praiseCount"`
	UserInfo     models.UserInfo `json:"userinfo"`
}

// 获取评论列表
func List(ctx *context.Context) {

	episodesIdStr := ctx.Input.Query("episodesId")
	limitStr := ctx.Input.Query("limit")
	offsetStr := ctx.Input.Query("offset")

	episodesId, _ := strconv.Atoi(episodesIdStr)
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if episodesId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定视频聚集"))
		return
	}
	if limit == 0 {
		limit = 12
	}

	num, comments, err := models.GetCommentList(episodesId, offset, limit)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	var data []CommentInfo
	var commentInfo CommentInfo
	for _, v := range comments {
		commentInfo.Id = v.Id
		commentInfo.Content = v.Content
		commentInfo.AddTime = v.AddTime
		commentInfo.AddTimeTitle = DateFormat(v.AddTime)
		commentInfo.UserId = v.UserId
		commentInfo.Stamp = v.Stamp
		commentInfo.PraiseCount = v.PraiseCount
		// 获取用户信息
		info, err := models.RedisGetUserInfo(v.UserId)
		if err != nil {
			ctx.JSONResp(ReturnError(4004, "没有相关内容"))
			return
		}
		commentInfo.UserInfo = info
		data = append(data, commentInfo)
	}

	ctx.JSONResp(ReturnSuccess(0, "success", data, num))
}

// 保存评论
func Save(ctx *context.Context) {
	content := ctx.Input.Query("content")
	uidStr := ctx.Input.Query("uid")
	episodesIdStr := ctx.Input.Query("episodesId")
	videoIdStr := ctx.Input.Query("videoId")

	uid, _ := strconv.Atoi(uidStr)
	episodesId, _ := strconv.Atoi(episodesIdStr)
	videoId, _ := strconv.Atoi(videoIdStr)

	if content == "" {
		ctx.JSONResp(ReturnError(4001, "内容不能为空"))
		return
	}
	if uid == 0 {
		ctx.JSONResp(ReturnError(4002, "请先登录"))
		return
	}
	if episodesId == 0 {
		ctx.JSONResp(ReturnError(4003, "必须指定剧集ID"))
		return
	}
	if videoId == 0 {
		ctx.JSONResp(ReturnError(4005, "必须指定视频ID"))
		return
	}

	err := models.SaveComment(content, uid, episodesId, videoId)
	if err != nil {
		ctx.JSONResp(ReturnError(5000, err))
	}
	ctx.JSONResp(ReturnSuccess(0, "success", "", 1))
}
