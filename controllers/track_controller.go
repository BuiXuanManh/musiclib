package controllers

import (
	"musiclib/models"
	"musiclib/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackController struct {
	trackService services.TrackService
}

func NewTrackController(trackService services.TrackService) *TrackController {
	return &TrackController{
		trackService: trackService,
	}
}
func CheckValidTrack(track *models.Track) bool {
	if track.Title == "" || track.Artist == "" || track.Genre == "" || track.Duration == "" || track.FileName == "" || track.ReleaseYear == "" {
		return false
	}
	return true
}

// CreateTrack 	godoc
// @Summary      CreateTrack
// @Description  create a Track
// @Tags         track
// @Accept       json
// @Produce      json
// @Param        track   body     dto.TrackDto  true  "Track data to create"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /track/create [post]
func (t *TrackController) CreateTrack(ctx *gin.Context) {
	var track models.Track
	if err := ctx.ShouldBindJSON(&track); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !CheckValidTrack(&track) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}
	if err := t.trackService.CreateTrack(&track); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, track)
}

// ListTracks godoc
// @Summary      List tracks
// @Description  get tracks
// @Tags         track
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Track
// @Router       /track/getAll [get]
func (t *TrackController) GetTracks(ctx *gin.Context) {
	tracks, err := t.trackService.GetTracks()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tracks)
}

// UpdateTrack 	godoc
// @Summary      UpdateTrack
// @Description  Update a track
// @Tags         track
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Update by Track ID"
// @Param        track   body     dto.TrackDto  true  "Track data to update"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /track/update/{id} [put]
func (t *TrackController) UpdateTrack(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var track models.Track
	if err := ctx.ShouldBindJSON(&track); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !CheckValidTrack(&track) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}
	if err := t.trackService.UpdateTrack(&id, &track); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, track)
}

// DeleteTrack 	godoc
// @Summary      DeleteTrack
// @Description  delete a track
// @Tags         track
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Delete by Track ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /track/delete/{id} [delete]
func (t *TrackController) DeleteTrack(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := t.trackService.DeleteTrack(&id); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "track deleted"})
}

// GetTrack 	godoc
// @Summary      GetTrack
// @Description  Get a track
// @Tags         track
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Find by Track ID"
// @Success      200  {object}   models.Track
// @Router       /track/get/{id} [get]
func (t *TrackController) FindTrack(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	track, err := t.trackService.FindTrack(&id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, track)
}
func (t *TrackController) RegisterTrackRouter(rt *gin.RouterGroup) {
	router := rt.Group("/track")
	router.POST("/create", t.CreateTrack)
	router.GET("/getAll", t.GetTracks)
	router.PUT("/update/:id", t.UpdateTrack)
	router.DELETE("/delete/:id", t.DeleteTrack)
	router.GET("/get/:id", t.FindTrack)
}
