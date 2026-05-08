package nekosbest

type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type GIF struct {
	URL        string     `json:"url"`
	Dimensions Dimensions `json:"dimensions"`
	AnimeName  string     `json:"anime_name"`
}

type Image struct {
	URL        string     `json:"url"`
	Dimensions Dimensions `json:"dimensions"`
	ArtistName string     `json:"artist_name"`
	ArtistHREF string     `json:"artist_href"`
	SourceURL  string     `json:"source_url"`
}

type GIFResponse struct {
	Results []GIF `json:"results"`
}

type ImageResponse struct {
	Results []Image `json:"results"`
}
