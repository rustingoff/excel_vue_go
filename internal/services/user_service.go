package services

import (
	"github.com/rustingoff/excel_vue_go/internal/models"
	"github.com/rustingoff/excel_vue_go/internal/repositories"
	"github.com/rustingoff/excel_vue_go/packages/token"
	"log"
)

type UserService interface {
	CreateUser(user models.User) error
	DeleteUser(id string) error
	GetUserById(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)

	Login(email, token string) error
}

type userService struct {
	repo repositories.UserRepository
}

func GetUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
} //nolint:typechecking

func (service *userService) CreateUser(user models.User) error {
	hashedPassword, err := token.GenerateHashPassword(user.Password)
	if err != nil {
		log.Println("[ERR]: failed to generate hash")
		return err
	}

	user.Password = hashedPassword
	user.Active = true
	return service.repo.CreateUser(user)
}

func (service *userService) DeleteUser(id string) error {
	return service.repo.DeleteUser(id)
}

func (service *userService) GetUserById(id string) (models.User, error) {
	return service.repo.GetUserById(id)
}

func (service *userService) GetUserByEmail(email string) (models.User, error) {
	return service.repo.GetUserByEmail(email)
}

func (service *userService) Login(email, token string) error {
	return service.repo.Login(email, token)
}
