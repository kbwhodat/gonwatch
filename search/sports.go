package search

import (
	"encoding/json"
	"net/http"

	"gonwatch/common"

	// "fmt"
	"io"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

type SportResponse []struct {
	Id	        string	`json:"id"`
	Name	    string	`json:"name"`
}
func ListSports() []common.SportsGenreTypeList {

	url := "https://streami.su/api/sports"

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

	var result SportResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot marshall the json")
	}

	var s common.SportsGenreTypeList

	myList := []common.SportsGenreTypeList{}
	for _, item := range result {
		s.SportsGenreID   = item.Id
		s.SportsGenreName = item.Name
		s.SportsType      = "sports"

		myList = append(myList, s)
	}

	return myList
}
