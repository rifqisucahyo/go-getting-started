package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type ResCheckout struct {
	Product []Cart    `json:"data"`
	Detail  ResDetail `json:"detail"`
}
type ResDetail struct {
	Disabled          bool   `json:"disabled"`
	PaymentMethod     int64  `json:"payment_method"`
	PaymentMethodName string `json:"payment_method_name"`
	Qty               int64  `json:"qty"`
	ShippingFee       int64  `json:"shipping_fee"`
	Weight            int64  `json:"weight"`
	SubPrice          int64  `json:"sub_price"`
	TotalPrice        int64  `json:"total_price"`
	RedirectLink      string `json:"redirect"`
}

type Order struct {
	OrderID          int64     `gorm:"primary_key;unique" json:"order_id" db:"order_id"`
	UserID           uint32    `gorm:"size:255" json:"user_id"`
	Name             string    `gorm:"size:255;not null" json:"name"`
	Phone            string    `gorm:"size:14;not null" json:"phone"`
	Email            string    `gorm:"size:120;" json:"email"`
	SubdistrictID    int64     `gorm:"size:11;not null" json:"subdistrict_id"`
	DistrictID       int64     `gorm:"size:11;not null" json:"district_id"`
	ProvinceID       int64     `gorm:"size:11;not null" json:"province_id"`
	Address          string    `json:"address"`
	PostalCode       int64     `gorm:"size:6" json:"postal_code"`
	CourierID        int64     `gorm:"size:2;not null" json:"courier_id"`
	Voucher          string    `gorm:"size:100" json:"voucher"`
	DiscAmount       int64     `gorm:"size:11" json:"disc_amount"`
	ShippingFee      int64     `gorm:"size:11" json:"shipping_fee"`
	TransactionTotal int64     `gorm:"size:11" json:"transaction_total"`
	OrderstatID      int64     `gorm:"size:11" json:"orderstat_id"`
	PaymentID        int64     `gorm:"size:11" json:"payment_id"`
	VANumber         string    `gorm:"size:11" json:"va_number"`
	CreatedBy        string    `gorm:"size:20" json:"created_by"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy        string    `gorm:"size:20" json:"updated_by"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
type OrderItems struct {
	ID              int64     `gorm:"primary_key;auto_increment;unique" json:"id" db:"orderitem_id"`
	OrderID         int64     `gorm:"size:11" json:"order_id"`
	ProductID       int64     `gorm:"size:11" json:"product_id"`
	Price           int64     `gorm:"size:11" json:"price"`
	Weight          int64     `gorm:"size:11" json:"weight"`
	Qty             int64     `gorm:"size:11" json:"qty"`
	SubPrice        int64     `gorm:"size:11" json:"sub_price"`
	OrderstatitemID int64     `gorm:"size:11" json:"orderstatitem_id"`
	CreatedBy       string    `gorm:"size:20" json:"created_by"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy       string    `gorm:"size:20" json:"updated_by"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type OrderProfile struct {
	// UserID uint32 `gorm:"size:255" json:"user_id"`
	Name  string `gorm:"size:255;not null" json:"name"`
	Phone string `gorm:"size:14;not null" json:"phone"`
	Email string `gorm:"size:120;" json:"email"`
}
type OrderAddress struct {
	SubdistrictID   int64  `gorm:"size:11;not null" json:"subdistrict_id"`
	DistrictID      int64  `gorm:"size:11;not null" json:"district_id"`
	ProvinceID      int64  `gorm:"size:11;not null" json:"province_id"`
	SubdistrictName string `gorm:"size:100;not null" json:"subdistrict_name"`
	DistrictName    string `gorm:"size:100;not null" json:"district_name"`
	ProvinceName    string `gorm:"size:100;not null" json:"province_name"`
	Address         string `json:"address"`
	PostalCode      int64  `gorm:"size:6" json:"postal_code"`
}
type OrderShipping struct {
	CourierID      int64  `gorm:"size:2;not null" json:"id"`
	CourierCode    string `gorm:"size:2;not null" json:"code"`
	CourierName    string `gorm:"size:2;not null" json:"name"`
	CourierService string `gorm:"size:2;not null" json:"service"`
	CourierDesc    string `gorm:"size:2;not null" json:"desc"`
	ShippingFee    int64  `gorm:"size:11" json:"shipping_fee"`
	CourierResi    string `gorm:"size:2;not null" json:"resi"`
}
type OrderDiscount struct {
	Voucher    string `gorm:"size:100" json:"voucher"`
	DiscAmount int64  `gorm:"size:11" json:"value"`
}
type OrderPayment struct {
	PaymentID       int64  `gorm:"size:11" json:"id"`
	PaymentCode     string `gorm:"size:20" json:"code"`
	PaymentType     string `gorm:"size:20" json:"type"`
	PaymentIcon     string `gorm:"size:200" json:"icon"`
	PaymentName     string `gorm:"size:200" json:"name"`
	VANumber        string `gorm:"size:11" json:"va_number"`
	PaymentDeadline string `gorm:"size:11" json:"payment_deadline"`
}

type OrderDetail struct {
	OrderID          int64              `gorm:"primary_key;unique" json:"order_id" db:"order_id"`
	OrderstatID      int64              `gorm:"size:11" json:"orderstat_id"`
	OrderstatName    string             `gorm:"size:20" json:"order_status"`
	Profile          OrderProfile       `json:"profile"`
	Product          []OrderDetailItems `json:"products"`
	Address          OrderAddress       `json:"address"`
	Shipping         OrderShipping      `json:"courier"`
	Discount         OrderDiscount      `json:"discount"`
	Payment          OrderPayment       `json:"payment"`
	SubPriceTotal    int64              `gorm:"size:11" json:"sub_price_total"`
	TransactionTotal int64              `gorm:"size:11" json:"transaction_total"`
	CreatedBy        string             `gorm:"size:20" json:"created_by"`
	CreatedAt        string             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy        string             `gorm:"size:20" json:"updated_by"`
	UpdatedAt        string             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
type OrderDetailItems struct {
	ID              int64  `gorm:"primary_key;auto_increment;unique" json:"id" db:"orderitem_id"`
	OrderID         int64  `gorm:"size:11" json:"order_id"`
	ProductID       int64  `gorm:"size:11" json:"product_id"`
	Price           int64  `gorm:"size:11" json:"price"`
	Weight          int64  `gorm:"size:11" json:"weight"`
	Qty             int64  `gorm:"size:11" json:"qty"`
	OrderstatitemID int64  `gorm:"size:11" json:"orderstatitem_id"`
	StatusName      string `gorm:"size:11" json:"status_name"`
	Name            string `gorm:"size:255" json:"name"`
	Img             string `gorm:"size:200" json:"img"`
	Link            string `gorm:"size:200" json:"link"`
	SubPrice        int64  `gorm:"size:200" json:"sub_price"`
}

type OrderDetailData struct {
	OrderID          int64 `gorm:"primary_key;unique" json:"order_id" db:"order_id"`
	OrderstatID      int64 `gorm:"size:11" json:"orderstat_id"`
	OrderstatName    string
	UserID           uint32 `gorm:"size:255" json:"user_id"`
	Name             string `gorm:"size:255;not null" json:"name"`
	Phone            string `gorm:"size:14;not null" json:"phone"`
	Email            string `gorm:"size:120;" json:"email"`
	SubdistrictID    int64  `gorm:"size:11;not null" json:"subdistrict_id"`
	DistrictID       int64  `gorm:"size:11;not null" json:"district_id"`
	ProvinceID       int64  `gorm:"size:11;not null" json:"province_id"`
	SubdistrictName  string `gorm:"size:100;not null" json:"subdistrict_name"`
	DistrictName     string `gorm:"size:100;not null" json:"district_name"`
	ProvinceName     string `gorm:"size:100;not null" json:"province_name"`
	Address          string `json:"address"`
	PostalCode       int64  `gorm:"size:6" json:"postal_code"`
	CourierID        int64  `gorm:"size:2;not null" json:"courier_id"`
	CourierCode      string `gorm:"size:2;not null" json:"courier_code"`
	CourierService   string `gorm:"size:2;not null" json:"courier_service"`
	CourierName      string `gorm:"size:2;not null" json:"courier_name"`
	CourierDesc      string `gorm:"size:2;not null" json:"courier_desc"`
	CourierResi      string `gorm:"size:2;not null" json:"courier_resi"`
	ShippingFee      int64  `gorm:"size:11" json:"shipping_fee"`
	Voucher          string `gorm:"size:100" json:"voucher"`
	DiscAmount       int64  `gorm:"size:11" json:"value"`
	PaymentID        int64  `gorm:"size:11" json:"payment_id"`
	PaymentCode      string `gorm:"size:20" json:"payment_code"`
	PaymentType      string `gorm:"size:20" json:"payment_type"`
	PaymentIcon      string `gorm:"size:200" json:"payment_icon"`
	PaymentName      string `gorm:"size:200" json:"payment_name"`
	VANumber         string `gorm:"size:11" json:"va_number"`
	SubPriceTotal    int64  `gorm:"size:11" json:"sub_price_total"`
	TransactionTotal int64  `gorm:"size:11" json:"transaction_total"`
	PaymentDeadline  string `gorm:"size:11" json:"payment_deadline"`
	CreatedBy        string `gorm:"size:20" json:"created_by"`
	CreatedAt        string `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy        string `gorm:"size:20" json:"updated_by"`
	UpdatedAt        string `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
type MyOrder struct {
	OrderID          int64              `gorm:"primary_key;unique" json:"order_id" db:"order_id"`
	OrderstatID      int64              `gorm:"size:11" json:"orderstat_id"`
	OrderstatName    string             `gorm:"size:20" json:"order_status"`
	Product          []OrderDetailItems `json:"products"`
	TransactionTotal int64              `gorm:"size:11" json:"transaction_total"`
	CreatedAt        string             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        string             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (o *Order) OrderSave(tx *gorm.DB, resCheckout ResCheckout) (*Order, error) {
	var err error
	err = tx.Debug().Create(&o).Error
	if err != nil {
		return &Order{}, err
	}

	orderItems := []OrderItems{}
	for _, row := range resCheckout.Product {
		if row.Checked {
			price := row.Price
			if row.PriceDisc > 0 {
				price = row.PriceDisc
			}

			orderItems = append(orderItems, OrderItems{
				OrderID:         o.OrderID,
				ProductID:       row.ID,
				Price:           price,
				Qty:             row.Qty,
				Weight:          row.Weight,
				SubPrice:        row.SubPrice,
				OrderstatitemID: 1,
				CreatedBy:       o.CreatedBy,
			})

			DeleteCart(tx, o.UserID, row.ID)
		}
	}
	err = OrderItemSave(tx, orderItems)
	if err != nil {
		return &Order{}, err
	}
	return o, nil
}

func OrderItemSave(tx *gorm.DB, orderItems []OrderItems) error {
	fmt.Println("2. Proses Save Item ===========>>")

	var err error
	valueStrings := make([]string, 0, len(orderItems))
	valueArgs := make([]interface{}, 0, len(orderItems)*3)
	for _, post := range orderItems {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, post.OrderID)
		valueArgs = append(valueArgs, post.ProductID)
		valueArgs = append(valueArgs, post.Price)
		valueArgs = append(valueArgs, post.Qty)
		valueArgs = append(valueArgs, post.Weight)
		valueArgs = append(valueArgs, post.SubPrice)
		valueArgs = append(valueArgs, post.OrderstatitemID)
		valueArgs = append(valueArgs, post.CreatedBy)
		valueArgs = append(valueArgs, post.CreatedAt)
	}
	stmt := fmt.Sprintf("INSERT INTO order_items (order_id, product_id, price, qty, weight, sub_price, orderstatitem_id, created_by, created_at) VALUES %s",
		strings.Join(valueStrings, ","))

	err = tx.Debug().Exec(stmt, valueArgs...).Error
	return err
}

func UpdateOrder(db *gorm.DB, OrderID int64, VANumber string) error {
	var err error
	err = db.Debug().Exec("UPDATE orders SET va_number = ? WHERE order_id = ?", VANumber, OrderID).Error
	return err
}

func GetLastOrderID(db *gorm.DB) (int64, error) {
	var err error
	type OrderLast struct {
		Total int64 `db:"total"`
	}
	orderLast := OrderLast{}
	err = db.Debug().Raw("SELECT count(*) as total FROM orders WHERE DATE_FORMAT(created_at, '%Y%m%d') = DATE_FORMAT(now(), '%Y%m%d')").Scan(&orderLast).Error
	return orderLast.Total, err
}

func (o *OrderDetailData) OrderDetail(db *gorm.DB, OrderID int64) (*OrderDetailData, error) {
	var err error
	orderDetailData := OrderDetailData{}
	err = db.Debug().Raw("SELECT a.order_id, a.orderstat_id, b.status_name as orderstat_name, a.user_id, a.name, phone, email, a.subdistrict_id, c.name as subdistrict_name, a.district_id, d.name as district_name, a.province_id, e.name as province_name, a.address, a.postal_code, a.courier_id, f.service as courier_service, f.description as courier_desc, g.code as courier_code, g.name as courier_name, a.resi_code, a.voucher, a.disc_amount, a.payment_id, h.code as payment_code, h.payment_type, h.icon as payment_icon, h.name as payment_name, a.va_number, (SELECT sum(sub_price) FROM order_items WHERE order_id = ?) as sub_price_total, shipping_fee, transaction_total,DATE_FORMAT((a.created_at + INTERVAL 1 DAY), '%M %d, %Y %H:%i:%s') as payment_deadline, DATE_FORMAT(a.created_at, '%Y %M %d %H:%i') as created_at, a.created_by, DATE_FORMAT(a.updated_at, '%Y %M %d %H:%i') as updated_at, a.updated_by FROM orders a JOIN order_status b ON a.orderstat_id = b.orderstat_id JOIN m_subdistrict c ON a.subdistrict_id = c.subdistrict_id JOIN m_district d ON a.district_id = d.district_id JOIN m_province e ON a.province_id = e.province_id JOIN m_couriers f ON a.courier_id = f.courier_id JOIN m_courier_groups g ON f.courierg_id = g.courierg_id JOIN m_payments h ON a.payment_id = h.payment_id WHERE a.order_id = ? ", OrderID, OrderID).Scan(&orderDetailData).Error

	return &orderDetailData, err
}

func (o *OrderDetailItems) OrderItems(db *gorm.DB, OrderID int64) ([]OrderDetailItems, error) {
	var err error
	orderDetailItems := []OrderDetailItems{}
	err = db.Debug().Raw("SELECT a.orderitem_id as id, a.order_id, a.product_id, a.price, a.weight, a.qty, a.sub_price, a.orderstatitem_id, c.status_name, d.name, d.img, CONCAT('/product/', d.slug, '-', d.product_id) as link from order_items a LEFT JOIN orders b ON a.order_id = b.order_id LEFT JOIN order_status_item c ON a.orderstatitem_id = c.orderstatitem_id LEFT JOIN products d ON a.product_id = d.product_id WHERE a.deleted_at is null AND a.order_id = ? ", OrderID).Scan(&orderDetailItems).Error

	return orderDetailItems, err
}

func (val *MyOrder) GetOrders(db *gorm.DB, userID uint32, status int64, q string, start int64, limit int64) ([]MyOrder, error) {
	var err error
	myOrder := []MyOrder{}
	if status > 0 {
		err = db.Debug().Raw(`SELECT order_id, a.orderstat_id, b.status_name as orderstat_name, a.transaction_total, DATE_FORMAT(created_at, '%d %M %y') as created_at, DATE_FORMAT(updated_at, '%Y%m%d') as updated_at FROM orders a LEFT JOIN order_status b ON a.orderstat_id = b.orderstat_id WHERE a.deleted_at is null AND a.user_id = ? AND a.orderstat_id = ? AND a.orderstat_id is not null ORDER BY a.order_id DESC LIMIT ?, ?`, userID, status, start, limit).Scan(&myOrder).Error
	} else {
		err = db.Debug().Raw(`SELECT order_id, a.orderstat_id, b.status_name as orderstat_name, a.transaction_total, DATE_FORMAT(created_at, '%d %M %y') as created_at, DATE_FORMAT(updated_at, '%Y%m%d') as updated_at FROM orders a LEFT JOIN order_status b ON a.orderstat_id = b.orderstat_id WHERE a.deleted_at is null AND a.user_id = ? AND a.orderstat_id is not null ORDER BY a.order_id DESC LIMIT ?, ?`, userID, start, limit).Scan(&myOrder).Error
	}
	return myOrder, err
}
