package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
	redisClient "youku/service/redis"
)

type User struct {
	Id       int
	Name     string
	Password string
	Status   int
	AddTime  int64
	Mobile   string
	Avatar   string
}
type UserInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	AddTime int64  `json:"addTime"`
	Avatar  string `json:"avatar"`
}

func init() {
	orm.RegisterModel(new(User))
}

// 根据手机号判断用户是否存在
func IsUserMobile(mobile string) bool {
	o := orm.NewOrm()
	user := User{Mobile: mobile}
	err := o.Read(&user, "Mobile")
	if err == orm.ErrNoRows {
		return false
	} else if err == orm.ErrMissPK {
		return false
	}
	return true
}

// 保存用户
func UserSave(mobile, password string) error {
	o := orm.NewOrm()
	var user User
	user.Name = ""
	user.Password = password
	user.Mobile = mobile
	user.Status = 1
	user.AddTime = time.Now().Unix()
	_, err := o.Insert(&user)
	return err
}

// 登录功能
func IsMobileLogin(mobile, password string) (int, string) {
	o := orm.NewOrm()
	var user User
	err := o.QueryTable("user").
		Filter("mobile", mobile).
		Filter("password", password).One(&user)
	if err == orm.ErrNoRows {
		return 0, ""
	} else if err == orm.ErrMissPK {
		return 0, ""
	}
	return user.Id, user.Name
}

// 根据用户Id获取user
func GetUserInfo(userId int) (UserInfo, error) {
	o := orm.NewOrm()
	var u UserInfo
	err := o.Raw("select id,name,add_time,avatar from user where id = ? limit 1", userId).QueryRow(&u)
	return u, err
}

// 增加redis缓存 - 根据用户ID获取用户信息
func RedisGetUserInfo(uid int) (UserInfo, error) {
	var u UserInfo
	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "user:id:" + strconv.Itoa(uid)
	// 判断redis是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if !exists {
		o := orm.NewOrm()
		o.Raw("select id,name,add_time,avatar from user where id = ? limit 1", uid).QueryRow(&u)
		_, err := conn.Do("hmset", redis.Args{}.Add(redisKey).AddFlat(u)...)
		if err == nil {
			conn.Do("expire", redisKey, 86400)
		}
	} else {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &u)
	}
	return u, err
}
