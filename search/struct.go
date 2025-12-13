package search

type MovieResponse struct {
	Results []struct {
		Id	          int64	   `json:"id"`
		Title	      string   `json:"title"`
		Overview	  string   `json:"overview"`
		ReleaseDate	  string   `json:"release_date"`
		OriginCountry []string `json:"origin_country,omitempty"`
		Rating        float64 `json:"vote_average"`
	}
}

type TvResponse struct {
	Results []struct {
		Id	          int64	   `json:"id"`
		Title	      string   `json:"name"`
		Overview	  string   `json:"overview"`
		ReleaseDate	  string   `json:"first_air_date"`
		OriginCountry []string `json:"origin_country,omitempty"`
		Rating        float64 `json:"vote_average"`
	}
}

type AnimeResponse struct {
	Results []struct {
		Id	        int64	`json:"id"`
		Title	    string	`json:"name"`
		Overview	string	`json:"overview"`
		ReleaseDate	string	`json:"first_air_date"`
		Country	    string	`json:"original_language"`
		Rating      float64 `json:"vote_average"`
	}
}
