package models

import (
	"github.com/jinzhu/gorm"
)

type Payment struct {
	Id          int32  `gorm:"primary_key;auto_increment" json:"id" db:"id"`
	PaymentType string `gorm:"size:20;not null" json:"payment_type"`
	Code        string `gorm:"size:20;not null" json:"code"`
	RekNumber   string `gorm:"size:20;not null" json:"rek_number"`
	Name        string `gorm:"size:30;not null" json:"name"`
	Icon        string `gorm:"size:100" json:"icon"`
	Description string `gorm:"size:255" json:"description"`
	IsDisabled  bool   `gorm:"size:255" json:"is_disabled"`
}

func (val *Payment) GetPayments(db *gorm.DB) (*[]Payment, error) {
	var err error
	payments := []Payment{}
	db.Raw("SELECT a.payment_id as id, a.payment_type, a.code, a.rek_number, a.name, a.icon, a.description, a.is_disabled FROM m_payments a WHERE a.is_publish = 1").Scan(&payments)
	return &payments, err
}

func (val *Payment) GetPaymentByID(db *gorm.DB, paymentID int64) (*Payment, error) {
	var err error
	payment := Payment{}

	err = db.Raw("SELECT a.payment_id as id, a.payment_type, a.code, a.rek_number, a.name, a.icon, a.description, a.is_disabled FROM m_payments a WHERE a.is_publish = 1 AND a.payment_id = ?", paymentID).Scan(&payment).Debug().Error
	return &payment, err
}
