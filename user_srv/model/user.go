package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32          `gorm:"primary key" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"updatedAt"`
	DeleteAt  gorm.DeletedAt `gorm:"column:delete_time" json:"deleteAt"`
	IsDeleted bool           `gorm:"column:is_deleted" json:"isDeleted"`
}

type User struct {
	BaseModel
	Email    string     `json:"email" gorm:"type:varchar(100);unique"`
	Password string     `json:"password" gorm:"type:varchar(100)"`
	NickName string     `json:"nickName" gorm:"type:varchar(100)"`
	Birthday *time.Time `gorm:"type:datetime" json:"birthday"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6)" json:"gender"`
	Desc     string     `gorm:"type:text" json:"desc"`
	Image    string     `gorm:"type:varchar(200)" json:"image"`
	Role     int        `gorm:"column:role;default:1;type:int" json:"role"`
}
