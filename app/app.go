package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"test/model"
)

type App struct {
	model IModel
}

type IModel interface {
	GetOneAdvert(id int, fields ...string) model.Ad
	GetAdverts(sortField string, order string) []model.Ad
	InsertAdvert(ad model.Ad) (uint, error)
}
type response struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func Init(m IModel) *App {
	return &App{model: m}
}

func (a *App) GetJsonAdvert(idStr string, optFields ...string) string {

	if strings.HasPrefix(idStr, "id") {
		idStr = idStr[2:]
	} else {
		return "Link is incomplete"
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err.Error()
	}
	availableOptFields := []string{"body", "photos"}
	for _, v := range optFields {
		if !find(availableOptFields, v) {
			return "Only 'body' and 'photos' fields are available"
		}
	}
	ad := a.model.GetOneAdvert(id, optFields...)

	b, err := json.Marshal(ad)
	if err != nil {
		panic(err)
	}

	return string(b)
}
func (a *App) GetJsonAdverts(sortField string, order string) string {
	availableSortFields := []string{"price", "created_at"}
	availableOrder := []string{"asc", "desc"}

	if !find(availableSortFields, sortField) {
		return "Only 'price' and 'created' fields are available for sorting"
	}
	if !find(availableOrder, order) {
		return "Only 'asc' and 'desc' orders are available"
	}

	ads := a.model.GetAdverts(sortField, order)
	b, err := json.Marshal(ads)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (a *App) CreateAdvert(data *multipart.Form) (string, error) {
	var photos []model.Photo

	title := data.Value["title"][0]
	body := data.Value["body"][0]
	price, err := strconv.Atoi(data.Value["price"][0])
	if err != nil {
		return "Can't create an advert", err
	}

	ad := model.Ad{Title: title, Body: body, Price: price}

	files := data.File["myFile"]
	for i, _ := range files {
		photo, err := a.createPhoto(files[i])
		if err != nil {
			fmt.Println("Cannot create file")
			return "Cannot create file", err
		}
		photos = append(photos, photo)
	}
	ad.Photos = photos
	res, err := a.model.InsertAdvert(ad)
	if err != nil {
		return "Cannot save advert", err
	}
	response := response{Id: fmt.Sprint(res), Status: "ok"}

	resp, err := json.Marshal(response)
	if err != nil {
		return "Cannot form json", err
	}

	return string(resp), nil

}

func (a *App) createPhoto(handler *multipart.FileHeader) (model.Photo, error) {
	extMap := map[string]string{"image/jpeg": "jpg", "image/gif": "gif", "image/png": "png"}
	var photo model.Photo

	file, err := handler.Open()
	defer file.Close()
	if err != nil {
		return photo, err
	}

	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		os.Mkdir("./files", os.ModeDir)
	}

	tempFile, err := ioutil.TempFile("files", "upload-*."+extMap[handler.Header.Get("Content-Type")])
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	photo.Url = "http://localhost:8080/" + filepath.ToSlash(tempFile.Name())
	return photo, nil
}
