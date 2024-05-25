package dao

import (
	"context"

	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/model"

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

func (t *TrainDao) FindServerList(tid int, page int64, pageSize int64) (servers []*model.Server, pages int64, totalCount int64, err error) {
	err = global.DB.Model(&model.Server{}).Where("task_id = ?", tid).Count(&totalCount).Order("id desc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&servers).Error
	pages = totalCount / pageSize
	if totalCount%(pageSize) != 0 {
		pages++
	}
	return
}

func (t *TrainDao) FindClientList(tid int, page int64, pageSize int64) (clients []*model.Client, pages int64, totalCount int64, err error) {
	err = global.DB.Model(&model.Client{}).Where("task_id = ?", tid).Count(&totalCount).Order("id desc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&clients).Error
	pages = totalCount / pageSize
	if totalCount%(pageSize) != 0 {
		pages++
	}
	return
}
