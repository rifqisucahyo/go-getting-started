package models

import (
	"github.com/jinzhu/gorm"
)

type Verification struct {
	VerificationID int64
	Timer          int64
}

func CheckPhoneOTP(db *gorm.DB, phone string) (int, error) {
	var err error
	var total int
	db.Raw("SELECT verification_id FROM verifications WHERE type = 'PHONE' AND phone = ?", phone).Count(&total)
	return total, err
}

func PhoneOTP(db *gorm.DB, phone string) (Verification, error) {
	var err error
	verification := Verification{}
	err = db.Debug().Raw("SELECT verification_id, TIMESTAMPDIFF(SECOND, now(), expires) as timer FROM verifications WHERE phone = ?", phone).Scan(&verification).Error
	return verification, err
}

func CreatePhoneOTP(db *gorm.DB, code int64, phone string) error {
	var err error
	err = db.Debug().Exec("INSERT INTO verifications (type, expires, code, phone) VALUES ('PHONE', (now() + INTERVAL 1 MINUTE), ?, ?)", code, phone).Error
	return err
}

func UpdatePhoneOTP(db *gorm.DB, code int64, phone string) error {
	var err error
	err = db.Debug().Exec("UPDATE verifications SET expires = (now() + INTERVAL 1 MINUTE), code = ? WHERE type = 'PHONE' AND phone = ?", code, phone).Error
	return err
}

func PhoneOTPVerification(db *gorm.DB, code int64, phone string) (Verification, error) {
	var err error
	verification := Verification{}
	err = db.Debug().Raw("SELECT verification_id, TIMESTAMPDIFF(SECOND, now(), expires) as timer FROM verifications WHERE phone = ? AND code = ?", phone, code).Scan(&verification).Error
	return verification, err
}
