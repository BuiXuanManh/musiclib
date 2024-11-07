package implements

import (
	"context"
	"musiclib/models"
	"musiclib/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrackImpl struct {
	trackCollection *mongo.Collection
	ctx             context.Context
}

func (t *TrackImpl) CreateTrack(track *models.Track) error {
	_, err := t.trackCollection.InsertOne(t.ctx, track)
	return err
}
func NewTrackService(trackCollection *mongo.Collection, ctx context.Context) services.TrackService {
	return &TrackImpl{
		trackCollection: trackCollection,
		ctx:             ctx,
	}
}
func (t *TrackImpl) GetTracks() ([]models.Track, error) {
	var tracks []models.Track
	cursor, err := t.trackCollection.Find(t.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(t.ctx, &tracks); err != nil {
		return nil, err
	}
	return tracks, nil
}
func (t *TrackImpl) UpdateTrack(trackId *primitive.ObjectID, track *models.Track) error {
	filter := bson.M{"_id": trackId}
	update := bson.M{"$set": track}
	_, err := t.trackCollection.UpdateOne(t.ctx, filter, update)
	return err
}
func (t *TrackImpl) DeleteTrack(trackId *primitive.ObjectID) error {
	_, err := t.trackCollection.DeleteOne(t.ctx, bson.M{"_id": trackId})
	return err
}
func (t *TrackImpl) FindTrack(trackId *primitive.ObjectID) (*models.Track, error) {
	var track *models.Track
	filter := bson.M{"_id": trackId}
	err := t.trackCollection.FindOne(t.ctx, filter).Decode(&track)
	return track, err
}
