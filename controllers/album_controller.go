package controllers

import (
	"musiclib/models"
	"musiclib/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlbumController struct {
	albumService services.AlbumService
}

func NewAlbumController(albumService services.AlbumService) *AlbumController {
	return &AlbumController{
		albumService: albumService,
	}
}
func CheckValidAlbum(album *models.Album) bool {
	if album.AlbumCover == "" || album.Title == "" {
		return false
	}
	return true
}

// CreateAlbum	godoc
// @Summary      CreateAlbum
// @Description  create a Album
// @Tags         album
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        album   body     dto.AlbumDto  true  "Album data to create"
// @Router       /album/create [post]
func (a *AlbumController) CreateAlbum(ctx *gin.Context) {
	var album *models.Album
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !CheckValidAlbum(album) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}
	if err := a.albumService.CreateAlbum(album); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, album)
}

// ListAlbum godoc
// @Summary      List albums
// @Description  get albums
// @Tags         album
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Album
// @Router       /album/getAll [get]
func (a *AlbumController) GetAlbums(ctx *gin.Context) {
	albums, err := a.albumService.GetAlbums()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, albums)
}

// UpdateAlbum 	godoc
// @Summary      UpdateAlbum
// @Description  Update a album
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Update by Album ID"
// @Param        album   body     dto.AlbumDto  true  "Album data to update"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /album/update/{id} [put]
func (a *AlbumController) UpdateAlbum(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var album models.Album
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.albumService.UpdateAlbum(&id, &album); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, album)
}

// DeleteAlbum 	godoc
// @Summary      DeleteAlbum
// @Description  delete a album
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Delete by Album ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /album/delete/{id} [delete]
func (a *AlbumController) DeleteAlbum(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.albumService.DeleteAlbum(&id); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}

// GetAlbum 	godoc
// @Summary      GetAlbum
// @Description  Get a album
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Find by album ID"
// @Success      200  {object}   models.Album
// @Router       /album/get/{id} [get]
func (a *AlbumController) FindAlbum(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	album, err := a.albumService.FindAlbum(&id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, album)
}

// GetTrackAndAlbum 	godoc
// @Summary      GetTrackAndAlbum
// @Description  Get list tracks and albums by keyword
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        keyword  query  string  true  "Search by keyword"
// @Router       /album/search [get]
func (a *AlbumController) FindTracksAndAlbums(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	albums, tracks, err := a.albumService.FindTracksAndAlbums(&keyword)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"albums": albums, "tracks": tracks})
}

// AddTrackToAlbum	godoc
// @Summary      AddTrackToAlbum
// @Description  Add a Track to Album
// @Tags         album
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id  path  string  true  "Find by album ID"
// @Param        track   body     dto.TrackDto  true  "Track data to create"
// @Router       /album/add_track/{id} [post]
func (a *AlbumController) AddTrackToAlbum(ctx *gin.Context) {
	albumId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var track models.Track
	if err := ctx.ShouldBindJSON(&track); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.albumService.AddTrackToAlbum(&albumId, &track); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, track)
}

// RemoveTrackFromAlbum 	godoc
// @Summary      RemoveTrackFromAlbum
// @Description  Remove a Track from Album
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Find by Album ID"
// @Param        trackId  path  string  true  "Remove by Track ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /album/remove_track/{id}/{trackId} [put]
func (a *AlbumController) RemoveTrackFromAlbum(ctx *gin.Context) {
	albumId, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trackId, err := primitive.ObjectIDFromHex(ctx.Param("trackId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.albumService.RemoveTrackFromAlbum(&albumId, &trackId); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "track removed from album"})
}

func (a *AlbumController) RegisterAlbumRouter(rt *gin.RouterGroup) {
	router := rt.Group("/album")
	router.POST("/create", a.CreateAlbum)
	router.GET("/getAll", a.GetAlbums)
	router.GET("/find/:id", a.FindAlbum)
	router.PUT("/update/:id", a.UpdateAlbum)
	router.DELETE("/delete/:id", a.DeleteAlbum)
	router.GET("/search", a.FindTracksAndAlbums)
	router.POST("/add_track/:id", a.AddTrackToAlbum)
	router.PUT("/remove_track/:id/:trackId", a.RemoveTrackFromAlbum)
}
