package entity

type (
	Tenor struct {
		ID               uint64 `gorm:"column:id"`
		TenorDescription string `gorm:"column:tenor_description"`
	}
)

func (t *Tenor) TableName() string {
	return "tenors"
}
