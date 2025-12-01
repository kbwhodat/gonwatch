package search

import (
	"encoding/json"
	"net/http"
	"strings"

	"gonwatch/common"

	// "fmt"
	"io"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

type MatchesSportsResponse []struct {
	Id	        string	 `json:"id"`
	Name	    string	 `json:"title"`
	Sources     []struct {
		SourceName string `json:"source"`
		Source_id string `json:"id"`
	} `json:"sources"`
}
type EmbedUrlResponse []struct {
	EmbdUrl string `json:"embedUrl"`
	Viewers int64  `json:"viewers"`

}
func ListSportMatches(sport string) []common.SportsGenreTypeList {

	url := "https://streami.su/api/matches/" + sport
	log.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

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

	var result MatchesSportsResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot marshall the json")
	}

	var s common.SportsGenreTypeList

	myList := []common.SportsGenreTypeList{}
	for _, item := range result {
		s.SportsGenreID    = item.Id
		s.SportsGenreName  = item.Name

		// initializing sports struct to be used
		var sportSources []struct {
			SportsSourceName string
			SportsSourceId string
		}

		// using anonymous struct here. Named struct would be an easier implementation as I would not have to redefine the shape of the struct
		for _, sources := range item.Sources {
			sportSources = append(sportSources, struct {
				SportsSourceName string
				SportsSourceId string
			}{
				SportsSourceName: sources.SourceName,
				SportsSourceId: sources.Source_id,
			})
		}

		s.SportsType   = "sports"
		s.SportSources = sportSources

		myList = append(myList, s)
	}

	return myList
}

func ListStreams(streams []string) []common.SportsGenreTypeList {
	var s common.SportsGenreTypeList

	myList := []common.SportsGenreTypeList{}
	for _, item := range streams {
		key, value, _ := strings.Cut(item, ":")
		log.Println("key:" + key)
		log.Println("value:" + value)
		s.SportsGenreID    = value
		s.SportsGenreName  = key

		s.SportsType   = "streams"
		myList = append(myList, s)
	}

	return myList
}

func GetStreamLink(stream_id string, stream_path string) string {
	url := "https://streami.su/api/stream/" + stream_id + "/" + stream_path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

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

	var result EmbedUrlResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot marshall the json")
	}


	var embedUrl string
	for _, item := range result {
		embedUrl = item.EmbdUrl
		break
	}

	return embedUrl

}
