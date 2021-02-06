package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Phone     string    `gorm:"size:14;unique" json:"phone"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Address struct {
	ID              int64  `gorm:"primary_key;auto_increment;not null;unique" json:"id"`
	Label           string `gorm:"size:255" json:"label"`
	Name            string `gorm:"size:100;not null" json:"name"`
	Phone           string `gorm:"size:14" json:"phone"`
	Address         string `gorm:"size:255" json:"address"`
	ProvinceID      int64  `gorm:"size:11;not null" json:"province_id"`
	ProvinceName    string `gorm:"size:100" json:"province_name"`
	DistrictID      int64  `gorm:"size:11;not null" json:"district_id"`
	DistrictName    string `gorm:"size:100" json:"district_name"`
	SubdistrictID   int64  `gorm:"size:11;not null" json:"subdistrict_id"`
	SubdistrictName string `gorm:"size:100" json:"subdistrict_name"`
	Region          string `json:"region"`
	PostalCode      int64  `gorm:"size:6" json:"postal_code"`
	IsPrimary       bool   `gorm:"size:1;not null" json:"is_primary"`
}

type UserAddress struct {
	ID            int64     `gorm:"primary_key;auto_increment;not null;unique" json:"id"`
	UserID        uint32    `gorm:"size:11" json:"user_id"`
	Label         string    `gorm:"size:255" json:"label"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Phone         string    `gorm:"size:14" json:"phone"`
	Address       string    `gorm:"size:255" json:"address"`
	ProvinceID    int64     `gorm:"size:11;not null" json:"province_id"`
	DistrictID    int64     `gorm:"size:11;not null" json:"district_id"`
	SubdistrictID int64     `gorm:"size:11;not null" json:"subdistrict_id"`
	Postal_code   int64     `gorm:"size:6" json:"postal_code"`
	IsPrimary     bool      `gorm:"size:1;not null" json:"is_primary"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	if u.Password != "" && u.Username != "" && u.Email != "" {
		hashedPassword, err := Hash(u.Password)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Phone = html.EscapeString(strings.TrimSpace(u.Phone))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Name == "" {
			return errors.New("Required Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Name == "" {
			return errors.New("Required Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (val *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	user := []User{}
	db.Raw("SELECT a.user_id as id, a.username, a.name, a.email, a.phone, date_format(a.created_at, '%d %M %Y %H:%i') as created_at, date_format(a.updated_at, '%d %M %Y %H:%i') as updated_at FROM users a").Scan(&user)
	return &user, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("userID = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("userID = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"username":  u.Username,
			"name":      u.Name,
			"email":     u.Email,
			"phone":     u.Phone,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("user_id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("user_id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) FindUserByEmail(db *gorm.DB, email string) (interface{}, error) {
	var err error
	result := User{}

	db.Raw("SELECT user_id as id, username, name, email FROM users WHERE email = ?", email).Scan(&result)
	return &result, err
}

func GetAddress(db *gorm.DB, userID uint32, q string, start int64, limit int64) (*[]Address, error) {
	var err error
	address := []Address{}
	// fmt.Println(`SELECT a.uaddress_id as id, a.label, a.name, a.phone, a.address, a.province_id, b.name as province_name, a.district_id, c.name as district_name, a.subdistrict_id, d.name as subdistrict_name, a.is_primary, date_format(a.created_at, '%d %M %Y %H:%i') as created_at, date_format(a.updated_at, '%d %M %Y %H:%i') as updated_at, a.postal_code as postalcode, CONCAT(a.name, ', ',b.name, ', ',c.name) as region FROM user_addresses a JOIN m_province b ON a.province_id = b.province_id JOIN m_district c ON a.district_id = c.district_id JOIN m_subdistrict d ON a.subdistrict_id = d.subdistrict_id WHERE a.user_id = ` + fmt.Sprint(userID) + ` AND (a.label like '%` + q + `%' OR a.name like '%` + q + `%' OR a.phone like '%` + q + `%' OR a.address like '%` + q + `%' OR b.name like '%` + q + `%' OR c.name like '%` + q + `%' OR d.name like '%` + q + `%') ORDER BY a.is_primary DESC, a.updated_at DESC, a.uaddress_id DESC LIMIT ` + fmt.Sprint(start) + `, ` + fmt.Sprint(limit))
	db.Raw(`SELECT a.uaddress_id as id, a.label, a.name, a.phone, a.address, a.province_id, b.name as province_name, a.district_id, c.name as district_name, a.subdistrict_id, d.name as subdistrict_name, a.is_primary, date_format(a.created_at, '%d %M %Y %H:%i') as created_at, date_format(a.updated_at, '%d %M %Y %H:%i') as updated_at, a.postal_code as postalcode, CONCAT(a.name, ', ',b.name, ', ',c.name) as region FROM user_addresses a JOIN m_province b ON a.province_id = b.province_id JOIN m_district c ON a.district_id = c.district_id JOIN m_subdistrict d ON a.subdistrict_id = d.subdistrict_id WHERE a.user_id = ? AND (a.label like '%`+q+`%' OR a.name like '%`+q+`%' OR a.phone like '%`+q+`%' OR a.address like '%`+q+`%' OR b.name like '%`+q+`%' OR c.name like '%`+q+`%' OR d.name like '%`+q+`%') ORDER BY a.is_primary DESC, a.updated_at DESC, a.uaddress_id DESC LIMIT ?, ?`, userID, start, limit).Scan(&address)
	return &address, err
}

func (a *UserAddress) Prepare(userID uint32) {
	a.ID = 0
	a.UserID = userID
	a.Label = html.EscapeString(strings.TrimSpace(a.Label))
	a.Name = html.EscapeString(strings.TrimSpace(a.Name))
	a.Phone = html.EscapeString(strings.TrimSpace(a.Phone))
	a.Address = html.EscapeString(strings.TrimSpace(a.Address))
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *UserAddress) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if a.Name == "" {
			return errors.New("Nama wajib diisi")
		}
		return nil
	default:
		if a.Name == "" {
			return errors.New("Nama wajib diisi")
		}
		return nil
	}
}

func UpdateAddressPrimary(db *gorm.DB, userID uint32) error {
	var err error
	db.Exec("UPDATE user_addresses SET is_primary = 0 WHERE is_primary = 1 AND user_id = ?", userID)
	return err
}

func (a *UserAddress) SaveAddress(db *gorm.DB) (*UserAddress, error) {

	var err error
	err = db.Debug().Create(&a).Error
	if err != nil {
		return &UserAddress{}, err
	}
	return a, nil
}

func (u *UserAddress) UpdateAddress(db *gorm.DB, uid uint32, uaddressID int64) (*UserAddress, error) {
	db = db.Debug().Model(&UserAddress{}).Where("user_id = ?", uid).Where("uaddress_id = ?", uaddressID).Take(&UserAddress{}).UpdateColumns(
		map[string]interface{}{
			"label":         u.Label,
			"name":          u.Name,
			"phone":         u.Phone,
			"address":       u.Address,
			"provinceID":    u.ProvinceID,
			"districtID":    u.DistrictID,
			"subdistrictID": u.SubdistrictID,
			"postal_code":   u.Postal_code,
			"is_primary":    u.IsPrimary,
			"update_at":     time.Now(),
		},
	)
	if db.Error != nil {
		return &UserAddress{}, db.Error
	}
	return u, nil
}

func (u *UserAddress) DelAddress(db *gorm.DB, uid uint32, uaddressID int64) (int64, error) {

	db = db.Debug().Model(&UserAddress{}).Where("user_id = ?", uid).Where("uaddress_id = ?", uaddressID).Take(&UserAddress{}).Delete(&UserAddress{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
