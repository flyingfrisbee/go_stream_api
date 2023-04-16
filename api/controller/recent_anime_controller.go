package controller

import (
	"fmt"
	"go_stream_api/api/common"
	db "go_stream_api/repository/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get recently updated / released anime
// @Description This endpoint supports pagination
// @Tags Anime
// @Produce json
// @Success 200 {object} common.baseResponse{data=[]domain.RecentAnime}
// @Router /anime/recent [get]
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
// @Param page query int false "Desired page" default(1)
func RecentAnimeHandler(c *gin.Context) {
	values := c.Request.URL.Query()
	page, err := strconv.Atoi(values["page"][0])
	if err != nil {
		page = 1
	}
	recentAnimes, err := db.Conn.Pg.GetRecentAnimes(page)
	if err != nil {
		msg := fmt.Sprintf("Error occured when fetching recent anime page %d", page)
		common.WrapWithBaseResponse(c, nil, msg, http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Success fetching recent anime page %d", page)
	common.WrapWithBaseResponse(c, recentAnimes, msg, http.StatusOK)
}
