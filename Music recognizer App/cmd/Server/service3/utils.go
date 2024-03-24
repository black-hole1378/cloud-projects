package main

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/models"
	packageRequest "awesomeProject/internal/request"
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func execute() {
	requests, err := database.SelectFromDB("failure", "Status")
	if err != nil {
		log.Fatal(err)
	}
	for _, request := range requests {
		if err := sendEmail(request.Email, "templates/failure.html", []models.Music{}); err != nil {
			log.Println(err)
		}
	}
	database.Delete("failure", "Status")
	requests, err = database.SelectFromDB("ready", "Status")
	for _, request := range requests {
		musics, err := packageRequest.SendRecommendedTracks(request.SongID)
		if err != nil {
			log.Println(err)
		}
		database.Update(fmt.Sprintf("%d", request.RequestID), "done", "Status")
		if err := sendEmail(request.Email, "templates/index.html", musics); err != nil {
			log.Println(err)
		}
	}
}

func sendEmail(to string, htmlPath string, musics []models.Music) error {
	var body bytes.Buffer
	te, err := template.ParseFiles(htmlPath)
	if err != nil {
		return err
	}
	if err = te.Execute(&body, musics); err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MY_EMAIL"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Music lists!")
	m.SetBody("text/html", body.String())
	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 587, os.Getenv("MY_EMAIL"), os.Getenv("EMAIL_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
