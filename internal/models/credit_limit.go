package models

import (
	"time"
)

type (
	CreditLimitCreateRequest struct {
		CustomerID    uint64  `json:"customer_id" form:"customer_id"`
		CreditLimit   float64 `json:"credit_limit" form:"credit_limit"`
		TenorID       uint64  `json:"tenor_id" form:"tenor_id"`
		StartDate     string  `json:"start_date" form:"start_date"`
		EndDate       string  `json:"end_date" form:"end_date"`
		StartDateTime time.Time
		EndDateTime   time.Time
	}

	CreditLimitUpdateRequest struct {
		ID            uint64  `param:"id"`
		CustomerID    uint64  `json:"customer_id" form:"customer_id"`
		CreditLimit   float64 `json:"credit_limit" form:"credit_limit"`
		TenorID       uint64  `json:"tenor_id" form:"tenor_id"`
		StartDate     string  `json:"start_date" form:"start_date"`
		EndDate       string  `json:"end_date" form:"end_date"`
		StartDateTime time.Time
		EndDateTime   time.Time
	}
)
