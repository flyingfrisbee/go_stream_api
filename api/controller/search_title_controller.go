package controller

import (
	"fmt"
	"go_stream_api/api/common"
	ws "go_stream_api/repository/webscraper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Search anime title
// @Description Will find all anime with matching keywords
// @Tags Anime
// @Produce json
// @Success 200 {object} common.baseResponse{data=[]webscraper.TitleSearchResult}
// @Router /anime/search [get]
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
// @Param keywords query string true "Keywords of anime title" default(Naruto)
func SearchTitleHandler(c *gin.Context) {
	values := c.Request.URL.Query()
	// Remove leading and trailing white spaces
	keywords := strings.Trim(values["keywords"][0], " ")
	if keywords == "" {
		common.WrapWithBaseResponse(c, nil, "Keyword must at least have characters other than whitespace", http.StatusBadRequest)
		return
	}
	result := ws.ScrapeAnimeTitlesByKeyword(keywords)

	msg := fmt.Sprintf("Success finding anime titles with related keywords: %s", keywords)
	common.WrapWithBaseResponse(c, result, msg, http.StatusOK)
}
