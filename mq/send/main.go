package main

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"youku/models"
	"youku/service/rabbitmq"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	db, _ := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql",db)


	rabbitmq.Consumer("","youku_send_message_user",callback)

}


func callback(s string) (err error){
	type Data struct {
		UserId int
		MessageId int64
	}
	var data Data
	err = json.Unmarshal([]byte(s), &data)
	if err!=nil{
		fmt.Println(err)
		return
	}

	err = models.SendMessageUser(data.UserId, data.MessageId)
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("msg is: ", s)
	return 
}