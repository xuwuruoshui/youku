package models

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"time"
	"youku/service/rabbitmq"
)

type Comment struct {
	Id          int
	Content     string
	AddTime     int64
	UserId      int
	Stamp       int
	Status      int
	PraiseCount int
	EpisodesId  int
	VideoId     int
}

func init() {
	orm.RegisterModel(new(Comment))
}

func GetCommentList(episodesId int, offset int, limit int) (int64, []Comment, error) {
	o := orm.NewOrm()
	var comments []Comment
	num, _ := o.Raw("SELECT id FROM comment WHERE status=1 AND episodes_id=?", episodesId).QueryRows(&comments)
	_, err := o.Raw("SELECT id,content,add_time,user_id,stamp,praise_count,episodes_id FROM comment WHERE status=1 AND episodes_id=? ORDER BY add_time DESC LIMIT ?,?", episodesId, offset, limit).QueryRows(&comments)
	return num, comments, err
}

func SaveComment(content string, uid, episodesId, videoId int) error {
	o := orm.NewOrm()
	var comment Comment
	comment.Content = content
	comment.UserId = uid
	comment.EpisodesId = episodesId
	comment.VideoId = videoId
	comment.Stamp = 0
	comment.Status = 1
	comment.AddTime = time.Now().Unix()
	_, err := o.Insert(&comment)
	if err == nil {
		// 修改视频的总评论数
		o.Raw("UPDATE video SET comment=comment+1 WHERE id=?", videoId).Exec()
		// 修改视频剧集的评论数
		o.Raw("UPDATE video_episodes SET comment=comment+1 WHERE id=?", episodesId).Exec()
		// 更新redis排行榜 - 通过MQ实现
		// 创建一个简单模式的MQ
		// 把要传递的数据转换为json字符串
		videoObj := map[string]int{
			"VideoId": videoId,
		}
		videoJson, _ := json.Marshal(videoObj)
		rabbitmq.Publish("", "youku_top", string(videoJson))
		
		// 延时添加浏览数(没有浏览字段，就将就使用评论数了,10秒后再添加一次评论数)
		videoCountObj := map[string]int{
			"VideoId": videoId,
			"EpisodesId": episodesId,
		}
		videoCountJson, _ := json.Marshal(videoCountObj)
		rabbitmq.PublicDlx("youku.comment.count",string(videoCountJson))
	}
	return err
}
