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

	var errCode mixin.ErrorCode

	inParam := &LoginRequest{}

	defer func() {
		model.AddLoginLog(model.LoginLog{
			CreatedAt: time.Now().Unix(),
			Name:      inParam.UserName,
			Ip:        r.RemoteAddr,
			Ua:        r.UserAgent(),
			Result:    errCode,
		})
	}()

	if err := this.validator.Validate(r, inParam); err != nil {
		logrus.Errorf("[Service.login_handle] validate err: %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorClientInvalidArgument)
		return
	}

	if mactch, _ := regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", inParam.UserName); mactch {
		inParam.UserName, errCode = model.GetUserName(inParam.UserName, "")
	} else if mactch, _ = regexp.MatchString("^[1][3578]\\d{9}$", inParam.UserName); mactch {
		inParam.UserName, errCode = model.GetUserName("", inParam.UserName)
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
