package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web/context"
	"strconv"
	"youku/models"
	es "youku/service/elasticsearch"
)

// 顶部广告
func ChannelAdvert(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")

	channelId, _ := strconv.Atoi(channelIdStr)
	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}

	num, videos, err := models.GetChannelAdvert(channelId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "请求数据失败，请稍后重试~"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", videos, num))
}

// 热播
func ChannelHotList(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")

	channelId, _ := strconv.Atoi(channelIdStr)
	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}
	numb, videos, err := models.GetChannelHostList(channelId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", videos, numb))
}

// 频道页-根据频道地区获取推荐视频
func ChannelRecommendRegion(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")
	regionIdStr := ctx.Input.Query("regionId")

	channelId, _ := strconv.Atoi(channelIdStr)
	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}
	regionId, _ := strconv.Atoi(regionIdStr)
	if regionId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道地区"))
		return
	}
	nums, videos, err := models.GetChannelRecommendRegionList(channelId, regionId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", videos, nums))
}

//频道页-根据频道类型获取推荐视频
func ChannelRecommendTypeList(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")
	typeIdStr := ctx.Input.Query("typeId")

	channelId, _ := strconv.Atoi(channelIdStr)
	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}
	typeId, _ := strconv.Atoi(typeIdStr)
	if typeId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道类型"))
		return
	}
	nums, videos, err := models.GetChannelRecommendTypeList(channelId, typeId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", videos, nums))
}

// 获取视频列表
func ChannelVideo(ctx *context.Context) {
	channelIdStr := ctx.Input.Query("channelId")
	regionIdStr := ctx.Input.Query("regionId")
	typeIdStr := ctx.Input.Query("typeId")
	limitStr := ctx.Input.Query("limit")
	offsetStr := ctx.Input.Query("offset")
	end := ctx.Input.Query("end")
	sort := ctx.Input.Query("sort")

	channelId, _ := strconv.Atoi(channelIdStr)
	regionId, _ := strconv.Atoi(regionIdStr)
	typeId, _ := strconv.Atoi(typeIdStr)
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if channelId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道"))
		return
	}
	if limit == 0 {
		// 默认12条
		limit = 6
	}
	nums, videos, err := models.GetChannelVideoListEs(channelId, regionId, typeId, offset, limit, end, sort)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", videos, nums))
}

// 获取视频详情
func VideoInfo(ctx *context.Context) {
	videoIdStr := ctx.Input.Query("videoId")
	videoId, _ := strconv.Atoi(videoIdStr)
	if videoId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道id"))
	}
	video, err := models.RedisGetVideoInfo(videoId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "请求数据失败，请稍后重试"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", video, 1))
}

// 获取视频剧集列表
func VideoEpisodesList(ctx *context.Context) {
	videoIdStr := ctx.Input.Query("videoId")
	videoId, _ := strconv.Atoi(videoIdStr)
	if videoId == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定频道id"))
	}
	num, episodes, err := models.RedisGetVideoEpisodesList(videoId)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "请求数据失败，请稍后重试"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", episodes, num))
}

// 我的视频管理
func UserVideo(ctx *context.Context) {
	uidStr := ctx.Input.Query("uid")
	uid, _ := strconv.Atoi(uidStr)
	if uid == 0 {
		ctx.JSONResp(ReturnError(4001, "必须指定用户"))
		return
	}

	nums, video, err := models.GetUserVideo(uid)
	if err != nil {
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
	}
	ctx.JSONResp(ReturnSuccess(0, "success", video, nums))
}

// 保存用户上传视频信息
func VideoSave(ctx *context.Context) {
	playUrl := ctx.Input.Query("playUrl")
	title := ctx.Input.Query("title")
	subTitle := ctx.Input.Query("subTitle")
	channelIdStr := ctx.Input.Query("channelId")
	typeIdStr := ctx.Input.Query("typeId")
	regionIdStr := ctx.Input.Query("regionId")
	uidStr := ctx.Input.Query("uid")
	aliyunVideoId := ctx.Input.Query("aliyunVideoId")

	channelId, _ := strconv.Atoi(channelIdStr)
	typeId, _ := strconv.Atoi(typeIdStr)
	regionId, _ := strconv.Atoi(regionIdStr)
	uid, _ := strconv.Atoi(uidStr)

	if uid == 0 {
		ctx.JSONResp(ReturnError(4001, "请先登录"))
		return
	}
	if playUrl == "" {
		ctx.JSONResp(ReturnError(4002, "视频地址不能为空"))
		return
	}
	err := models.SaveVideo(title, subTitle, channelId, regionId, typeId, playUrl, uid, aliyunVideoId)
	if err != nil {
		ctx.JSONResp(ReturnError(5000, err))
		return
	}
	ctx.JSONResp(ReturnSuccess(0, "success", nil, 1))
}

// 导入es脚本
func SendEs(ctx *context.Context) {

	_, data, _ := models.GetAllList()
	for _, v := range data {
		body := es.E{
			"id":                   v.Id,
			"title":                v.Title,
			"sub_title":            v.SubTitle,
			"add_time":             v.AddTime,
			"img":                  v.Img,
			"img1":                 v.Img1,
			"episodes_count":       v.EpisodesCount,
			"is_end":               v.IsEnd,
			"channel_id":           v.ChannelId,
			"status":               v.Status,
			"region_id":            v.RegionId,
			"type_id":              v.TypeId,
			"episodes_update_time": v.EpisodesUpdateTime,
			"comment":              v.Comment,
			"user_id":              v.UserId,
			"is_recommend":         v.IsRecommend,
		}
		es.Add("youku_video", "video-"+strconv.Itoa(v.Id), body)
	}
}

// 搜索
func Search(ctx *context.Context) {
	keyword := ctx.Input.Query("keyword")
	limitStr := ctx.Input.Query("limit")
	offsetStr := ctx.Input.Query("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if keyword == "" {
		ctx.JSONResp(ReturnError(4001, "关键字不能为空"))
		return
	}

	if limit == 0 {
		limit = 12
	}

	sort := []map[string]string{
		map[string]string{"id": "desc"},
	}

	query := es.E{
		"bool": es.E{
			"must": es.E{
				"match": es.E{
					"title": keyword,
				},
			},
		},
	}

	res, err := es.Search("youku_video", query, offset, limit, sort)
	if err!=nil{
		ctx.JSONResp(ReturnError(4002, "查询错误"))
		return
	}
	total := res.Hits.Total.Value
	var data []models.Video
	for _, v := range res.Hits.Hits {
		var itemData models.Video
		err := json.Unmarshal(v.Source,&itemData)
		if err==nil{
			data = append(data,itemData)
		}
	}
	if total<=0{
		ctx.JSONResp(ReturnError(4004, "没有相关内容"))
		return
	}
	ctx.JSONResp(ReturnSuccess(0,"success",data,int64(total)))
}
