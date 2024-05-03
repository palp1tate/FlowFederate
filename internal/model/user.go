package model

type UserInfo struct {
	Uuid       string `json:"uuid" gorm:"column:uuid"`
	UserName   string `json:"user_name" gorm:"column:user_name"`
	Password   string `json:"password" gorm:"column:password"`
	Role       int    `json:"role" gorm:"column:role"`
	State      int    `json:"state" gorm:"column:state"`
	CreateTime string `json:"create_time" gorm:"column:create_time"`
}

func (m *UserInfo) TableName() string {
	return "user_info"
}
