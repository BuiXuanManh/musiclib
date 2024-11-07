package services

import (
	"musiclib/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackService interface {
	CreateTrack(*models.Track) error
	GetTracks() ([]models.Track, error)
	FindTrack(*primitive.ObjectID) (*models.Track, error)
	UpdateTrack(*primitive.ObjectID, *models.Track) error
	DeleteTrack(*primitive.ObjectID) error
}
