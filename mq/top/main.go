package main

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"youku/models"
	"youku/service/rabbitmq"
	redisClient "youku/service/redis"
)

func main() {
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	db, _ := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql",db)
	
	
	rabbitmq.Consumer("","youku_top",callback)
	
}


func callback(s string){
	type Data struct {
		VideoId int
	}
	
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	videoInfo, err := models.RedisGetVideoInfo(data.VideoId)
	if err!=nil{
		fmt.Println(err)
		return
	}
	conn := redisClient.PoolConnect()
	defer conn.Close()
	// 更新排行榜
	redisChannelKey := "video:top:channel:channelId:"+strconv.Itoa(videoInfo.ChannelId)
	redisTypeKey := "video:top:type:typeId:"+strconv.Itoa(videoInfo.TypeId)
	conn.Do("zincrby",redisChannelKey,1,data.VideoId)
	conn.Do("zincrby",redisTypeKey,1,data.VideoId)
	
	fmt.Println("msg is: ", s)
}