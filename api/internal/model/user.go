package model

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Username  string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Mobile    string    `gorm:"not null"`
	Avatar    string    `gorm:"not null;default:http://se1437foq.hn-bkt.clouddn.com/avatar.jpg"`
	Role      int       `gorm:"not null;default:1;check: role in(1,2)"` // 1:普通用户 2:管理员
}
