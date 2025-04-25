package models

import (
	"github.com/voxtmault/mentoring/library-project/pkg/pagination"
	"github.com/voxtmault/mentoring/library-project/pkg/shared_models"
)

type User struct {
	shared_models.Model `gorm:"embedded"`
	Name                string                       `json:"name" gorm:"not null;default:'John Doe'" validate:"required"`
	Gender              string                       `json:"gender" gorm:"not null" validate:"required"`
	Username            string                       `json:"username" gorm:"unique;not null"`
	Email               string                       `json:"email" gorm:"unique;not null"`
	Password            string                       `json:"-" gorm:"not null"`
	Salt                string                       `json:"-" gorm:"not null"`
	Address             []*Address                   `json:"address" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsValidated         bool                         `json:"is_validated" gorm:"default:false"`
	ValidatedAt         shared_models.CustomNullTime `json:"validated_at" gorm:"default:null"`
}

type UserFilter struct {
	Username    string `query:"username" validate:"omitempty,alphanumunicode"`
	Email       string `query:"email" validate:"omitempty"`
	IsValidated *bool  `query:"is_validated" validate:"omitempty"`

	pagination.PaginationFilter
}

type Address struct {
	shared_models.Model `gorm:"embedded"`
	UserID              uint   `json:"user_id" gorm:"not null"`
	Street              string `json:"street" gorm:"not null" validate:"required"`
	City                string `json:"city" gorm:"not null" validate:"required"`
	State               string `json:"state" gorm:"not null" validate:"required"`
	ZipCode             string `json:"zip_code" gorm:"not null" validate:"required"`
	Country             string `json:"country" gorm:"not null" validate:"required"`
}
