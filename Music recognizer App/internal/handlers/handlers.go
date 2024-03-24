package handlers

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/models"
	"awesomeProject/internal/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand/v2"
	"net/http"
)

func HomeHandler(ec echo.Context) error {
	return ec.HTML(http.StatusOK, "<h1>Welcome to Home page!</h1>")
}

func UploadHandler(ec echo.Context) error {
	email := ec.FormValue("email")
	file, err := ec.FormFile("file")
	reID := rand.IntN(100000)
	if err != nil {
		database.Insert(models.Request{RequestID: reID, Email: email, Status: "failure", SongID: ""})
		return ec.HTML(http.StatusServiceUnavailable, "<h1>Something bad happened!"+err.Error()+"</h1>")
	}
	err = utils.UploadFile(file, fmt.Sprintf("%d", reID)+".mp3")
	if err != nil {
		database.Insert(models.Request{RequestID: reID, Email: email, Status: "failure", SongID: ""})
		return ec.HTML(http.StatusInternalServerError, "<h1>Something bad happened!"+err.Error()+"</h1>")
	}
	database.Insert(models.Request{RequestID: reID, Email: email, Status: "pending", SongID: ""})
	if err != nil {
		database.Insert(models.Request{RequestID: reID, Email: email, Status: "failure", SongID: ""})
		return ec.HTML(http.StatusInternalServerError, "<h1>Something bad happened!"+err.Error()+"</h1>")
	}
	err = utils.Publish(fmt.Sprintf("%d", reID))
	if err != nil {
		database.Insert(models.Request{RequestID: reID, Email: email, Status: "failure", SongID: ""})
		return ec.HTML(http.StatusInternalServerError, "<h1>Something bad happened!"+err.Error()+"</h1>")
	}
	return ec.HTML(http.StatusOK, "<h1>your request has been successfully answered!</h1>")
}
