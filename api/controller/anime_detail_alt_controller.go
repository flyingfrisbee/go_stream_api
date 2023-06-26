package controller

import (
	"fmt"
	"go_stream_api/api/common"
	ws "go_stream_api/repository/webscraper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type animeDetailAltRequest struct {
	Title    string `json:"title" binding:"required" example:"naruto"`
	Endpoint string `json:"endpoint" binding:"required" example:"/category/naruto"`
}

// @Summary Anime detail alternative
// @Description This endpoint will scrape the detail instead of fetching from database
// @Tags Anime
// @Produce json
// @Success 200 {object} common.baseResponse{data=domain.Anime}
// @Router /anime/detail-alt [post]
// @Param request body controller.animeDetailAltRequest true "request body"
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
func AnimeDetailAlternativeHandler(c *gin.Context) {
	defer common.RecoverWhenPanic(c, "Cannot reach scraping server")

	var request animeDetailAltRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("Success fetching anime detail: %s", request.Title)
	common.WrapWithBaseResponse(c, ws.ScrapeDetailAlternative(request.Title, request.Endpoint), msg, http.StatusOK)
}
