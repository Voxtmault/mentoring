package interfaces

import (
	"context"

	"github.com/voxtmault/mentoring/library-project/internal/users/models"
	"github.com/voxtmault/mentoring/library-project/pkg/http_utility"
)

type User interface {
	GetUsers(ctx context.Context, filter *models.UserFilter) (res *http_utility.HTTPResponse, data []*models.User)
	GetUserByID(ctx context.Context, id uint) (res *http_utility.HTTPResponse, data *models.User)
	CreateUser(ctx context.Context, user *models.User) (res *http_utility.HTTPResponse, data *models.User)
}
