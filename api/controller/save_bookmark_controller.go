package controller

import (
	"fmt"
	"go_stream_api/api/common"
	db "go_stream_api/repository/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type saveBookmarkRequest struct {
	AnimeID       int    `json:"anime_id" binding:"required"`
	UserToken     string `json:"user_token" binding:"required"`
	LatestEpisode string `json:"latest_episode" binding:"required"`
}

// @Summary Save bookmark
// @Description User can receive notification from apps when their bookmarked anime has new update
// @Tags Bookmark
// @Produce json
// @Success 200 {object} common.baseResponse
// @Router /bookmark [post]
// @Param request body controller.saveBookmarkRequest true "request body"
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
func SaveBookmarkHandler(c *gin.Context) {
	var request saveBookmarkRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.Conn.Pg.InsertBookmark(request.AnimeID, request.UserToken, request.LatestEpisode)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Success bookmark anime with id: %d", request.AnimeID)
	common.WrapWithBaseResponse(c, nil, msg, http.StatusOK)
}
