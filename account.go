package main

import (
	"fmt"
	"goaccount/mixin"
	"goaccount/model"
	"net/http"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

func (this *Service) create_user_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.User{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Debugf("[Service.create_user_handle] validate %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.CreateUser(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

func (this *Service) login_handle(w http.ResponseWriter, r *http.Request) {

	var errCode, logResult mixin.ErrorCode

	inParam := &LoginRequest{}

	defer func(errCode mixin.ErrorCode) {
		model.AddLoginLog(model.LoginLog{
			CreatedAt: time.Now().Unix(),
			Name:      inParam.UserName,
			Ip:        r.RemoteAddr,
			Ua:        r.UserAgent(),
			Result:    errCode,
		})
	}(logResult)

	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[Service.login_handle] validate err: %s", err.Error())
		logResult = mixin.ErrorClientInvalidArgument
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if mactch, _ := regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", inParam.UserName); mactch {
		inParam.UserName, _ = model.GetUserName(inParam.UserName, "")
	} else if mactch, _ = regexp.MatchString("^[1][3578]\\d{9}$", inParam.UserName); mactch {
		inParam.UserName, _ = model.GetUserName("", inParam.UserName)
	}

	user, errCode := model.CheckPassword(inParam.UserName, inParam.Password)
	if errCode != mixin.StatusOK {
		logrus.Errorf("[Service.login_handle] owner.Login name:%s, password:%s, err_code: %d", inParam.UserName, inParam.Password, errCode)
		this.ResponseErrCode(w, mixin.ErrorClientUserOrPassword)
		return
	}

	token, err := this._jwt.PublicJWT().Encode(fmt.Sprint(user.ID), user.UserName, int64(user.Role))
	if err != nil {
		logrus.Errorf("[Service.login_handle] gen token error %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorServerCreateToken)
		return
	}

	response := &LoginResponse{
		UserId:   user.ID,
		UserName: inParam.UserName,
		Token:    token,
		Email:    user.Email,
		Phone:    user.Phone,
		RealName: user.RealName,
	}
	this.ResponseOK(w, response)
}

func (this *Service) update_user_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.User{}
	if err := this.validator.Validate(r, inParam); err != nil ||
		inParam.Password != "" || inParam.UserName != "" || inParam.Role != 0 {
		logrus.Debugf("[Service.create_user_handle] validate %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.UpdateUser(*inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)
}

func (this *Service) delete_user_handle(w http.ResponseWriter, r *http.Request) {

	inParam := &ResetPasswordRequest{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Debugf("[Service.reset_password_handle] validate %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.DeleteUser(inParam.UserID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}

	this.ResponseOK(w, nil)

}

func (this *Service) update_password_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &UpdatePasswordRequest{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Debugf("[Service.update_password_handle] validate %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.UpdatePassword(inParam.UserName, inParam.OldPassword, inParam.NewPassword)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *Service) reset_password_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &ResetPasswordRequest{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Debugf("[Service.reset_password_handle] validate %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	errCode := model.ResetPassword(inParam.UserID)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *Service) list_user_handle(w http.ResponseWriter, r *http.Request) {
	user, errCode := model.ListAllUser()
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, user)
}

func (this *Service) login_log_handle(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("username")
	resp := model.QueryLoginLog(userName)
	fmt.Fprintf(w, `{"login_log":%v}`, resp)
}

func (this *Service) logout_handle(w http.ResponseWriter, r *http.Request) {
	this.ResponseOK(w, nil)
}

// 实名认证
func (this *Service) real_name_auth_handle(w http.ResponseWriter, r *http.Request) {
	inParam := &model.RealNameAuthInfo{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[Service.real_name_auth_handle] validate err: %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	errCode := model.CreateRealNameAuth(inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *Service) update_real_name_auth(w http.ResponseWriter, r *http.Request) {
	inParam := &model.RealNameAuthInfo{}
	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[Service.real_name_auth_handle] validate err: %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	errCode := model.UpdateRealNameAuth(inParam)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *Service) real_name_auth_info(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("username")
	if userName == "" {
		logrus.Errorf("[Service.real_name_auth_info] username is empty")
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	resp, errCode := model.GetRealNameAuthInfo(userName)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, resp)
}

func (this *Service) pass_real_name_auth(w http.ResponseWriter, r *http.Request) {
	inParams := &PassRealNameAuth{}
	if err := this.validator.Validate(r, inParams); err != nil {
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	errCode := model.RealNameAuthing(inParams.ID, inParams.Result, inParams.Comment)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, nil)
}

func (this *Service) real_name_auth_list(w http.ResponseWriter, r *http.Request) {
	var inParam map[string]interface{}
	if err := this.validator.Validate(r, &inParam); err != nil {
		logrus.Errorf("[Service.real_name_auth_list] validate err: %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}
	from, to := inParam["from"].(int), inParam["to"].(int)
	delete(inParam, "from")
	delete(inParam, "to")

	resp, errCode := model.GetAuthReviewList(inParam, from, to)
	if errCode != mixin.StatusOK {
		this.ResponseErrCode(w, errCode)
		return
	}
	this.ResponseOK(w, resp)
}
