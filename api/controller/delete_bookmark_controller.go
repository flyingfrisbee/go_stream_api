package controller

import (
	"fmt"
	"go_stream_api/api/common"
	db "go_stream_api/repository/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteBookmarkRequest struct {
	AnimeID   int    `json:"anime_id" binding:"required"`
	UserToken string `json:"user_token" binding:"required"`
}

// @Summary Delete bookmark
// @Description User can delete bookmark to opt out from notification in the future
// @Tags Bookmark
// @Produce json
// @Success 200 {object} common.baseResponse
// @Router /delete [post]
// @Param request body controller.deleteBookmarkRequest true "request body"
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
func DeleteBookmarkHandler(c *gin.Context) {
	var request deleteBookmarkRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.Conn.Pg.DeleteBookmark(request.AnimeID, request.UserToken)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Success unbookmark anime with id: %d", request.AnimeID)
	common.WrapWithBaseResponse(c, nil, msg, http.StatusOK)
}
