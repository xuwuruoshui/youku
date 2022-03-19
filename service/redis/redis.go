package redisClient

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

func PoolConnect() redis.Conn {
	pool := &redis.Pool{ //实例化一个连接池
		MaxIdle: 5000, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   10000, //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300,   //连接关闭时间 300秒 （300秒不使用自动关闭）
		Wait:        true,  //超过最大连接数时，是等待还是报错
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			s, err := beego.AppConfig.String("redisdb")
			if err != nil {
				fmt.Println("redis初始化错误")
			}
			return redis.Dial("tcp", s)
		},
	}
	return pool.Get()
}
