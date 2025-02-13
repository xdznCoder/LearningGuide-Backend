package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32          `gorm:"int,primary key" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"updatedAt"`
	DeleteAt  gorm.DeletedAt `gorm:"column:delete_time" json:"deleteAt"`
	IsDeleted bool           `gorm:"column:is_deleted" json:"isDeleted"`
}
