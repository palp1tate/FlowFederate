package model

type UserInfo struct {
	ID       int    `gorm:"autoIncrement:false"`
	UserName string `gorm:"size:255;uniqueIndex;primaryKey"`
	Password string `gorm:"size:255"`
	Role     int    `gorm:"default:1;check:Role in (0,1)"`
	State    int    `gorm:"default:0;check:State in (0,1)"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}
