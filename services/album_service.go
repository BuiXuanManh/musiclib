package services

import (
	"musiclib/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate go-mockgen-tool --type AlbumService
type AlbumService interface {
	CreateAlbum(*models.Album) error
	GetAlbums() ([]models.Album, error)
	FindAlbum(*primitive.ObjectID) (*models.Album, error)
	UpdateAlbum(*primitive.ObjectID, *models.Album) error
	DeleteAlbum(*primitive.ObjectID) error
	FindTracksAndAlbums(*string) ([]models.Album, []models.Track, error)
	AddTrackToAlbum(*primitive.ObjectID, *models.Track) error
	AddExistedTrackToAlbum(*primitive.ObjectID, *primitive.ObjectID) error
	RemoveTrackFromAlbum(*primitive.ObjectID, *primitive.ObjectID) error
}
