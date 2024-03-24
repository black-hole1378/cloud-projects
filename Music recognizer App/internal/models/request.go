package models

type Request struct {
	RequestID int
	Email     string
	Status    string
	SongID    string
}

type Trucks struct {
	Counts int    `json:"totalCount"`
	Items  []Item `json:"items"`
}

type Item struct {
	Data map[string]interface{} `json:"data"`
}

type Album struct {
	Track Trucks `json:"tracks"`
}

type Music struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"preview_url"`
}

type RecommendedTrack struct {
	Tracks []Music `json:"tracks"`
}
