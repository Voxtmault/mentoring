package shared_models

import (
	"log/slog"

	"gorm.io/gorm"
)

type Gender struct {
	Model       `gorm:"embedded"`
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description" gorm:"not null"`
}

func (g *Gender) Seeder(db *gorm.DB) error {
	genders := []Gender{
		{Name: "Male", Description: ""},
		{Name: "Female", Description: ""},
		{Name: "Other", Description: "Use for unspecified / unknown case"},
	}

	for _, gender := range genders {
		if err := db.FirstOrCreate(&gender, gender).Error; err != nil {
			return err
		}
	}
	slog.Debug("gender seeded successfully")

	return nil
}
