package entity

import "time"

type CreditLimit struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	CustomerID  uint64     `gorm:"not null"`
	CreditLimit float64    `gorm:"type:decimal(16,2);not null"`
	TenorID     uint       `gorm:"not null"`
	StartDate   time.Time  `gorm:"not null"`
	EndDate     *time.Time `gorm:"default:null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	Customer    Customer   `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Tenor       Tenor      `gorm:"foreignKey:TenorID"`
}

func (c *CreditLimit) TableName() string {
	return "credit_limits"
}
