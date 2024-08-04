package models

import (
	"mime/multipart"
	"time"
)

type (
	UserCreateRequest struct {
		Email    string `json:"email" form:"email"`
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	UserUpdateRequest struct {
		ID       uint64 `param:"id"`
		Email    string `json:"email" form:"email"`
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	UserGetRequest struct {
		ID       uint64 `query:"id"`
		Email    string `query:"email"`
		Username string `query:"username"`
	}
)

type (
	UserResponse struct {
		ID       uint64 `json:"id"`
		Email    string `json:"email" form:"email"`
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
)

type (
	CustomerCreateRequest struct {
		IdNumber         string `json:"id_number" form:"id_number"`
		FullName         string `json:"full_name" form:"full_name"`
		LegalName        string `json:"legal_name" form:"legal_name"`
		BirthdayLoc      string `json:"birthday_loc" form:"birthday_loc"`
		BirthdayDate     string `json:"birthday_date" form:"birthday_date"`
		BirthdayDateTime time.Time
		Salary           float64               `json:"salary" form:"salary"`
		IdPic            *multipart.FileHeader `json:"id_pic" form:"id_pic"`
		SelfPic          *multipart.FileHeader `json:"self_pic" form:"self_pic"`
		IdPicUrl         string
		SelfPicUrl       string
	}

	CustomerUpdateRequest struct {
		ID           uint64                `param:"id"`
		IdNumber     string                `json:"id_number" form:"id_number"`
		FullName     string                `json:"full_name" form:"full_name"`
		LegalName    string                `json:"legal_name" form:"legal_name"`
		BirthdayLoc  string                `json:"birthday_loc" form:"birthday_loc"`
		BirthdayDate string                `json:"birthday_date" form:"birthday_date"`
		Salary       float64               `json:"salary" form:"salary"`
		IdPic        *multipart.FileHeader `json:"id_pic" form:"id_pic"`
		SelfPic      *multipart.FileHeader `json:"self_pic" form:"self_pic"`
		IdPicUrl     string
		SelfPicUrl   string
	}
)

type (
	CustomerRegisterRequest struct {
		Email            string `json:"email" form:"email"`
		Username         string `json:"username" form:"username"`
		Password         string `json:"password" form:"password"`
		IdNumber         string `json:"id_number" form:"id_number"`
		FullName         string `json:"full_name" form:"full_name"`
		LegalName        string `json:"legal_name" form:"legal_name"`
		BirthdayLoc      string `json:"birthday_loc" form:"birthday_loc"`
		BirthdayDate     string `json:"birthday_date" form:"birthday_date"`
		BirthdayDateTime time.Time
		Salary           float64               `json:"salary" form:"salary"`
		IdPic            *multipart.FileHeader `json:"id_pic" form:"id_pic"`
		SelfPic          *multipart.FileHeader `json:"self_pic" form:"self_pic"`
		IdPicUrl         string
		SelfPicUrl       string
	}
)

type (
	CustomerRegisterResponse struct {
		UserCreateRequest
		IdNumber     string  `json:"id_number" form:"id_number"`
		FullName     string  `json:"full_name" form:"full_name"`
		LegalName    string  `json:"legal_name" form:"legal_name"`
		BirthdayLoc  string  `json:"birthday_loc" form:"birthday_loc"`
		BirthdayDate string  `json:"birthday_date" form:"birthday_date"`
		Salary       float64 `json:"salary" form:"salary"`
		IdPicUrl     string  `json:"id_pic"`
		SelfPicUrl   string  `json:"self_pic"`
	}
)

type (
	LoginRequest struct {
		Email    string `json:"email" form:"email"`
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	LoginResponse struct {
		UserID              uint64 `json:"user_id"`
		Username            string `json:"username"`
		Email               string `json:"email"`
		AccessToken         string `json:"access_token"`
		RefreshToken        string `json:"refresh_token"`
		ExpiredToken        int64  `json:"expired_token"`
		ExpiredRefreshToken int64  `json:"expired_refresh_token"`
	}
)
