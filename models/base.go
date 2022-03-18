package models

import "github.com/beego/beego/v2/client/orm"

type Region struct {
	Id   int
	Name string
}
type Type struct {
	Id   int
	Name string
}

func init() {
	orm.RegisterModel(new(Region))
	orm.RegisterModel(new(Type))
}

func GetChannelRegion(channelId int) (int64, []Region, error) {
	o := orm.NewOrm()
	var regions []Region
	nums, err := o.Raw("select id,name from channel_region where status=1 and channel_id=? order by sort desc", channelId).QueryRows(&regions)
	return nums, regions, err
}

func GetChannelType(channelId int) (int64, []Type, error) {
	o := orm.NewOrm()
	var types []Type
	nums, err := o.Raw("select id,name from channel_type where status=1 and channel_id=? order by sort desc", channelId).QueryRows(&types)
	return nums, types, err
}
