package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/voxtmault/mentoring/library-project/pkg/pagination"
)

type User struct {
	gorm.Model  `gorm:"embedded"`
	Username    string     `json:"username" gorm:"unique;not null"`
	Email       string     `json:"email" gorm:"unique;not null"`
	Password    string     `json:"-" gorm:"not null"`
	Salt        string     `json:"-" gorm:"not null"`
	Address     []*Address `json:"address" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsValidated bool       `json:"is_validated" gorm:"default:true"`
	ValidatedAt time.Time  `json:"validated_at" gorm:"default:null"`
}

type UserFilter struct {
	Username string `query:"username"`
	Email    string `query:"email"`
	IsActive *bool  `query:"is_active"`

	pagination.PaginationFilter
}

type Address struct {
	gorm.Model `gorm:"embedded"`
	UserID     uint   `json:"user_id" gorm:"not null"`
	Street     string `json:"street" gorm:"not null"`
	City       string `json:"city" gorm:"not null"`
	State      string `json:"state" gorm:"not null"`
	ZipCode    string `json:"zip_code" gorm:"not null"`
	Country    string `json:"country" gorm:"not null"`
}
