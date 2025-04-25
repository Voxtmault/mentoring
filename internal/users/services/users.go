package services

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/internal/users/interfaces"
	"github.com/voxtmault/mentoring/library-project/internal/users/models"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"github.com/voxtmault/mentoring/library-project/pkg/http_utility"
	"github.com/voxtmault/mentoring/library-project/pkg/pagination"
	"github.com/voxtmault/mentoring/library-project/pkg/security"
	"gorm.io/gorm"
)

type UserService struct {
	db  *gorm.DB
	cfg *config.AppConfig
}

// Ensure UserService implements the User interface
var _ interfaces.User = (*UserService)(nil)

// New initializes a new UserService instance and migrates the database
// schema for User and Address models.
func New(db *gorm.DB, cfg *config.AppConfig) *UserService {
	db.AutoMigrate(&models.User{}, &models.Address{})
	return &UserService{db: db, cfg: cfg}
}

func (s *UserService) GetUsers(ctx context.Context, filter *models.UserFilter) (res *http_utility.HTTPResponse, data []*models.User) {
	res = http_utility.New(s.cfg)
	data = make([]*models.User, 0)
	var metadata pagination.PaginationMetadata

	query := s.db.Model(&models.User{}).Preload("Address")

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}
	if filter.IsValidated != nil {
		query = query.Where("is_validated = ?", *filter.IsValidated)
	}

	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.InternalErrorMessage = http_utility.GeneralError
		res.ErrorStack = eris.Wrap(err, "failed to get users count")

		return
	}

	offset := (filter.PageNumber - 1) * filter.Limit
	if err := query.Limit(int(filter.Limit)).Offset(int(offset)).Find(&data).Error; err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.InternalErrorMessage = http_utility.GeneralError
		res.ErrorStack = eris.Wrap(err, "failed to get users")

		return
	}

	// Calculate the metadata
	pagination.CalculateMetadata(ctx, totalRecords, &metadata, &filter.PaginationFilter)

	res.StatusCode = http.StatusOK
	res.InternalMessage = http_utility.Success
	res.Metadata = metadata
	res.Data = data

	return
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (res *http_utility.HTTPResponse, data *models.User) {
	res = http_utility.New(s.cfg)

	if err := s.db.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			res.StatusCode = http.StatusNotFound
			res.InternalErrorMessage = http_utility.ResourceNotFound
			res.ErrorStack = eris.Wrap(err, "user not found")
		} else {
			res.StatusCode = http.StatusInternalServerError
			res.InternalErrorMessage = http_utility.GeneralError
			res.ErrorStack = eris.Wrap(err, "failed to get user")
		}

		return
	}

	res.StatusCode = http.StatusOK
	res.InternalMessage = http_utility.Success
	res.Data = data

	return
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (res *http_utility.HTTPResponse, data *models.User) {
	res = http_utility.New(s.cfg)

	// Generate salt and hash password
	password := security.GenerateRandomPassword(s.cfg.SecurityConfig.MinPasswordLength, &s.cfg.SecurityConfig)
	hashPw, salt, err := security.HashPassword(password, &s.cfg.SecurityConfig)
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.InternalErrorMessage = http_utility.GeneralError
		res.ErrorStack = eris.Wrap(err, "failed to hash password")

		return
	}

	user.Password = hashPw
	user.Salt = salt

	tx := s.db.Begin(&sql.TxOptions{})

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		res.StatusCode = http.StatusInternalServerError
		res.InternalErrorMessage = http_utility.GeneralError
		res.ErrorStack = eris.Wrap(err, "failed to create user")

		return
	}

	if err := tx.Commit().Error; err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.InternalErrorMessage = http_utility.GeneralError
		res.ErrorStack = eris.Wrap(err, "failed to commit transaction")

		return
	}

	data = user

	res.StatusCode = http.StatusCreated
	res.InternalMessage = http_utility.Created
	res.Data = data

	return
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) (res *http_utility.HTTPResponse, data *models.User) {
	res = http_utility.New(s.cfg)

	return res, nil
}
