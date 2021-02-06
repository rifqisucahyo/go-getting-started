package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
)

func (s *Server) OrderDetail(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["orderID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// userID, _ := auth.ExtractTokenID(r)

	orderDetailData := models.OrderDetailData{}
	orderData, err := orderDetailData.OrderDetail(s.DB, orderID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	orderItemData := models.OrderDetailItems{}
	orderItem, err := orderItemData.OrderItems(s.DB, orderID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	orderDetail := models.OrderDetail{}
	orderDetail.OrderID = orderData.OrderID
	orderDetail.OrderstatID = orderData.OrderstatID
	orderDetail.OrderstatName = orderData.OrderstatName
	orderDetail.Profile.Name = orderData.Name
	orderDetail.Profile.Email = orderData.Email
	orderDetail.Profile.Phone = orderData.Phone
	orderDetail.Profile.Phone = orderData.Phone
	orderDetail.Product = orderItem
	orderDetail.Address.SubdistrictID = orderData.SubdistrictID
	orderDetail.Address.DistrictID = orderData.DistrictID
	orderDetail.Address.SubdistrictName = orderData.SubdistrictName
	orderDetail.Address.DistrictName = orderData.DistrictName
	orderDetail.Address.ProvinceID = orderData.ProvinceID
	orderDetail.Address.ProvinceName = orderData.ProvinceName
	orderDetail.Address.Address = orderData.Address
	orderDetail.Address.PostalCode = orderData.PostalCode
	orderDetail.Shipping.CourierID = orderData.CourierID
	orderDetail.Shipping.CourierService = orderData.CourierService
	orderDetail.Shipping.CourierDesc = orderData.CourierDesc
	orderDetail.Shipping.CourierCode = orderData.CourierCode
	orderDetail.Shipping.CourierName = orderData.CourierName
	orderDetail.Shipping.CourierResi = orderData.CourierResi
	orderDetail.Shipping.ShippingFee = orderData.ShippingFee
	orderDetail.Discount.Voucher = orderData.Voucher
	orderDetail.Discount.DiscAmount = orderData.DiscAmount
	orderDetail.Payment.PaymentID = orderData.PaymentID
	orderDetail.Payment.PaymentCode = orderData.PaymentCode
	orderDetail.Payment.PaymentType = orderData.PaymentType
	orderDetail.Payment.PaymentIcon = orderData.PaymentIcon
	orderDetail.Payment.PaymentName = orderData.PaymentName
	orderDetail.Payment.VANumber = orderData.VANumber
	orderDetail.Payment.PaymentDeadline = orderData.PaymentDeadline
	orderDetail.SubPriceTotal = orderData.SubPriceTotal
	orderDetail.TransactionTotal = orderData.TransactionTotal
	orderDetail.CreatedAt = orderData.CreatedAt
	orderDetail.CreatedBy = orderData.CreatedBy
	orderDetail.UpdatedAt = orderData.UpdatedAt
	orderDetail.UpdatedBy = orderData.UpdatedBy

	responses.JSON(w, http.StatusOK, orderDetail)
}

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Start, Limit, Status int64
		Q                    string
	}
	_start, _ := strconv.ParseInt(r.FormValue("start")[0:], 10, 64)
	_limit, _ := strconv.ParseInt(r.FormValue("limit")[0:], 10, 64)
	_status, _ := strconv.ParseInt(r.FormValue("status")[0:], 10, 64)

	par := param{
		Start:  _start,
		Limit:  _limit,
		Status: _status,
		Q:      r.FormValue("q"),
	}

	userID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	orders := models.MyOrder{}
	_orders, err := orders.GetOrders(s.DB, userID, par.Status, par.Q, par.Start, par.Limit)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	for i, row := range _orders {
		orderItemData := models.OrderDetailItems{}
		orderItem, err := orderItemData.OrderItems(s.DB, row.OrderID)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		_orders[i].Product = orderItem

	}

	result := map[string]interface{}{
		"data":   _orders,
		"status": true,
	}

	responses.JSON(w, http.StatusOK, result)
}
