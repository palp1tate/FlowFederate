package dao

import (
	"context"

	"github.com/palp1tate/FlowFederate/internal/model"

	"github.com/palp1tate/FlowFederate/api/global"

	"gorm.io/gorm"
)

type TrainDao struct {
	*gorm.DB
}

func NewTrainDao(ctx context.Context) *TrainDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &TrainDao{global.DB.WithContext(ctx)}
}

func (t *TrainDao) FindTaskList(userName string, page int64, pageSize int64) (tasks []*model.Task, pages int64, totalCount int64, err error) {
	err = global.DB.Model(&model.Task{}).Where("user_name = ?", userName).Count(&totalCount).Order("id desc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&tasks).Error
	pages = totalCount / pageSize
	if totalCount%(pageSize) != 0 {
		pages++
	}
	return
}

func (t *TrainDao) FindTask(tid int) (task *model.Task, err error) {
	err = global.DB.Where("id = ?", tid).First(&task).Error
	return
}
