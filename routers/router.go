package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"youku/controllers"
)

func init() {
	userRouter()
	channelRouter()
	videoRouter()
	commentRouter()
	topRouter()
	BarrageRouter()
	AliyunRouter()
}

// 用户相关
func userRouter() {
	beego.Post("/register/save", controllers.SaveRegister)
	beego.Post("/login/do", controllers.LoginDo)
	beego.Post("/send/message", controllers.SendMessageDo)
}

// 频道相关
func channelRouter() {
	beego.Get("/channel/advert", controllers.ChannelAdvert)
	beego.Get("/channel/hot", controllers.ChannelHotList)
	beego.Get("/channel/recommend/region", controllers.ChannelRecommendRegion)
	beego.Get("/channel/recommend/type", controllers.ChannelRecommendTypeList)
	beego.Get("/channel/region", controllers.ChannelRegion)
	beego.Get("/channel/type", controllers.ChannelType)
	beego.Get("/channel/video", controllers.ChannelVideo)
}

// 视频获取
func videoRouter() {
	beego.Get("/video/info", controllers.VideoInfo)
	beego.Get("/video/episodes/list", controllers.VideoEpisodesList)
	beego.Get("/user/video", controllers.UserVideo)
	beego.Post("/video/save", controllers.VideoSave)

}

// 评论
func commentRouter() {
	beego.Get("/comment/list", controllers.List)
	beego.Post("/comment/save", controllers.Save)
}

// 排行榜
func topRouter() {
	beego.Get("/channel/top", controllers.ChannelTop)
	beego.Get("/type/top", controllers.TypeTop)
}

// 弹幕
func BarrageRouter() {
	beego.Get("/barrage/ws", controllers.BarrageWs)
	beego.Post("/barrage/save", controllers.BarrageSave)
}

// 阿里云
func AliyunRouter() {
	beego.Post("/aliyun/create/upload/video", controllers.CreatUploadVideo)
	beego.Post("/aliyun/refresh/upload/video", controllers.RefreshUploadVideo)
	beego.Post("/aliyun/video/play/auth", controllers.GetPlayAuth)
	beego.Post("/aliyun/video/callback", controllers.VideoCallback)
}
