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
)

func GetSeries(text string) []common.StreamTypeList {

	urlPrefix := []byte{104, 116, 116, 112, 115, 58, 47, 47, 97, 112, 105, 46, 116, 104, 101, 109, 111, 118, 105, 101, 100, 98, 46, 111, 114, 103, 47, 51, 47, 115, 101, 97, 114, 99, 104, 47, 116, 118, 63, 113, 117, 101, 114, 121, 61}
	urlSuffix := []byte{38, 105, 110, 99, 108, 117, 100, 101, 95, 97, 100, 117, 108, 116, 61, 116, 114, 117, 101, 38, 108, 97, 110, 103, 117, 97, 103, 101, 61, 101, 110, 45, 85, 83, 38, 112, 97, 103, 101, 61, 49}
	url := string(urlPrefix) + url.QueryEscape(text) + string(urlSuffix)

	bearerBytes := []byte{66, 101, 97, 114, 101, 114, 32, 101, 121, 74, 104, 98, 71, 99, 105, 79, 105, 74, 73, 85, 122, 73, 49, 78, 105, 74, 57, 46, 101, 121, 74, 104, 100, 87, 81, 105, 79, 105, 73, 119, 77, 122, 77, 48, 77, 71, 73, 48, 90, 68, 104, 107, 79, 68, 103, 53, 78, 68, 77, 120, 77, 122, 73, 52, 89, 50, 69, 119, 79, 68, 81, 48, 89, 122, 73, 51, 90, 106, 65, 51, 90, 105, 73, 115, 73, 109, 53, 105, 90, 105, 73, 54, 77, 84, 99, 120, 77, 122, 73, 122, 78, 106, 73, 120, 77, 67, 52, 53, 78, 122, 99, 115, 73, 110, 78, 49, 89, 105, 73, 54, 73, 106, 89, 50, 77, 87, 82, 108, 79, 71, 89, 121, 78, 84, 73, 52, 89, 106, 74, 108, 77, 68, 69, 48, 89, 84, 78, 108, 78, 84, 100, 109, 89, 121, 73, 115, 73, 110, 78, 106, 98, 51, 66, 108, 99, 121, 73, 54, 87, 121, 74, 104, 99, 71, 108, 102, 99, 109, 86, 104, 90, 67, 74, 100, 76, 67, 74, 50, 90, 88, 74, 122, 97, 87, 57, 117, 73, 106, 111, 120, 102, 81, 46, 100, 49, 122, 95, 101, 55, 122, 54, 105, 118, 76, 84, 50, 65, 49, 115, 75, 45, 101, 95, 98, 75, 98, 76, 119, 108, 71, 82, 112, 83, 71, 55, 102, 80, 57, 74, 81, 73, 55, 115, 69, 97, 111}
	var bearer = string(bearerBytes)

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

	var result TvResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot unmarshall the json")
	}

	var s common.StreamTypeList

	myList := []common.StreamTypeList{}
	for _, item := range result.Results {
		s.StreamID = item.Id
		s.StreamPlot = item.Overview
		s.StreamTitle = item.Title
		s.StreamYear = item.ReleaseDate
		s.StreamRating = item.Rating
		if len(item.OriginCountry) > 0 {
			s.StreamCountry = item.OriginCountry[0]
		}

		myList = append(myList, s)
	}

	return myList
}
