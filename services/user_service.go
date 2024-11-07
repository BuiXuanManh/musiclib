package services

import (
	"musiclib/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	CreateUser(*models.User) error
	GetUser(*string) (*models.User, error)
	UpdateUser(*models.User) error
	ChangePassword(*string, *string, *string) error
	DeleteUser(*primitive.ObjectID) error
	GetUserFromUsername(*string) (*models.User, error)
}
