package dto

type AlbumDto struct {
	Title      string `json:"album_title" bson:"album_title"`
	AlbumCover string `json:"album_cover" bson:"album_cover"`
}
