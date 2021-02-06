package models

import (
	"github.com/jinzhu/gorm"
)

type InfoBox struct {
	ID      uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name    string `gorm:"size:255;not null;unique" json:"name"`
	Img     string `gorm:"size:14" json:"img"`
	Content string `gorm:"size:100;not null;unique" json:"content"`
}

func (val *InfoBox) GetInfoBox(db *gorm.DB) (interface{}, error) {
	var err error
	result := []InfoBox{}
	db.Raw("SELECT infoboxID as id, name, img, content FROM infobox WHERE is_publish = ?", 1).Limit(6).Scan(&result)

	return &result, err
}

type Slider struct {
	ID   uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
	Img  string `gorm:"size:100" json:"img"`
	Link string `gorm:"size:255;not null;unique" json:"link"`
}

func (val *Slider) GetSlider(db *gorm.DB) (interface{}, error) {
	var err error
	result := []Slider{}
	db.Raw("SELECT sliderID as id, name, img, link FROM slider WHERE is_publish = ?", 1).Limit(5).Scan(&result)

	return &result, err
}
