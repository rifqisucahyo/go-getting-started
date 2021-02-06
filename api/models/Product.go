package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID            int32   `gorm:"primary_key;auto_increment" json:"id"`
	SkuCode       string  `gorm:"size:255;not null;unique" json:"sku_code"`
	Name          string  `gorm:"size:255;not null;unique" json:"name"`
	Price         int     `gorm:"size:11" json:"price"`
	PriceDisc     int     `gorm:"size:11" json:"price_disc"`
	Discount      float64 `gorm:"size:11" json:"discount"`
	PriceGrosir   int     `gorm:"size:11" json:"price_grosir"`
	Stock         int     `gorm:"size:11" json:"stock"`
	StatusProduct string  `gorm:"size:11" json:"status_product"`
	Img           string  `gorm:"size:100" json:"img"`
	Description   string  `gorm:"type:TEXT" json:"description"`
	Weight        int     `gorm:"size:11" json:"weight"`
	IsPo          bool    `gorm:"size:1" json:"is_po"`
	PoTime        int     `gorm:"size:2" json:"po_time"`
	IsWarranty    bool    `gorm:"size:1" json:"is_warranty"`
	IsNew         bool    `gorm:"size:1" json:"is_new"`
	Link          string  `gorm:"size:255" json:"link"`
	Favorite      bool    `gorm:"size:255" json:"favorite"`
	Rating        int32   `gorm:"size:11" json:"rating"`
	Review        int32   `gorm:"size:11" json:"review"`
	ProductCond   string  `gorm:"size:11" json:"condition"`
	Brand         string  `gorm:"size:11" json:"brand"`
	Sold          string  `gorm:"size:11" json:"sold"`
}

type ReqFilterProduct struct {
	StatusProduct int    `json:"status_product"`
	ProductCond   string `json:"condition"`
	Type, Q       string
	IsNew         bool `json:"is_new"`
	IsBest        bool `json:"is_best"`
}
type ReqSortingProduct struct {
	Lates      string `json:"latest"`
	BestSeller string `json:"bestseller"`
	Price      string `json:"price"`
}
type ReqProduct struct {
	Start, Limit int
	Filter       ReqFilterProduct
	Sorting      ReqSortingProduct
}

func (val *Product) GetProduct(db *gorm.DB, start int64, limit int64, type_product string) ([]Product, error) {
	var err error
	result := []Product{}
	row := db.Raw("SELECT product_id as id, sku_code, name, img, price, price_disc, 100-ROUND((price_disc / price)* 100, 0) as discount, price_grosir, stock, CASE WHEN stock = 2 THEN 'Ready' ELSE 'Oos' END as status_product, description, weight, is_po, po_time, is_warranty, is_new, CONCAT('/product/', slug, '-', product_id) as link, 1 as favorite, 2 as rating, 123 as review, cond AS product_cond FROM products WHERE is_publish = ? LIMIT ?, ?", 1, start, limit).Scan(&result)
	if type_product == "new" {
		row.Where("is_new = ?", 1)
	} else if type_product == "best" {
		row.Where("is_best = ?", 1)
	}

	return result, err
}

func (val *Product) GetProducts(db *gorm.DB, par ReqProduct) (interface{}, error) {
	var err error
	result := []Product{}
	_sqlWhere := ``
	if par.Filter.StatusProduct > 0 {
		_sqlWhere += ` AND status_product = ` + fmt.Sprint(par.Filter.StatusProduct)
	}
	if par.Filter.ProductCond != "" {
		_sqlWhere += ` AND cond = '` + fmt.Sprint(par.Filter.ProductCond) + `'`
	}
	if par.Filter.IsBest == true {
		_sqlWhere += ` AND is_best = 1`
	}
	if par.Filter.IsNew == true {
		_sqlWhere += ` AND is_new = 1`
	}

	_sqlOrder := ``
	if par.Sorting.Lates != "" {
		_sqlOrder += ` a.product_id ` + par.Sorting.Lates
	}
	if par.Sorting.BestSeller != "" {
		coma := ``
		if _sqlOrder != "" {
			coma = `, `
		}
		_sqlOrder += coma + ` a.sold ` + par.Sorting.BestSeller
	}
	if par.Sorting.Price != "" {
		coma := ``
		if _sqlOrder != "" {
			coma = `, `
		}
		_sqlOrder += coma + ` (CASE WHEN a.price_disc is not null THEN a.price_disc ELSE a.price END) ` + par.Sorting.Price
	}

	if _sqlOrder != "" {
		_sqlOrder = ` ORDER BY ` + _sqlOrder
	}

	err = db.Debug().Raw(`SELECT product_id as id, sku_code, name, img, price, price_disc, 100-ROUND((price_disc / price)* 100, 0) as discount, price_grosir, stock, CASE WHEN stock = 2 THEN 'Ready' ELSE 'Oos' END as status_product, description, weight, is_po, po_time, is_warranty, is_new, CONCAT('/product/', slug, '-', product_id) as link, 1 as favorite, 2 as rating, 123 as review, cond AS product_cond FROM products a WHERE is_publish = ? AND CONCAT(name, sku_code, slug, status_product) like '%`+par.Filter.Q+`%' `+_sqlWhere+_sqlOrder, 1).Offset(par.Start).Limit(par.Limit).Scan(&result).Error
	return &result, err
}

func (p *Product) ProductByID(db *gorm.DB, productID int32) (result Product, err error) {
	err = db.Debug().Model(Product{}).Where("product_id = ?", productID).Take(&p).Error
	if err != nil {
		return result, err
	}

	db.Raw("SELECT a.product_id as id, a.sku_code, a.name, a.img, a.price, (CASE WHEN a.cond = 'NEW' THEN 'Baru' ELSE 'Bekas' END) AS product_cond, a.price_disc, a.price_grosir, a.stock, CASE WHEN a.status_product = 2 THEN 'Ready' ELSE 'Oos' END as status_product, a.description, a.weight, a.is_po, a.po_time, a.is_warranty, a.is_new, CONCAT('/product/',slug, '-', a.product_id) as link, 1 as favorite, 2 as rating, 123 as review, a.cond AS product_cond, a.sold, b.name as brand FROM products a LEFT JOIN brands b ON a.brand_id = b.brand_id WHERE a.product_id = ? AND a.is_publish = ?", productID, 1).Take(&result)
	return result, err
}

func GalleryByID(db *gorm.DB, productID int32) (interface{}, error) {
	var err error
	type Gallery struct {
		Path string `gorm:"size:255" json:"path"`
	}
	result := []Gallery{}
	db.Raw("SELECT path FROM product_images WHERE product_id = ? ", productID).Scan(&result)
	return result, err
}
