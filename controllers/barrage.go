package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"youku/models"
)

type WsData struct {
	CurrentTime int `json:"currentTime"`
	EpisodesId  int `json:"episodesId"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 获取弹幕
func BarrageWs(ctx *context.Context) {
	var (
		conn     *websocket.Conn
		err      error
		data     []byte
		barrages []models.BarrageData
	)
	if conn, err = upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil); err != nil {
		goto ERR
	}

	for {
		if _, data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		var wsData WsData
		err := json.Unmarshal([]byte(data), &wsData)
		if err != nil {
			goto ERR
		}

		endTime := wsData.CurrentTime + 60*1000
		// 获取弹幕数据
		_, barrages, err = models.BarrageList(wsData.EpisodesId, wsData.CurrentTime, endTime)
		if err == nil {
			if err := conn.WriteJSON(barrages); err != nil {
				goto ERR
			}
		}
	}

ERR:
	conn.Close()
}

// 发弹幕
func BarrageSave(ctx *context.Context) {
	uidStr := ctx.Input.Query("uid")
	content := ctx.Input.Query("content")
	currentTimeStr := ctx.Input.Query("currentTime")
	episodesIdStr := ctx.Input.Query("episodesId")
	videoIdStr := ctx.Input.Query("videoId")

	uid, _ := strconv.Atoi(uidStr)
	currentTime, _ := strconv.Atoi(currentTimeStr)
	episodesId, _ := strconv.Atoi(episodesIdStr)
	videoId, _ := strconv.Atoi(videoIdStr)

	if content == "" {
		ctx.JSONResp(ReturnError(4001, "弹幕不能为空"))
		return
	}
	if uid == 0 {
		ctx.JSONResp(ReturnError(4002, "请先登录"))
		return
	}
	if episodesId == 0 {
		ctx.JSONResp(ReturnError(4003, "必须指定聚集ID"))
		return
	}
	if videoId == 0 {
		ctx.JSONResp(ReturnError(4005, "必须指定视频ID"))
		return
	}
	if currentTime == 0 {
		ctx.JSONResp(ReturnError(4006, "必须指定视频播放时间"))
		return
	}
	err := models.SaveBarrage(episodesId, videoId, currentTime, uid, content)
	if err != nil {
		ctx.JSONResp(ReturnError(5000, err))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", "", 1))
}
