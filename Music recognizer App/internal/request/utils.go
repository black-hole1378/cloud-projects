package request

import (
	"awesomeProject/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func Close(reader io.ReadCloser) {
	err := reader.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func initialize(audio io.ReadCloser) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload_file", "music.mp3")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return nil, writer, err
	}

	_, err = io.Copy(part, audio)
	if err != nil {
		fmt.Println("Error copying file data:", err)
		return nil, writer, err
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return nil, writer, err
	}

	return body, writer, nil
}

func initializeSpotify(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("SPOTIFY_API_HOST"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer Close(res.Body)
	body, _ := io.ReadAll(res.Body)
	return string(body), nil
}

func getMusic(myJson string) (string, error) {
	fmt.Println(myJson)
	var trackData map[string]interface{}
	err := json.Unmarshal([]byte(myJson), &trackData)
	if err != nil {
		return "", err
	}
	track := trackData["track"].(map[string]interface{})
	title := track["subtitle"].(string)
	return title, nil
}

func getSongID(myJson string) (string, error) {
	var album models.Album
	if err := json.Unmarshal([]byte(myJson), &album); err != nil {
		return "", err
	}
	fmt.Println("hello", album.Track.Items[0])
	return album.Track.Items[0].Data["id"].(string), nil
}

func getRecommendedTracks(myJson string) ([]models.Music, error) {
	var tracks models.RecommendedTrack
	if err := json.Unmarshal([]byte(myJson), &tracks); err != nil {
		return nil, err
	}
	fmt.Println(tracks.Tracks)
	return tracks.Tracks, nil
}
