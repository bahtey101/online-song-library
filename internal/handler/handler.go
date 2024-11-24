package handler

import (
	"net/http"
	"online-song-library/internal/service"
	"strconv"

	msong "online-song-library/internal/model/song"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	songs := router.Group("/songs")
	{
		songs.GET("/", handler.GetPaginatedSongs)
		songs.GET("/{id}", handler.GetPaginatedText)
		songs.POST("/", handler.CreateSong)
		songs.PUT("/{id}", handler.UpdateSong)
		songs.DELETE("/{id}", handler.DeleteSong)
	}

	return router
}

func (handler *Handler) GetPaginatedSongs(ctx *gin.Context) {
	logrus.Debug("GetPaginatedSongs: received request")

	fields := map[string]string{}

	group, ok := ctx.GetQuery("group")
	if ok {
		fields["group"] = group
	}

	song, ok := ctx.GetQuery("song")
	if ok {
		fields["song"] = song
	}

	releaseDate, ok := ctx.GetQuery("releaseDate")
	if ok {
		fields["releaseDate"] = releaseDate
	}

	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "1"))
	if offset < 1 {
		offset = 1
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	logrus.Debugf("GetPaginatedSongs: fields=%v, offset=%d, limit=%d",
		fields, offset, limit)

	songs, err := handler.service.GetPaginatedSongs(ctx, fields, offset, limit)
	if err != nil {
		logrus.Errorf("GetPaginatedSongs: failed to fetch songs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch songs",
		})
		return
	}

	logrus.Infof("GetPaginatedSongs: retrieved %d songs", len(*songs))

	ctx.JSON(http.StatusOK, gin.H{
		"songs": songs,
	})
}

func (handler *Handler) GetPaginatedText(ctx *gin.Context) {
	logrus.Debug("GetPaginatedText: received request")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		logrus.Error("GetPaginatedText: invalid song ID")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid song ID",
		})
	}

	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "1"))
	if offset < 1 {
		offset = 1
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	logrus.Debugf("GetPaginatedText: song ID=%d, offset=%d, limit=%d",
		id, offset, limit)

	text, err := handler.service.GetPaginatedText(
		ctx,
		msong.Song{ID: id},
		offset,
		limit,
	)
	if err != nil {
		logrus.Errorf("GetPaginatedText: failed to fetch text: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch text",
		})
		return
	}

	logrus.Infof("GetPaginatedText: successfully fetched text for song ID=%d", id)

	ctx.JSON(http.StatusOK, gin.H{
		"text": text,
	})
}

func (handler *Handler) CreateSong(ctx *gin.Context) {
	logrus.Debug("CreateSong: received request")

	var request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		logrus.Error("CreateSong: invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := handler.service.CreateSong(ctx, msong.Song{
		Group: request.Group,
		Song:  request.Song,
	}); err != nil {
		logrus.Errorf("CreateSong: failed to create song, error=%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create song",
		})
		return
	}

	logrus.Infof("CreateSong: successfully created song")
	ctx.JSON(http.StatusCreated, "")
}

func (handler *Handler) UpdateSong(ctx *gin.Context) {
	logrus.Debug("UpdateSong: received request")

	var request struct {
		Group       string `json:"group"`
		Song        string `json:"song"`
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		logrus.Error("UpdateSong: invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		logrus.Error("UpdateSong: invalid ID parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid song ID",
		})
		return
	}

	err = handler.service.UpdateSong(ctx, msong.Song{
		ID:          id,
		Group:       request.Group,
		Song:        request.Song,
		ReleaseDate: request.ReleaseDate,
		Verses:      msong.SplitIntoVerses(request.Text),
		Link:        request.Link,
	})
	if err != nil {
		logrus.Errorf("UpdateSong: failed to update song ID=%d, error=%v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update song",
		})
		return
	}

	logrus.Infof("UpdateSong: successfully updated song ID=%d", id)

	ctx.JSON(http.StatusNoContent, "")
}

func (handler *Handler) DeleteSong(ctx *gin.Context) {
	logrus.Debug("DeleteSong: received request")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		logrus.Error("DeleteSong: invalid ID parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid song ID",
		})
		return
	}

	if err := handler.service.DeleteSong(ctx, msong.Song{ID: id}); err != nil {
		logrus.Errorf("DeleteSong: failed to delete song: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete song",
		})
		return
	}

	logrus.Infof("DeleteSong: successfully deleted song ID=%d", id)

	ctx.JSON(http.StatusNoContent, "")
}
