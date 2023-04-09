package controller

import (
	"go_stream_api/api/common"
	db "go_stream_api/repository/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type syncAnimeRequest struct {
	AnimeIDs []int `json:"anime_ids" binding:"required"`
}

// @Summary Sync anime
// @Description To sync anime data from client to server
// @Tags Anime
// @Produce json
// @Success 200 {object} common.baseResponse{data=[]domain.SyncAnime}
// @Router /anime/sync [post]
// @Param request body controller.syncAnimeRequest true "request body"
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
func SyncAnimeHandler(c *gin.Context) {
	var request syncAnimeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusBadRequest)
		return
	}

	animes, err := db.Conn.Pg.GetSyncAnime(request.AnimeIDs)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	common.WrapWithBaseResponse(c, animes, "Success getting sync data", http.StatusOK)
}
