package dao

import (
	"context"

	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"github.com/palp1tate/FlowFederate/api/internal/model"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{global.DB.WithContext(ctx)}
}

func (u *UserDao) FindUserByMobile(mobile string) (user *model.User, err error) {
	err = u.Where("mobile = ? and role = ?", mobile, consts.User).First(&user).Error
	return
}

func (u *UserDao) CreateUser(user *model.User) (err error) {
	err = u.Create(&user).Error
	return
}

func (u *UserDao) FindUserById(id int) (user *model.User, err error) {
	err = u.Where("id = ? and role = ?", id, consts.User).First(&user).Error
	return
}

func (u *UserDao) UpdatePassword(user *model.User, password string) (err error) {
	err = u.Model(&user).Update("password", password).Error
	return
}

func (u *UserDao) UpdateUser(user *model.User) (err error) {
	err = u.Save(&user).Error
	return
}
