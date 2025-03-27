package handlers

import (
	"net/http"
	"song_library/entities"
	"song_library/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SongHandler struct {
	service *service.Service
	logger  *logrus.Logger
}

func NewSongHandler(logger *logrus.Logger, service *service.Service) *SongHandler {
	return &SongHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary Get songs
// @Description Get songs with filtering and pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group_name query string false "Group name"
// @Param song_name query string false "Song name"
// @Param release_date query string false "Release date"
// @Param link query string false "Link"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} entities.Song
// @Router /api/v1/songs [get]
func (sh *SongHandler) GetSongs(c *gin.Context) {
	group := c.Request.URL.Query().Get("group_name")
	song := c.Request.URL.Query().Get("song_name")
	date := c.Request.URL.Query().Get("release_date")
	link := c.Request.URL.Query().Get("link")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		sh.logger.Errorf("Can't convert string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		sh.logger.Errorf("Can't convert string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	filter := entities.Song{GroupName: group, SongName: song, ReleaseDate: date, Link: link}
	songs, err := sh.service.GetSongs(filter, page, limit)
	if err != nil {
		sh.logger.Printf("Error getting songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get song lyrics
// @Description Get song lyrics with pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]string "Lyrics text"
// @Router /api/v1/songs/{id}/text [get]
func (sh *SongHandler) GetSongText(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		sh.logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid song id"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		sh.logger.Errorf("Can't convert string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		sh.logger.Errorf("Can't convert string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	text, err := sh.service.GetSongText(c.Request.Context(), songID, page, limit)
	if err != nil {
		sh.logger.Errorf("Error getting song text: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get song"})
		return
	}

	c.JSON(http.StatusOK, text)
}

// @Summary Delete song
// @Description Delete a song by ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Router /api/v1/songs/{id} [delete]
func (sh *SongHandler) DeleteSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sh.logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid song id"})
		return
	}
	err = sh.service.DeleteSong(c.Request.Context(), songID)
	if err != nil {
		sh.logger.Errorf("Failed to delete song:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Update song
// @Description Update song details
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param song body entities.Song true "Song details"
// @Success 200 {object} entities.Song "Updated song"
// @Router /api/v1/songs/{id} [put]
func (sh *SongHandler) UpdateSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sh.logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid song id"})
		return
	}

	var song entities.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		sh.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	sh.logger.Info(song)

	song.ID = songID
	err = sh.service.UpdateSong(c.Request.Context(), &song)
	if err != nil {
		sh.logger.Errorf("Failed to update song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// @Summary Add song
// @Description Add a new song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body entities.Song true "Song details"
// @Success 201 {object} entities.Song "Added song"
// @Router /api/v1/songs [post]
func (sh *SongHandler) AddSong(c *gin.Context) {
	var song entities.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		sh.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request body"})
		return
	}

	err := sh.service.AddSong(c.Request.Context(), &song)
	if err != nil {
		sh.logger.Errorf("Failed to create song:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song"})
		return
	}

	c.JSON(http.StatusCreated, song)
}
