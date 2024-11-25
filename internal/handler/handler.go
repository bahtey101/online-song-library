package handler

import (
	"net/http"
	"online-song-library/internal/service"
	"strconv"

	msong "online-song-library/internal/model/song"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

// GetPaginatedSongs godoc
// @Summary Get  paginated list of songs
// @Description  Retrieve a paginated list of songs based on optional query parameters.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group       query   string  false  "Group name filter"
// @Param        song        query   string  false  "Song name filter"
// @Param        releaseDate query   string  false  "Release date filter"
// @Param        offset      query   int     false  "Page offset (default 1)"
// @Param        limit       query   int     false  "Number of items per page (default 10, max 100)"
// @Success      200 {array} song.Song
// @Failure      500 {string} string "failed to fetch songs"
// @Router       /songs/ [get]
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

// GetPaginatedText godoc
// @Summary      Get paginated song text
// @Description  Retrieve paginated text of a song by ID.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id     path   uint64  true   "Song ID"
// @Param        offset query  int     false  "Page offset (default 1)"
// @Param        limit  query  int     false  "Number of items per page (default 4, max 10)"
// @Success      200 {string} string "text"
// @Failure      400 {string} string "invalid song ID"
// @Failure      500 {string} string "failed to fetch text"
// @Router       /songs/{id} [get]
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

// CreateSong godoc
// @Summary      Create a new song
// @Description  Add a new song to the system.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        request body handler.CreateSong.Request true "group name and song name"
// @Success      201 {string} string "Created"
// @Failure      400 {string} string "invalid request body"
// @Failure      500 {string} string "failed to create song"
// @Router       /songs/ [post]
func (handler *Handler) CreateSong(ctx *gin.Context) {
	logrus.Debug("CreateSong: received request")

	type Request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}
	var req Request

	if err := ctx.BindJSON(&req); err != nil {
		logrus.Error("CreateSong: invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := handler.service.CreateSong(ctx, msong.Song{
		Group: req.Group,
		Song:  req.Song,
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

// UpdateSong godoc
// @Summary      Update an existing song
// @Description  Update the details of a song by ID.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id          path uint64 true  "Song ID"
// @Param        request body handler.UpdateSong.Request false "song fields"
// @Success      204 {string} string "No content"
// @Failure      400 {object} string "invalid request body"
// @Failure      500 {object} string "failed to update song"
// @Router       /songs/{id} [put]
func (handler *Handler) UpdateSong(ctx *gin.Context) {
	logrus.Debug("UpdateSong: received request")

	type Request struct {
		Group       string `json:"group"`
		Song        string `json:"song"`
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	var req Request

	if err := ctx.BindJSON(&req); err != nil {
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
		Group:       req.Group,
		Song:        req.Song,
		ReleaseDate: req.ReleaseDate,
		Verses:      msong.SplitIntoVerses(req.Text),
		Link:        req.Link,
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

// DeleteSong godoc
// @Summary      Delete a song
// @Description  Remove a song from the system by ID.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id   path   uint64  true   "Song ID"
// @Success      204 {string} string "No content"
// @Failure      400 {string} string "invalid song ID"
// @Failure      500 {string} string "failed to delete song"
// @Router       /songs/{id} [delete]
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
