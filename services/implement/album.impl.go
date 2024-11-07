package implements

import (
	"context"
	"musiclib/models"
	"musiclib/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlbumImpl struct {
	albumCollection *mongo.Collection
	trackCollection *mongo.Collection
	ctx             context.Context
}

func NewAlbumService(albumCollection *mongo.Collection, trackCollection *mongo.Collection, ctx context.Context) services.AlbumService {
	return &AlbumImpl{
		albumCollection: albumCollection,
		trackCollection: trackCollection,
		ctx:             ctx,
	}
}

func (a *AlbumImpl) CreateAlbum(album *models.Album) error {
	_, err := a.albumCollection.InsertOne(a.ctx, album)
	return err
}
func (a *AlbumImpl) GetAlbums() ([]models.Album, error) {
	var albums []models.Album
	cursor, err := a.albumCollection.Find(a.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(a.ctx, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}
func (a *AlbumImpl) UpdateAlbum(albumId *primitive.ObjectID, album *models.Album) error {
	filter := bson.M{"_id": albumId}
	update := bson.M{"$set": album}
	_, err := a.albumCollection.UpdateOne(a.ctx, filter, update)
	return err
}
func (a *AlbumImpl) DeleteAlbum(albumId *primitive.ObjectID) error {
	_, err := a.albumCollection.DeleteOne(a.ctx, albumId)
	return err
}
func (a *AlbumImpl) FindAlbum(albumId *primitive.ObjectID) (*models.Album, error) {
	var album *models.Album
	filter := bson.M{"_id": albumId}
	err := a.albumCollection.FindOne(a.ctx, filter).Decode(&album)
	return album, err
}
func (a *AlbumImpl) AddTrackToAlbum(albumId *primitive.ObjectID, track *models.Track) error {
	var album models.Album
	err := a.albumCollection.FindOne(a.ctx, bson.M{"_id": albumId}).Decode(&album)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if track.TrackId == "" {
		result, err := a.trackCollection.InsertOne(a.ctx, track)
		if err != nil {
			return err
		}
		track.TrackId = result.InsertedID.(primitive.ObjectID).Hex()
	}

	filter := bson.M{"_id": albumId}

	if err == mongo.ErrNoDocuments || len(album.Tracks) == 0 {
		update := bson.M{
			"$set": bson.M{
				"tracks": []models.Track{*track},
			},
		}
		_, err := a.albumCollection.UpdateOne(a.ctx, filter, update)
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"tracks": track,
		},
	}
	_, err = a.albumCollection.UpdateOne(a.ctx, filter, update)
	return err
}

func (a *AlbumImpl) AddExistedTrackToAlbum(albumId *primitive.ObjectID, trackId *primitive.ObjectID) error {
	var track models.Track
	err := a.trackCollection.FindOne(a.ctx, bson.M{"_id": trackId}).Decode(&track)
	if err != nil {
		return err
	}
	var album models.Album
	err = a.albumCollection.FindOne(a.ctx, bson.M{"_id": albumId}).Decode(&album)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": albumId}
	if err == mongo.ErrNoDocuments || len(album.Tracks) == 0 {
		update := bson.M{
			"$set": bson.M{
				"tracks": []models.Track{track},
			},
		}
		_, err := a.albumCollection.UpdateOne(a.ctx, filter, update)
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"tracks": track,
		},
	}
	_, err = a.albumCollection.UpdateOne(a.ctx, filter, update)
	return err
}

func (a *AlbumImpl) RemoveTrackFromAlbum(albumId *primitive.ObjectID, trackId *primitive.ObjectID) error {
	filter := bson.M{"_id": albumId}
	update := bson.M{"$pull": bson.M{"tracks": bson.M{"_id": trackId.Hex()}}}
	_, err := a.albumCollection.UpdateOne(a.ctx, filter, update)
	return err
}

func (a *AlbumImpl) FindTracksAndAlbums(keyword *string) ([]models.Album, []models.Track, error) {
	var albums []models.Album
	var tracks []models.Track

	// Tạo bộ lọc tìm kiếm gần đúng cho albums
	albumFilter := bson.M{
		"title": bson.M{"$regex": primitive.Regex{Pattern: *keyword, Options: "i"}},
	}

	// Tạo bộ lọc tìm kiếm gần đúng cho tracks
	trackFilter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": primitive.Regex{Pattern: *keyword, Options: "i"}}},
			{"artist": bson.M{"$regex": primitive.Regex{Pattern: *keyword, Options: "i"}}},
			{"album": bson.M{"$regex": primitive.Regex{Pattern: *keyword, Options: "i"}}},
			{"genre": bson.M{"$regex": primitive.Regex{Pattern: *keyword, Options: "i"}}},
		},
	}

	// Tìm albums dựa trên bộ lọc
	cursor, err := a.albumCollection.Find(a.ctx, albumFilter)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(a.ctx)

	// Giải nén kết quả tìm kiếm vào slice albums
	if err := cursor.All(a.ctx, &albums); err != nil {
		return nil, nil, err
	}

	// Tìm tracks dựa trên bộ lọc
	cursor, err = a.trackCollection.Find(a.ctx, trackFilter)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(a.ctx)

	// Giải nén kết quả tìm kiếm vào slice tracks
	if err := cursor.All(a.ctx, &tracks); err != nil {
		return nil, nil, err
	}

	return albums, tracks, nil
}
