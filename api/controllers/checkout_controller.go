package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
)

type ReqAddress struct {
	ProvinceID    int64  `json:"province_id"`
	DistrictID    int64  `json:"district_id"`
	SubdistrictID int64  `json:"subdistrict_id"`
	Address       string `json:"address"`
	PostalCode    int64  `json:"postal_code"`
}
type ReqProduct struct {
	Id      int64 `json:"id"`
	Price   int64 `json:"price"`
	Qty     int64 `json:"qty"`
	Checked bool  `json:"checked"`
	Delete  bool  `json:"delete"`
}
type ReqCheckout struct {
	Step     int64        `json:"step"`
	Name     string       `json:"name"`
	Email    string       `json:"email"`
	Phone    string       `json:"phone"`
	Voucher  string       `json:"voucher"`
	Courier  int64        `json:"courier"`
	Payment  int64        `json:"payment"`
	Address  ReqAddress   `json:"address"`
	Products []ReqProduct `json:"products"`
}

func checkoutProcess(server *Server, userID uint32, reqProduct []ReqProduct) error {
	var err error
	var totalDelete, totalUpdate, totalCreate int
	for _, row := range reqProduct {
		if row.Delete {
			err := models.DeleteCart(server.DB, userID, row.Id)
			if err != nil {
				panic(err)
			}
			totalDelete += 1
		} else {
			if userID != 0 {
				total, err := models.GetCartCount(server.DB, userID, row.Id)
				if err != nil {
					panic(err)
				}
				if total > 0 {
					err := models.UpdateCart(server.DB, userID, row.Id, row.Qty)
					if err != nil {
						panic(err)
					}
					totalUpdate += 1
				} else {
					err := models.CreateCart(server.DB, userID, row.Id, row.Qty)
					if err != nil {
						panic(err)
					}
					totalCreate += 1
				}
			}
		}
	}

	return err
}

func getDataCheckout(server *Server, userID uint32, reqProduct []ReqProduct) (models.ResCheckout, error) {

	var total_qty, total_price_product, total_weight int64
	var resCheckout models.ResCheckout
	var productIDList []int64
	for _, row := range reqProduct {
		if !row.Delete {
			productIDList = append(productIDList, row.Id)
		}
	}

	carts := models.Cart{}
	products, err := carts.GetCart(server.DB, int32(userID), productIDList)

	if err != nil {
		return resCheckout, err
	}

	for i, row := range products {
		for _, rowReq := range reqProduct {
			if row.ID == rowReq.Id {
				qty := InterfaceToInt64(rowReq.Qty)
				weight := InterfaceToInt64(row.Weight) * qty
				sub_price := InterfaceToInt64(0)
				if row.Stock > 0 {
					sub_price = InterfaceToInt64(row.Price) * qty
				}

				products[i].Qty = qty
				products[i].SubPrice = sub_price
				products[i].Disc = 0
				products[i].Checked = rowReq.Checked
				products[i].Weight = weight
				if row.PriceDisc != 0 {
					if row.Stock > 0 {
						sub_price = InterfaceToInt64(row.PriceDisc) * qty
					}

					products[i].SubPrice = sub_price
					products[i].Disc = math.Ceil(((InterfaceToFloat64(row.Price) - InterfaceToFloat64(row.PriceDisc)) / InterfaceToFloat64(row.Price)) * 100)
				}

				error_msg := make([]string, 0, 0)
				if row.Stock <= 0 {
					error_msg = append(error_msg, "Status produk habis")
				}
				if rowReq.Price != row.Price {
					// error_msg = append(error_msg, "Terjadi perubahan harga")
				}

				products[i].ErrorMsg = error_msg

				if rowReq.Checked && !rowReq.Delete {
					total_price_product += sub_price
					total_qty += qty
					total_weight += weight
				}
			}
		}
	}

	resCheckout.Product = products
	resCheckout.Detail.Qty = total_qty
	resCheckout.Detail.Weight = total_weight
	resCheckout.Detail.SubPrice = total_price_product
	resCheckout.Detail.TotalPrice = total_price_product
	resCheckout.Detail.PaymentMethod = 0
	resCheckout.Detail.PaymentMethodName = ""
	resCheckout.Detail.Disabled = true

	return resCheckout, nil
}

func (s *Server) saveCheckout(userID uint32, reqCheckout ReqCheckout, resCheckout models.ResCheckout) (*models.Order, error) {
	fmt.Println("1. Proses Save Order ===========>>")
	tx := s.DB.Begin()
	lastOrderID, err := models.GetLastOrderID(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	t := time.Now()
	date := fmt.Sprintf("%02d%02d%02d", t.Year(), t.Month(), t.Day())
	codeUnix := date + fmt.Sprint((lastOrderID + 1))
	_codeUnix, _ := strconv.ParseInt(codeUnix, 10, 64)

	order := models.Order{}
	order.OrderID = _codeUnix
	order.UserID = userID
	order.Name = reqCheckout.Name
	order.Phone = reqCheckout.Phone
	order.Email = reqCheckout.Email
	order.SubdistrictID = reqCheckout.Address.SubdistrictID
	order.DistrictID = reqCheckout.Address.DistrictID
	order.ProvinceID = reqCheckout.Address.ProvinceID
	order.Address = reqCheckout.Address.Address
	order.PostalCode = reqCheckout.Address.PostalCode
	order.CourierID = reqCheckout.Courier
	order.Voucher = reqCheckout.Voucher
	order.DiscAmount = 0
	order.ShippingFee = resCheckout.Detail.ShippingFee
	order.TransactionTotal = resCheckout.Detail.TotalPrice
	order.OrderstatID = 1
	order.PaymentID = reqCheckout.Payment
	if userID == 0 {
		order.CreatedBy = substr(reqCheckout.Name, 0, 20)
	} else {
		order.CreatedBy = fmt.Sprint(userID)
	}

	orderData, err := order.OrderSave(tx, resCheckout)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	respPayment, err := s.paymentCreate(reqCheckout, resCheckout, orderData.OrderID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	fmt.Println("4. Proses Update VA Number =======>>.")

	//Jika menggunakan VA
	if len(respPayment.VANumbers) > 0 {
		for _, row := range respPayment.VANumbers {
			if row.VANumber != "" {
				err = models.UpdateOrder(tx, orderData.OrderID, row.VANumber)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				break
			}
		}
	}

	tx.Commit()
	return orderData, nil
}

func (server *Server) Checkout(w http.ResponseWriter, r *http.Request) {
	var err error
	userID, _ := auth.ExtractTokenID(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var reqCheckout ReqCheckout
	json.Unmarshal([]byte(body), &reqCheckout)

	err = checkoutProcess(server, userID, reqCheckout.Products)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Proses checkout gagal."))
		return
	}

	result, err := getDataCheckout(server, userID, reqCheckout.Products)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if reqCheckout.Courier != 0 {
		_couriers, courier_code, err := server.getCourier(reqCheckout.Courier)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		_, shippingFee, err := shipmentCheck(fmt.Sprint(result.Detail.Weight), fmt.Sprint(reqCheckout.Address.SubdistrictID), courier_code, _couriers, "checkout")
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		result.Detail.ShippingFee = shippingFee
		result.Detail.TotalPrice = result.Detail.TotalPrice + result.Detail.ShippingFee

		if result.Detail.ShippingFee > 0 && result.Detail.Qty > 0 && result.Detail.SubPrice > 0 && result.Detail.TotalPrice > 0 && reqCheckout.Payment != 0 {
			result.Detail.Disabled = false
		}

		if reqCheckout.Step == 2 {
			respCheckout, err := server.saveCheckout(userID, reqCheckout, result)
			if err != nil {
				responses.ERROR(w, http.StatusBadRequest, err)
				return
			}
			result.Detail.RedirectLink = os.Getenv("HOST_WEB") + `/order/complete?id=` + fmt.Sprint(respCheckout.OrderID)
		}

	}

	responses.JSON(w, http.StatusOK, result)
}
