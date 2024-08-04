package models

import "time"

type UserEvent struct {
	ID                  uint64        `json:"id"`
	Username            string        `json:"username"`
	Email               string        `json:"email"`
	AccessToken         string        `json:"access_token"`
	RefreshToken        string        `json:"refresh_token"`
	ExpiredToken        time.Duration `json:"expired_token"`
	ExpiredRefreshToken time.Duration `json:"expired_refresh_token"`
}

func (u *UserEvent) GetId() uint64 {
	return u.ID
}
