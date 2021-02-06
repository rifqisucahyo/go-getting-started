package models

import (
	"github.com/jinzhu/gorm"
)

type Region struct {
	ProvinceID      int64  `gorm:"size:11;not null" db:"province_id" json:"province_id"`
	Provincename    string `gorm:"size:100;not null" json:"province_name"`
	DistrictID      int64  `gorm:"size:11;not null" json:"district_id"`
	Districtname    string `gorm:"size:100;not null" json:"district_name"`
	SubdistrictID   int64  `gorm:"size:11;not null" json:"subdistrict_id"`
	Subdistrictname string `gorm:"size:100;not null" json:"subdistrict_name"`
	Region          string `json:"region"`
}

func GetRegion(db *gorm.DB, q string, start int64, limit int64) (*[]Region, error) {
	var err error
	region := []Region{}
	db.Raw(`SELECT c.province_id, c.name as provincename, b.district_id, b.name as districtname, a.subdistrict_id, a.name as subdistrictname, CONCAT(a.name, ', ',b.name, ', ',c.name) as region, b.postal_code as postalcode FROM m_subdistrict a LEFT JOIN m_district b ON a.district_id = b.district_id LEFT JOIN m_province c ON b.province_id = c.province_id WHERE (a.name like '%`+q+`%' OR b.name like '%`+q+`%' OR c.name like '%`+q+`%') LIMIT ?, ?`, start, limit).Scan(&region)
	return &region, err
}
