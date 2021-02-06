package models

import (
	"github.com/jinzhu/gorm"
)

type Courier struct {
	Id          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Service     string `json:"service"`
	Description string `json:"description"`
}

func (val *Courier) GetCourier(db *gorm.DB, courierID int64) (*[]Courier, error) {
	var err error
	courier := []Courier{}
	if courierID != 0 {
		db.Raw("SELECT a.courier_id as id, b.code, b.name, a.service, a.service FROM m_couriers a LEFT JOIN m_courier_groups b ON a.courierg_id = b.courierg_id WHERE a.is_publish = 1 AND b.is_publish = 1 AND a.courier_id = ?", courierID).Scan(&courier)
	} else {
		db.Raw("SELECT a.courier_id as id, b.code, b.name, a.service, a.service FROM m_couriers a LEFT JOIN m_courier_groups b ON a.courierg_id = b.courierg_id WHERE a.is_publish = 1 AND b.is_publish = 1").Scan(&courier)
	}
	return &courier, err
}

func (val *Courier) GetCourierByID(db *gorm.DB, courierID int64) (*[]Courier, error) {
	var err error
	courier := []Courier{}
	db.Raw("SELECT a.courier_id as id, b.code, b.name, a.service, a.service FROM m_couriers a LEFT JOIN m_courier_groups b ON a.courierg_id = b.courierg_id WHERE a.is_publish = 1 AND b.is_publish = 1 AND courier_id = ?", courierID).Scan(&courier)
	return &courier, err
}
