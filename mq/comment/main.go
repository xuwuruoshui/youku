package main

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"youku/service/rabbitmq"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	db, _ := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql",db)
	
	
	rabbitmq.ComsumerDlx("youku.comment.count","youku_comment_count","youku.comment.count.dlx","youku_comment_count_dlx",10000,callback)
	
}


func callback(s string) (err error){
	type Data struct {
		VideoId int
		EpisodesId int
	}
	
	var data Data
	err = json.Unmarshal([]byte(s),&data)
	if err!=nil{
		fmt.Println(err)
		return 
	}
	o := orm.NewOrm()
	// 修改视频的总评论数
	o.Raw("UPDATE video SET comment=comment+1 WHERE id=?", data.VideoId).Exec()
	// 修改视频剧集的评论数
	o.Raw("UPDATE video_episodes SET comment=comment+1 WHERE id=?", data.EpisodesId).Exec()
	// 更新redis排行榜 - 通过MQ实现
	// 创建一个简单模式的MQ
	// 把要传递的数据转换为json字符串
	videoObj := map[string]int{
		"VideoId": data.VideoId,
	}
	videoJson, _ := json.Marshal(videoObj)
	rabbitmq.Publish("", "youku_top", string(videoJson))

	fmt.Println("msg is: ", s)
	return 
}