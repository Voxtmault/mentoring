package services

import (
	"github.com/voxtmault/mentoring/library-project/internal/users/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUsers(filter *models.UserFilter) ([]*models.User, error) {
	var user []*models.User
	query := s.db.Model(&models.User{})

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := s.db.First(&user, id).Error; err != nil {
		return User{}, err
	}
	return user, nil
}
