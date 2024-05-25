package model

import "time"

type Client struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	ClientID     string    `json:"client_id" gorm:"column:client_id"`
	TaskID       int       `json:"task_id" gorm:"column:task_id"`
	Task         Task      `gorm:"foreignKey:TaskID;AssociationForeignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
	Model        string    `json:"model" gorm:"column:model"`
	Dataset      string    `json:"dataset" gorm:"column:dataset"`
	Type         string    `json:"type" gorm:"column:type"`
	Status       string    `json:"status" gorm:"column:status"`
	CurrentRound int       `json:"current_round" gorm:"column:current_round"`
	TotalRound   int       `json:"total_round" gorm:"column:total_round"`
	Progress     string    `json:"progress" gorm:"column:progress"`
	Accuracy     float32   `json:"accuracy" gorm:"column:accuracy"`
	Loss         float32   `json:"loss" gorm:"column:loss"`
	Cpu          string    `json:"cpu" gorm:"column:cpu"`
	Memory       string    `json:"memory" gorm:"column:memory"`
	Disk         string    `json:"disk" gorm:"column:disk"`
}
