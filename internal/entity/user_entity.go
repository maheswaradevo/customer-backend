package entity

import "time"

type User struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	Email     string    `gorm:"column:email"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
