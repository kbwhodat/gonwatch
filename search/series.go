package search

//http://ky-iptv.com/player_api.php?username=walkerrodney216%40gmail.com&password=Stream1958&action=get_vod_info&vod_id=819792

import (
	"encoding/json"
	"net/http"

	"gonwatch/common"

	// "fmt"
	"io"
	"log"
	"net/url"

	_ "github.com/marcboeker/go-duckdb"
)

func checkForErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type TheSeriesDBResponse struct {
	Results []struct {
		Id	        int64	`json:"id"`
		Title	    string	`json:"name"`
		Overview	string	`json:"overview"`
		ReleaseDate	string	`json:"first_air_date"`
	}
}
func GetSeries(text string) []common.StreamTypeList {

	url := "https://api.themoviedb.org/3/search/tv?query=" + url.QueryEscape(text) + "&include_adult=true&language=en-US&page=1"
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

	var result TheSeriesDBResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot marshall the json")
	}

	var s common.StreamTypeList

	myList := []common.StreamTypeList{}
	for _, item := range result.Results {
		s.StreamID    = item.Id
		s.StreamPlot  = item.Overview
		s.StreamTitle = item.Title
		s.StreamYear  = item.ReleaseDate

		myList = append(myList, s)
	}

	return myList
}
