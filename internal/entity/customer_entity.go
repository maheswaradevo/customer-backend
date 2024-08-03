package entity

import "time"

type Customer struct {
	ID           uint64     `gorm:"column:id;primaryKey"`
	IdNumber     string     `gorm:"column:nik_number"`
	FullName     string     `gorm:"column:full_name"`
	LegalName    string     `gorm:"column:legal_name"`
	BirthdayLoc  string     `gorm:"column:birthday_loc"`
	BirthdayDate time.Time  `gorm:"column:birthday_date"`
	Salary       float64    `gorm:"column:salary"`
	IdPic        string     `gorm:"column:id_pic"`
	SelfPic      string     `gorm:"column:self_pic"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (c *Customer) TableName() string {
	return "customers"
}
