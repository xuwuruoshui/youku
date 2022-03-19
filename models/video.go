package models

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
	redisClient "youku/service/redis"
)

type Video struct {
	Id                 int
	Title              string
	SubTitle           string
	AddTime            int64
	Img                string
	Img1               string
	EpisodesCount      int
	IsEnd              int
	ChannelId          int
	Status             int
	RegionId           int
	TypeId             int
	EpisodesUpdateTime int64
	Comment            int
	UserId             int
	IsRecommend        int
}
type VideoData struct {
	Id            int
	Title         string
	SubTitle      string
	AddTime       int64
	Img           string
	Img1          string
	EpisodesCount int
	IsEnd         int
	Comment       int
}
type Episodes struct {
	Id            int
	Title         string
	AddTime       int64
	Num           int
	PlayUrl       string
	Comment       int
	AliyunVideoId string
}

func init() {
	orm.RegisterModel(new(VideoData))
	orm.RegisterModel(new(Video))
	orm.RegisterModel(new(Episodes))
}

// 视频基本信息
func GetChannelHostList(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var video []VideoData
	num, err := o.Raw("select id,title,sub_title,add_time,img,img1,episodes_count,is_end from video where "+
		"status=1 and is_hot=1 and channel_id=? order by episodes_update_time desc limit 9", channelId).QueryRows(&video)
	return num, video, err
}

func GetChannelRecommendRegionList(channelId, regionId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var video []VideoData
	num, err := o.Raw("select id,title,sub_title,add_time,img,img1,episodes_count,is_end from video where "+
		"status=1 and is_recommend=1 and region_id=? and channel_id=? order by episodes_update_time desc limit 9", regionId, channelId).QueryRows(&video)
	return num, video, err
}

func GetChannelRecommendTypeList(channelId, typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var video []VideoData
	num, err := o.Raw("select id,title,sub_title,add_time,img,img1,episodes_count,is_end from video where "+
		"status=1 and is_recommend=1 and type_id=? and channel_id=? order by episodes_update_time desc limit 9", typeId, channelId).QueryRows(&video)
	return num, video, err
}

func GetChannelVideoList(channelId, regionId, typeId, offset, limit int, end, sort string) (int64, []orm.Params, error) {
	o := orm.NewOrm()
	var videos []orm.Params
	qs := o.QueryTable("video")
	qs = qs.Filter("channel_id", channelId)
	qs = qs.Filter("status", 1)
	if regionId > 0 {
		qs = qs.Filter("region_id", regionId)
	}
	if typeId > 0 {
		qs = qs.Filter("type_id", typeId)
	}
	if end == "n" {
		qs = qs.Filter("is_end", 0)
	} else {
		qs = qs.Filter("is_end", 1)
	}

	if sort == "episodesUpdateTime" {
		qs = qs.OrderBy("-episodes_update_time")
	} else if sort == "comment" {
		qs = qs.OrderBy("-comment")
	} else {
		qs = qs.OrderBy("-add_time")
	}
	nums, _ := qs.Count()
	qs = qs.Limit(limit, offset)
	_, err := qs.Values(&videos, "id", "title", "sub_title", "add_time", "img", "img1",
		"episodes_count", "is_end")
	return nums, videos, err
}

// 视频信息episodes
func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.Raw("select * from video where id =? limit 1", videoId).QueryRow(&video)
	return video, err
}

// 添加Redis缓存，获取视频详情
func RedisGetVideoInfo(videoId int) (Video, error) {
	var video Video
	conn := redisClient.PoolConnect()
	defer conn.Close()

	// 定义redis key
	redisKey := "video:id:" + strconv.Itoa(videoId)
	// 判断redis中是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if !exists {
		o := orm.NewOrm()
		o.Raw("select * from video where id =? limit 1", videoId).QueryRow(&video)
		_, err = conn.Do("hmset", redis.Args{}.Add(redisKey).AddFlat(video)...)
		if err == nil {
			_, err = conn.Do("expire", redisKey, 86400)
		}
	} else {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &video)
	}
	return video, err
}

func GetVideoEpisodesList(videoId int) (int64, []Episodes, error) {
	o := orm.NewOrm()
	var episodes []Episodes
	num, err := o.Raw("select id,title,add_time,num,play_url,comment,aliyun_video_id from video_episodes where video_id=? "+
		"and status =1 order by num asc", videoId).QueryRows(&episodes)
	return num, episodes, err
}

func RedisGetVideoEpisodesList(videoId int) (int64, []Episodes, error) {
	var (
		episodes []Episodes
		num      int64
		err      error
	)

	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:episodes:videoId:" + strconv.Itoa(videoId)
	exist, _ := redis.Bool(conn.Do("exists", redisKey))

	if !exist {
		o := orm.NewOrm()
		num, err = o.Raw("select id,title,add_time,num,play_url,comment,aliyun_video_id from video_episodes where video_id=? "+
			"and status =1 order by num asc", videoId).QueryRows(&episodes)
		for _, v := range episodes {
			jsonValue, err := json.Marshal(v)
			if err == nil {
				conn.Do("rpush", redisKey, jsonValue)
			}
		}
		conn.Do("expire", redisKey, 86400)
	} else {
		num, err = redis.Int64(conn.Do("llen", redisKey))
		if err != nil {
			return 0, nil, err
		}
		values, err := redis.Values(conn.Do("lrange", redisKey, "0", "-1"))
		if err != nil {
			return 0, nil, err
		}
		var episodesInfo Episodes
		for _, v := range values {
			err := json.Unmarshal(v.([]byte), &episodesInfo)
			if err != nil {
				return 0, nil, err
			}
			episodes = append(episodes, episodesInfo)
		}
	}

	return num, episodes, err
}

// 频道排行榜
func GetChannelTop(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE status=1 AND channel_id=? ORDER BY comment DESC LIMIT 10", channelId).QueryRows(&videos)
	return num, videos, err
}

// redis频道排行榜
func RedisGetChannelTop(channelId int) (int64, []VideoData, error) {
	conn := redisClient.PoolConnect()
	defer conn.Close()

	var (
		num    int64
		videos []VideoData
		err    error
	)

	redisKey := "video:top:channel:channelId:" + strconv.Itoa(channelId)
	exists, err := redis.Bool(conn.Do("exists", redisKey))

	if !exists {
		o := orm.NewOrm()
		num, err = o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end,comment FROM video WHERE status=1 AND channel_id=? ORDER BY comment DESC LIMIT 10", channelId).QueryRows(&videos)
		for _, v := range videos {
			conn.Do("zadd", redisKey, v.Comment, v.Id)
		}
		conn.Do("expire", redisKey, 86400*30)
	} else {
		res, _ := redis.Values(conn.Do("zrevrange", redisKey, "0", "10", "WITHSCORES"))
		num := 0
		for k, v := range res {
			if k&1 == 0 {
				videoId, _ := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.IsEnd = videoInfo.IsEnd
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.Comment = videoInfo.Comment
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}

	}

	return num, videos, err
}

//类型排行榜
func GetTypeTop(typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE status=1 AND type_id=? ORDER BY comment DESC LIMIT 10", typeId).QueryRows(&videos)
	return num, videos, err
}

// redis频类型排行榜
func RedisGetTypeTop(typeId int) (int64, []VideoData, error) {
	conn := redisClient.PoolConnect()
	defer conn.Close()

	var (
		num    int64
		videos []VideoData
		err    error
	)

	redisKey := "video:top:type:typeId:" + strconv.Itoa(typeId)
	exists, err := redis.Bool(conn.Do("exists", redisKey))

	if !exists {
		o := orm.NewOrm()
		num, err = o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE status=1 AND type_id=? ORDER BY comment DESC LIMIT 10", typeId).QueryRows(&videos)
		for _, v := range videos {
			conn.Do("zadd", redisKey, v.Comment, v.Id)
		}
		conn.Do("expire", redisKey, 86400*30)
	} else {
		res, _ := redis.Values(conn.Do("zrevrange", redisKey, "0", "10", "WITHSCORES"))
		num := 0
		for k, v := range res {
			if k&1 == 0 {
				videoId, _ := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.IsEnd = videoInfo.IsEnd
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.Comment = videoInfo.Comment
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}

	}

	return num, videos, err
}

func GetUserVideo(uid int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count, is_end FROM video WHERE user_id=? ORDER BY add_time DESC", uid).QueryRows(&videos)
	return num, videos, err
}

func SaveVideo(title string, subTitle string, channelId int, regionId int, typeId int, playUrl string, user_id int, aliyunVideoId string) error {
	o := orm.NewOrm()
	var video Video
	time := time.Now().Unix()
	video.Title = title
	video.SubTitle = subTitle
	video.AddTime = time
	video.Img = ""
	video.Img1 = ""
	video.EpisodesCount = 1
	video.IsEnd = 1
	video.ChannelId = channelId
	video.Status = 1
	video.RegionId = regionId
	video.TypeId = typeId
	video.EpisodesUpdateTime = time
	video.Comment = 0
	video.UserId = user_id
	videoId, err := o.Insert(&video)
	if err == nil {
		if aliyunVideoId != "" {
			playUrl = ""
		}
		_, err = o.Raw("INSERT INTO video_episodes (title,add_time,num,video_id,play_url,status,comment,aliyun_video_id) VALUES (?,?,?,?,?,?,?,?)", subTitle, time, 1, videoId, playUrl, 1, 0, aliyunVideoId).Exec()
		//fmt.Println(err)
	}
	return err
}

func SaveAliyunVideo(videoId, log string) error {
	o := orm.NewOrm()
	_, err := o.Raw("INSERT INTO aliyun_video (video_id, log, add_time) VALUES (?,?,?)", videoId, log, time.Now().Unix()).Exec()

	return err
}
