package model

import (
	"goaccount/mixin"
	"time"

	"github.com/sirupsen/logrus"
)

type RealNameAuthInfo struct {
	ID             uint `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UserName       string `json:"user_name,omitempty"`
	RealName       string `json:"real_name,omitempty"`
	IdNumber       string `json:"id_number,omitempty"`
	IdCardPositive string `json:"id_card_positive,omitempty"`
	IdCardNegative string `json:"id_card_negative,omitempty"`
	IdCardWithhand string `json:"id_card_withhand,omitempty"`
	State          string `json:"state,omitempty"` //authing passing reject
	Comment        string `json:"comment,omitempty"`
}

func CreateRealNameAuth(info *RealNameAuthInfo) mixin.ErrorCode {
	if err := DB.Create(info).Error; err != nil {
		logrus.Errorf("[User.CreateRealNameAuth] create error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func UpdateRealNameAuth(info *RealNameAuthInfo) mixin.ErrorCode {
	info.State = "authing"
	info.Comment = ""
	if err := DB.Save(info).Error; err != nil {
		logrus.Errorf("[User.UpdateRealNameAuth] create error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}

func GetRealNameAuthInfo(userName string) (RealNameAuthInfo, mixin.ErrorCode) {
	var info RealNameAuthInfo
	if err := DB.Where("user_name = ?", userName).First(&info).Error; err != nil {
		logrus.Errorf("[User.GetRealNameAuthInfo] error %s", err.Error())
		return info, mixin.ErrorServerDb
	}

	return info, mixin.StatusOK
}

func GetAuthReviewList(searchMap map[string]interface{}, from, to int) ([]RealNameAuthInfo, mixin.ErrorCode) {
	var infos []RealNameAuthInfo
	if err := DB.Where(searchMap).Find(&infos).Error; err != nil {
		logrus.Errorf("[User.GetAuthReviewList] error %s", err.Error())
		return infos, mixin.ErrorServerDb
	}
	return infos, mixin.StatusOK
}

func RealNameAuthing(id uint, result, comment string) mixin.ErrorCode {

	if result != "pass" {
		result = "reject"
	}
	info := RealNameAuthInfo{
		ID:      id,
		State:   result,
		Comment: comment,
	}

	if err := DB.Model(&info).Updates(info).Error; err != nil {
		logrus.Errorf("[User.RealNameAuthing] error %s", err.Error())
		return mixin.ErrorServerDb
	}
	return mixin.StatusOK
}
