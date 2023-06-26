package controller

import (
	"go_stream_api/api/common"
	ws "go_stream_api/repository/webscraper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type videoURLRequest struct {
	EpisodeEndpoint string `json:"episode_endpoint" binding:"required" example:"/naruto-episode-1"`
}

type videoURLResponse struct {
	VideoURL string `json:"video_url"`
}

// @Summary Video URL
// @Description Fetch video URL by episode endpoint
// @Tags Anime
// @Produce json
// @Success 200 {object} common.baseResponse{data=controller.videoURLResponse}
// @Router /anime/video-url [post]
// @Param request body controller.videoURLRequest true "request body"
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
func VideoURLHandler(c *gin.Context) {
	defer common.RecoverWhenPanic(c, "Cannot reach scraping server")

	var request videoURLRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, "Cannot parse request body", http.StatusBadRequest)
		return
	}

	response := videoURLResponse{
		VideoURL: ws.ScrapeVideoURL(request.EpisodeEndpoint),
	}
	common.WrapWithBaseResponse(c, response, "Success fetching video url", http.StatusOK)
}
