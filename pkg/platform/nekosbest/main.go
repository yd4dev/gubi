package nekosbest

import (
	"encoding/json"
	"net/http"
)

//go:generate go run genCategories.go

const baseURL = "https://nekos.best/api/v2/"

var client = &http.Client{}

func FetchGIF(category gifCategory) (GIF, error) {
	req, err := http.NewRequest("GET", baseURL+string(category), nil)
	if err != nil {
		return GIF{}, err
	}
	req.Header.Set("User-Agent", "gubi-go/1.0 (github.com/yd4dev/gubi)")

	resp, err := client.Do(req)
	if err != nil {
		return GIF{}, err
	}
	defer resp.Body.Close()

	var result GIFResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GIF{}, err
	}

	return result.Results[0], nil
}

func FetchImage(category imageCategory) (Image, error) {
	req, err := http.NewRequest("GET", baseURL+string(category), nil)
	if err != nil {
		return Image{}, err
	}
	req.Header.Set("User-Agent", "gubi-go/1.0 (github.com/yd4dev/gubi)")

	resp, err := client.Do(req)
	if err != nil {
		return Image{}, err
	}
	defer resp.Body.Close()

	var result ImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Image{}, err
	}

	return result.Results[0], nil
}
