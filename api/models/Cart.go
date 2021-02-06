package models

import (
	"os"

	"github.com/jinzhu/gorm"
)

type Cart struct {
	ID            int64    `gorm:"primary_key;auto_increment" json:"id"`
	Name          string   `gorm:"size:255;not null" json:"name"`
	Img           string   `gorm:"size:100" json:"img"`
	Link          string   `gorm:"size:255" json:"link"`
	StatusProduct string   `gorm:"size:11" json:"status_product"`
	Qty           int64    `gorm:"size:11" json:"qty"`
	Stock         int64    `gorm:"size:11" json:"stock"`
	Weight        int64    `gorm:"size:11" json:"weight"`
	Price         int64    `gorm:"size:11" json:"price"`
	PriceDisc     int64    `gorm:"size:11" json:"price_disc"`
	Disc          float64  `gorm:"size:11" json:"disc"`
	SubPrice      int64    `gorm:"size:11" json:"sub_price"`
	Checked       bool     `gorm:"size:11" json:"checked"`
	ErrorMsg      []string `gorm:"size:11" json:"error_msg"`
}

func (val *Cart) GetCart(db *gorm.DB, userID int32, productID []int64) ([]Cart, error) {
	var err error
	carts := []Cart{}
	db.Raw("SELECT a.product_id as id, name, price, price_disc, (CASE WHEN b.qty is not null THEN b.qty ELSE 0 END) as qty, a.stock, a.weight, (CASE WHEN b.qty is not null THEN (price*b.qty) ELSE 0 END) as sub_price, CASE WHEN stock > 0 THEN 'Ready' ELSE 'Oos' END as status_product, img, CONCAT('"+os.Getenv("HOST_WEB")+"', '/product/', slug, '-', a.product_id) as link FROM products a LEFT JOIN carts b ON a.product_id = b.product_id WHERE a.product_id IN (?) OR b.user_id = ? GROUP by b.product_id ORDER BY b.cart_id DESC, a.name ASC", productID, userID).Scan(&carts)
	return carts, err
}

func GetCartCount(db *gorm.DB, userID uint32, productID int64) (int, error) {
	var err error
	var total int
	db.Raw("SELECT cart_id FROM carts WHERE user_id = ? AND product_id = ?", userID, productID).Count(&total)
	return total, err
}

func DeleteCart(db *gorm.DB, userID uint32, productID int64) error {
	var err error
	db.Exec("DELETE FROM carts WHERE user_id = ? AND product_id = ?", userID, productID)
	return err
}

func UpdateCart(db *gorm.DB, userID uint32, productID int64, qty int64) error {
	var err error
	db.Exec("UPDATE carts SET qty = ? WHERE user_id = ? AND product_id = ?", qty, userID, productID)
	return err
}

func CreateCart(db *gorm.DB, userID uint32, productID int64, qty int64) error {
	var err error
	db.Exec("INSERT INTO carts (qty, user_id, product_id) VALUES (?,?,?)", qty, userID, productID)
	return err
}
