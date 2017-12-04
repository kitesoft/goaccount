package model

import (
	"goaccount/mixin"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint32    `gorm:"primary_key" json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UserName    string    `gorm:"type:varchar(128);index;unique" json:"username"`
	Password    string    `json:"password,omitempty" valid:"length(6|100)"`
	RawPassword string    `json:"raw_password,omitempty"`
	Email       string    `gorm:"type:varchar(128);index;unique" json:"email"`
	Phone       string    `gorm:"type:varchar(128);index;unique" json:"phone"`
	RealName    string    `json:"real_name"`
	Role        int32     `json:"role"`
}

func UserInfo(userName string) (User, mixin.ErrorCode) {
	var user User
	if err := DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
		logrus.Errorf("[User.Create] get user indo error %s", err.Error())
		return user, mixin.ErrorServerDb
	}
	return user, mixin.StatusOK
}

func CreateUser(user User) mixin.ErrorCode {
	if user.UserName == "" || user.Password == "" {
		return mixin.ErrorClientInvalidArgument
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[User.Create] bcrypt.GenerateFromPassword err: %s", err.Error())
		return mixin.ErrorServerCreateSecret
	}
	user.RawPassword = user.Password
	user.Password = string(hashedPassword)

	if err := DB.Create(&user).Error; err != nil {
		logrus.Errorf("[User.Create] create record error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func DeleteUser(id []int) mixin.ErrorCode {

	if err := DB.Where("id IN (?)", id).Delete(User{}).Error; err != nil {
		logrus.Errorf("[User.Delete] error %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func UpdateUser(user User) mixin.ErrorCode {

	if err := DB.Model(&user).Updates(user).Error; err != nil {
		logrus.Errorf("[User.Update] Updates err: %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func ResetPassword(userID []int) mixin.ErrorCode {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[User.Update] bcrypt.GenerateFromPassword err: %s", err.Error())
		return mixin.ErrorServerCreateSecret
	}

	if err := DB.Table("users").Where("id IN (?)", userID).Updates(map[string]string{"password": string(hashedPassword)}).Error; err != nil {
		logrus.Errorf("[User.Update] update err: %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func UpdatePassword(userName, oldPwd, newPwd string) mixin.ErrorCode {

	user, errCode := CheckPassword(userName, oldPwd)
	if errCode != mixin.StatusOK {
		return errCode
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[User.Update] bcrypt.GenerateFromPassword err: %s", err.Error())
		return mixin.ErrorServerCreateSecret
	}

	if err := DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		logrus.Errorf("[User.Update] update err: %s", err.Error())
		return mixin.ErrorServerDb
	}

	return mixin.StatusOK
}

func CheckPassword(userName, password string) (User, mixin.ErrorCode) {
	user, errCode := UserInfo(userName)
	if errCode != mixin.StatusOK {
		return user, errCode
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Errorf("[CheckPassword] bcrypt.CompareHashAndPassword %s, %s", user.Password, password)
		return user, mixin.ErrorClientUserOrPassword
	}
	return user, mixin.StatusOK
}

func GetUserName(email, phone string) (string, mixin.ErrorCode) {
	var user User
	var err error

	switch {
	case email != "":
		err = DB.Where("email = ?", email).First(&user).Error
	case phone != "":
		err = DB.Where("phone = ?", email).First(&user).Error
	default:
		return "", mixin.ErrorClientInvalidArgument
	}

	if err != nil {
		logrus.Errorf("[GetUserName] %s", err.Error())
		return "", mixin.ErrorServerDb
	}

	return user.UserName, mixin.StatusOK
}

func ListAllUser() ([]User, mixin.ErrorCode) {
	var users []User
	rows, err := DB.Raw("SELECT id, created_at, user_name, email, phone, real_name, role FROM users").Rows()
	if err != nil {
		return users, mixin.ErrorServerDb
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UserName, &user.Email, &user.Phone, &user.RealName, &user.Role)
		if err != nil {
			return users, mixin.ErrorServerDb
		}
		users = append(users, user)
	}

	return users, mixin.StatusOK
}
