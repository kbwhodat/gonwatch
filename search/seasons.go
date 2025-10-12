package search

//http://ky-iptv.com/player_api.php?username=walkerrodney216%40gmail.com&password=Stream1958&action=get_vod_info&vod_id=819792

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gonwatch/common"

	// "fmt"
	"io"
	"log"
	"net/url"

	_ "github.com/marcboeker/go-duckdb"
)

type TheSeasonsDBResponse struct {
	Seasons []struct {
		SeasonId	   int64  `json:"id"`
		SeasonTitle    string `json:"name"`
		SeasonOverview string `json:"overview"`
		EpisodeCount   int64  `json:"episode_count"`
		ReleaseDate	   string `json:"air_date"`
		SeasonNumber   int64 `json:"season_number"`
	}
}
func GetSeasons(tmdbid int64) []common.SeasonsTypeList {

	url := "https://api.themoviedb.org/3/tv/" + url.QueryEscape(strconv.Itoa(int(tmdbid))) + "?language=en-US"
	var bearer = "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIwMzM0MGI0ZDhkODg5NDMxMzI4Y2EwODQ0YzI3ZjA3ZiIsIm5iZiI6MTcxMzIzNjIxMC45NzcsInN1YiI6IjY2MWRlOGYyNTI4YjJlMDE0YTNlNTdmYyIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.d1z_e7z6ivLT2A1sK-e_bKbLwlGRpSG7fP9JQI7sEao"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result TheSeasonsDBResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot marshall the json")
	}

	var s common.SeasonsTypeList

	myList := []common.SeasonsTypeList{}
	for _, item := range result.Seasons {
		if item.SeasonNumber != 0 {
			s.SeriesID          = tmdbid
			s.EpisodeCount      = item.EpisodeCount
			s.SeasonTitle       = item.SeasonTitle
			s.SeasonReleaseDate = item.ReleaseDate
			s.SeasonNumber      = strconv.Itoa(int(item.SeasonNumber))
			// s.SeasonID          = tmdbid

			myList = append(myList, s)
		}
	}

	return myList
}
