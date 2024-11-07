package models

type Album struct {
	AlbumId    string  `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string  `json:"album_title" bson:"album_title"`
	AlbumCover string  `json:"album_cover" bson:"album_cover"`
	Tracks     []Track `json:"tracks" bson:"tracks"`
}
