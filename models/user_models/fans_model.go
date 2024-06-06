package user_models

import "live/common/models"

type FansModel struct {
	models.Model
	AnchorId uint32 `json:"anchorId"` //主播id
	UserId   uint32 `json:"userId"`   //粉丝id
	Level    int8   `json:"level"`    //等级
	IsFans   bool   `json:"isFans"`   //是否关注
}
