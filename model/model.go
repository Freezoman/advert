package model

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Model struct {
	db *gorm.DB
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func getConn() *gorm.DB {
	dsn := "root:12131213@tcp(localhost:3306)/addb?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
		log.Fatalln("Unable to connect to DB")
	}
	return db
}

func Init() *Model {
	db := getConn()
	db.AutoMigrate(&Ad{})
	db.AutoMigrate(&Photo{})

	return &Model{
		db: db,
	}
}

func (m *Model) GetOneAdvert(id int, optFields ...string) Ad {
	var ad Ad
	var photos []Photo
	fields := []string{"title", "price"}
	for i := range optFields {
		if optFields[i] != "photos" {
			fields = append(fields, optFields[i])
		}
	}

	fmt.Println(fields)
	res := m.db.Select(fields).Where(id).First(&ad)
	if res.Error != nil {
		log.Print(res.Error)
		// log.Fatalf("advertisement with id = %b is not found", id)
	}
	if find(optFields, "photos") {
		res = m.db.Select("url").Where(&Photo{AdID: id}).Find(&photos)
		if res.Error != nil {
			log.Print(res.Error)
			// log.Fatalf("photo with adid = %b is not found", id)
		}
		ad.Photos = photos
	}
	return ad
}

func (m *Model) GetAdverts(sortField string, order string) []Ad {
	var ads []Ad
	res := m.db.Order(sortField+" "+order).Select("title", "price").Find(&ads)
	if res.Error != nil {
		log.Print(res.Error)
	}

	for id := range ads {
		var photos []Photo
		res := m.db.Select("url").Find(&photos)
		if res.Error != nil {
			log.Print(res.Error)
		}
		ads[id].Photos = photos
	}
	return ads
}

func (m *Model) InsertAdvert(ad Ad) (uint, error) {
	res := m.db.Create(&ad)
	if res.Error != nil {
		return 0, res.Error
	}
	return ad.ID, nil
}
