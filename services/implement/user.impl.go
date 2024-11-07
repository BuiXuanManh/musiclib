package implements

import (
	"context"
	"errors"
	"musiclib/helper"
	"musiclib/models"
	"musiclib/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) services.UserService {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.User) error {
	us, err := u.GetUserFromUsername(&user.Username)
	if us != nil {
		return errors.New("user already exists")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	_, err = u.userCollection.InsertOne(u.ctx, user)
	return err
}

func (u *UserServiceImpl) GetUser(userId *string) (*models.User, error) {
	var user *models.User
	id, err := primitive.ObjectIDFromHex(*userId)
	if err != nil {
		return nil, err
	}
	query := bson.D{bson.E{Key: "_id", Value: id}}
	err = u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	id, err := primitive.ObjectIDFromHex(user.UserId)
	if err != nil {
		return err
	}
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	update := bson.D{
		bson.E{Key: "$set",
			Value: bson.D{
				bson.E{Key: "username", Value: user.Username},
				bson.E{Key: "password", Value: user.Password},
			},
		},
	}
	result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) ChangePassword(username *string, OldPassword *string, NewPassword *string) error {
	// Retrieve the user by ID
	filter := bson.D{bson.E{Key: "username", Value: username}}
	var existingUser *models.User
	err := u.userCollection.FindOne(u.ctx, filter).Decode(&existingUser)
	if err != nil {
		return err
	}
	// Compare password with password hash
	if !helper.CheckPassword(existingUser.Password, *OldPassword) {
		return errors.New("wrong password")
	}
	hashedPassword, err := helper.HashPassword(*NewPassword)
	if err != nil {
		return err
	}
	existingUser.Password = string(hashedPassword)

	update := bson.D{
		bson.E{Key: "$set",

			Value: bson.D{
				bson.E{Key: "password", Value: existingUser.Password},
			},
		},
	}

	result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(userId *primitive.ObjectID) error {
	filter := bson.D{bson.E{Key: "_id", Value: userId}}
	result, _ := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}

func (u *UserServiceImpl) GetUserFromUsername(username *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "username", Value: username}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}
