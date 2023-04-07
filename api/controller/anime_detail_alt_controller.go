package controller

import (
	"fmt"
	"go_stream_api/api/common"
	ws "go_stream_api/repository/webscraper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Anime detail alternative
// @Description This endpoint will scrape the detail instead of fetching from database
// @Tags Anime
// @Produce json
// @Success 200 {object} common.baseResponse{data=webscraper.anime}
// @Router /anime/detail-alt [post]
// @Param request body webscraper.TitleSearchResult true "request body"
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
func AnimeDetailAlternativeHandler(c *gin.Context) {
	var request ws.TitleSearchResult
	err := c.ShouldBindJSON(&request)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, "Cannot parse request body", http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("Success fetching anime detail: %s", request.Title)
	common.WrapWithBaseResponse(c, ws.ScrapeDetailAlternative(request), msg, http.StatusOK)
}
