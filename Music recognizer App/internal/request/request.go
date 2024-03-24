package request

import (
	"awesomeProject/internal/models"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func SendRequestShazam(audio io.ReadCloser) (string, error) {
	body, writer, err := initialize(audio)
	req, _ := http.NewRequest("POST", os.Getenv("SHAZAM_API"), body)
	req.Header.Set("content-type", writer.FormDataContentType())
	req.Header.Set("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Set("X-RapidAPI-Host", os.Getenv("RAPID_API_HOST"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer Close(res.Body)
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return getMusic(string(responseBody))
}

func SendRequestSpotify(trackTitle string) (string, error) {
	url := "https://spotify23.p.rapidapi.com/search/?q=" + trackTitle + "&type=track&numberOfTopResults=3"
	result, err := initializeSpotify(url)
	if err != nil {
		return "", err
	}
	return getSongID(result)
}

func SendRecommendedTracks(songID string) ([]models.Music, error) {
	url := "https://spotify23.p.rapidapi.com/recommendations/?limit=20&seed_tracks=" + songID
	result, err := initializeSpotify(url)
	if err != nil {
		return nil, err
	}
	return getRecommendedTracks(result)
}
