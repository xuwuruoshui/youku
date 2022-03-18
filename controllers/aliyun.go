package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vod"
	"github.com/beego/beego/v2/server/web/context"
	"youku/models"
)

var (
	accessKeyId     = "xxx"
	accessKeySecret = "xxx"
)

type JSONS struct {
	RequestId     string
	UploadAddress string
	UploadAuth    string
	VideoId       string
}
type PlayJSONS struct {
	PlayAuth string
}

// 获取音/视频上传地址和凭证
func CreatUploadVideo(ctx *context.Context) {

	title := ctx.Input.Query("title")
	desc := ctx.Input.Query("desc")
	fileName := ctx.Input.Query("fileName")
	coverUrl := ctx.Input.Query("coverUrl")
	tags := ctx.Input.Query("tags")

	client, err := InitVodClient(accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}

	response, err := MyCreateUploadVideo(client, title, desc, fileName, coverUrl, tags)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.GetHttpContentString())

	data := &JSONS{
		response.RequestId,
		response.UploadAddress,
		response.UploadAuth,
		response.VideoId,
	}
	ctx.JSONResp(data)
}

// 刷新音/视频上传凭证
func RefreshUploadVideo(ctx *context.Context) {

	videoId := ctx.Input.Query("videoId")

	client, err := InitVodClient(accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	response, err := MyRefreshUploadVideo(client, videoId)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpContentString())
	//fmt.Printf("UploadAddress: %s\n UploadAuth: %s", response.UploadAddress, response.UploadAuth)
	data := &JSONS{
		response.RequestId,
		response.UploadAddress,
		response.UploadAuth,
		response.VideoId,
	}
	ctx.JSONResp(data)
}

// 获取视频播放凭证
func GetPlayAuth(ctx *context.Context) {
	videoId := ctx.Input.Query("videoId")
	client, err := InitVodClient(accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	response, err := MyGetPlayAuth(client, videoId)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpContentString())
	//fmt.Printf("%s: %s\n", response.VideoMeta, response.PlayAuth)
	data := &PlayJSONS{
		response.PlayAuth,
	}
	ctx.JSONResp(data)
}

func InitVodClient(accessKeyId string, accessKeySecret string) (client *vod.Client, err error) {

	// 点播服务接入地域
	regionId := "cn-shanghai"

	// 创建授权对象
	credential := &credentials.AccessKeyCredential{
		accessKeyId,
		accessKeySecret,
	}

	// 自定义config
	config := sdk.NewConfig()
	config.AutoRetry = true     // 失败是否自动重试
	config.MaxRetryTime = 3     // 最大重试次数
	config.Timeout = 3000000000 // 连接超时，单位：纳秒；默认为3秒

	// 创建vodClient实例
	return vod.NewClientWithOptions(regionId, config, credential)
}

func MyCreateUploadVideo(client *vod.Client, title string, desc string, fileName string, coverUrl string, tags string) (response *vod.CreateUploadVideoResponse, err error) {
	request := vod.CreateCreateUploadVideoRequest()
	request.Title = title
	request.Description = desc
	request.FileName = fileName
	//request.CateId = "-1"
	request.CoverURL = coverUrl
	request.Tags = tags
	request.AcceptFormat = "JSON"
	return client.CreateUploadVideo(request)
}

func MyRefreshUploadVideo(client *vod.Client, videoId string) (response *vod.RefreshUploadVideoResponse, err error) {
	request := vod.CreateRefreshUploadVideoRequest()
	request.VideoId = videoId
	request.AcceptFormat = "JSON"

	return client.RefreshUploadVideo(request)
}

func MyGetPlayAuth(client *vod.Client, videoId string) (response *vod.GetVideoPlayAuthResponse, err error) {
	request := vod.CreateGetVideoPlayAuthRequest()
	request.VideoId = videoId
	request.AcceptFormat = "JSON"

	return client.GetVideoPlayAuth(request)
}

// 回调接口
type CallbackData struct {
	EventTime   string
	EventType   string
	VideoId     string
	Status      string
	Exteng      string
	StreamInfos []CallbackStreamInfosData
}
type CallbackStreamInfosData struct {
	Status     string
	Bitrate    int
	Definition string
	Duration   int
	Encrypt    bool
	FileUrl    string
	Format     string
	Fps        int
	Height     int
	Size       int
	Width      int
	JobId      string
}

func VideoCallback(ctx *context.Context) {
	var ob CallbackData
	r := ctx.Input.RequestBody
	json.Unmarshal(r, &ob)

	models.SaveAliyunVideo(ob.VideoId, string(r))
	ctx.WriteString("success")
}
