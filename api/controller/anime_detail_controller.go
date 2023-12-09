package controller

import (
	"fmt"
	"go_stream_api/api/common"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	ws "go_stream_api/repository/webscraper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get anime detail
// @Description By providing anime id, this endpoint will return the details
// @Tags Anime
// @Produce json
// @Success 200 {object} common.baseResponse{data=domain.Anime}
// @Router /anime/detail/{anime_id} [get]
// @Param Authorization header string true "Insert your auth token" default(Bearer <Add access token here>)
// @Param anime_id path int true "Anime ID"
func AnimeDetailHandler(c *gin.Context) {
	animeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf("path param {id} has to be of type integer, received: %s", c.Param("id"))
		common.WrapWithBaseResponse(c, nil, msg, http.StatusBadRequest)
		return
	}

	eps, err := db.Conn.Pg.GetEpisodes(animeID)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch from mongoDB instead
	// if len(eps) == 0 {
	// 	eps, err = db.Conn.Mongo.GetEpisodes(animeID)
	// 	if err != nil {
	// 		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	anime, err := db.Conn.Pg.GetAnimeDetail(animeID)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	anime.Episodes = eps

	shouldFetchEpisodesFromScraper := len(eps) == 0 || len(eps) == 30
	if shouldFetchEpisodesFromScraper {
		url := fmt.Sprintf(env.EpisodesURLFormat, anime.ID)
		ws.ScrapeEpisodes(&anime, url)
		ws.ReverseEpisodesOrder(&anime)
	}

	msg := fmt.Sprintf("Success fetching anime detail with id: %d", animeID)
	common.WrapWithBaseResponse(c, anime, msg, http.StatusOK)
}
