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
	
	// 获取uid channel
	uidChan := make(chan int,limit)
	closeChane := make(chan bool,5)
	resChan := make(chan models.UserInfo,limit)
	// 把获取到的uid放到channel中
	go func() {
		for _, v := range comments {
			uidChan <- v.UserId
		}
		close(uidChan)
	}()
	// 开5个goroutine处理uidChannel中的信息
	for i:=0;i<5;i++ {
		
		go channelGetUserInfo(uidChan,resChan,closeChane)
	}
	// 判断是否执行完成,组合信息
	go func() {
		for i:=0;i<5;i++{
			<-closeChane
		}
		close(resChan)
		close(closeChane)
	}()
	
	userInfoMap := make(map[int]models.UserInfo)
	for r := range resChan {
		userInfoMap[r.Id] = r
	}
	
	for _,v := range comments{
			commentInfo.Id = v.Id
			commentInfo.Content = v.Content
			commentInfo.AddTime = v.AddTime
			commentInfo.AddTimeTitle = DateFormat(v.AddTime)
			commentInfo.UserId = v.UserId
			commentInfo.Stamp = v.Stamp
			commentInfo.PraiseCount = v.PraiseCount
			commentInfo.UserInfo = userInfoMap[v.UserId]
			data = append(data, commentInfo)
	}

	ctx.JSONResp(ReturnSuccess(0, "success", data, num))
}

func channelGetUserInfo(uidChan chan int,resChan chan models.UserInfo,closeChan chan bool){
	for uid:= range uidChan{
		res, err := models.RedisGetUserInfo(uid)
		if err == nil{
			resChan <- res
		}
	}
	closeChan<-true
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
