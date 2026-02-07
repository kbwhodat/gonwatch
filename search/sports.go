package search

import (
	"encoding/json"
	"net/http"

	"gonwatch/common"

	"io"
	"log"
)

type SportResponse []struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ListSports() []common.SportsGenreTypeList {

	url := string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 115, 116, 114, 101, 97, 109, 105, 46, 115, 117, 47, 97, 112, 105, 47, 115, 112, 111, 114, 116, 115})

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

	myList := []common.SportsGenreTypeList{}
	for _, item := range result {
		var s common.SportsGenreTypeList
		s.SportsGenreID = item.Id
		s.SportsGenreName = item.Name
		s.SportsType = "sports"

		myList = append(myList, s)
	}

	return myList
}
