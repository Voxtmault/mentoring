package interfaces

type User interface {
	GetUserByID(id uint) (User, error)
}
