package app

import (
	"encoding/json"
	"mime/multipart"
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

func (a *App) CreateAdvert(data multipart.Form) string {
	title := data.Value["title"]
	body := data.Value["body"]
	price := data.Value["price"]
}

func (a *App) createPhoto(file multipart.File, handler *multipart.FileHeader) model.Photo {

}
