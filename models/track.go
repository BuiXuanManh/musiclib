package models

type Track struct {
	TrackId     string `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string `json:"music_title" bson:"music_title"`
	Artist      string `json:"artist" bson:"artist"`
	Genre       string `json:"genre" bson:"genre"`
	ReleaseYear string `json:"release_year" bson:"release_year"`
	Duration    string `json:"duration" bson:"duration"`
	FileName    string `json:"file_name" bson:"file_name"`
}
