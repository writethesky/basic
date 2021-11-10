package dao

import (
	"basic/internal"
	"basic/repository/entity"
)

func GetUserIDByUsername(username string) (userID int64, err error) {
	userInfo, err := GetUserInfo(username)
	if nil != err {
		return
	}
	userID = int64(userInfo.ID)
	return
}

func CreateUser(user entity.User) (err error) {
	return internal.DB.Save(&user).Error
}

func GetUserInfo(username string) (userInfo entity.User, err error) {
	err = internal.DB.First(&userInfo, "username = ?", username).Error
	return
}

func SetPassword(id int64, password, salt string) error {
	return internal.DB.Model(&entity.User{}).Where("id=?", id).Updates(entity.User{
		Password: password,
		Salt:     salt,
	}).Error
}
