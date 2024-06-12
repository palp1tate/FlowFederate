package model

import "time"

type Task struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserName  string    `json:"user_name" gorm:"column:user_name"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Model     string    `json:"model" gorm:"column:model"`
	Dataset   string    `json:"dataset" gorm:"column:dataset"`
	Type      string    `json:"type" gorm:"column:type"`
	Status    string    `json:"status" gorm:"column:status"`
	Progress  string    `json:"progress" gorm:"column:progress"`
	Accuracy  string    `json:"accuracy" gorm:"column:accuracy;size:2048"`
	Loss      string    `json:"loss" gorm:"column:loss;size:2048"`
}
