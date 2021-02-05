package model

import (
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model `json:"-"`
	Url        string `json:"url"`
	AdID       int    `json:"-"`
}

type Ad struct {
	gorm.Model `json:"-"`
	Title      string  `json:"title"`
	Body       string  `json:"body,omitempty"`
	Price      int     `json:"price"`
	Photos     []Photo `json:"photos,omitempty"`
}
