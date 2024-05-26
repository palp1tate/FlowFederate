package dao

import (
	"context"

	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"github.com/palp1tate/FlowFederate/api/internal/model"

	"gorm.io/gorm"
)

type AdminDao struct {
	*gorm.DB
}

func NewAdminDao(ctx context.Context) *AdminDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &AdminDao{global.DB.WithContext(ctx)}
}

func (u *AdminDao) FindAdminByMobile(mobile string) (user *model.User, err error) {
	err = u.Where("mobile = ? and role = ?", mobile, consts.Admin).First(&user).Error
	return
}

func (u *AdminDao) FindAdminById(id int) (user *model.User, err error) {
	err = u.Where("id = ? and role = ?", id, consts.Admin).First(&user).Error
	return
}

func (u *AdminDao) FindUserList(page int64, pageSize int64) (users []*model.User, pages int64, totalCount int64, err error) {
	err = global.DB.Model(&model.User{}).Where("role = ?", consts.User).Count(&totalCount).Order("id desc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&users).Error
	pages = totalCount / pageSize
	if totalCount%(pageSize) != 0 {
		pages++
	}
	return
}

func (u *AdminDao) DeleteUser(user *model.User) (err error) {
	err = global.DB.Delete(&user).Error
	return
}
