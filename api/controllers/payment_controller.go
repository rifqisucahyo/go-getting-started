package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
)

type TrxTransactionDetails struct {
	GrossAmount int64  `json:"gross_amount"`
	OrderID     string `json:"order_id"`
}

type TrxCustomerDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type TrxItemDetails struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
	Name     string `json:"name"`
}

type TrxBankTransfer struct {
	Bank     string                 `json:"bank"`
	VANumber string                 `json:"va_number"`
	FreeText map[string]interface{} `json:"free_text"`
}

type TrxCheckout struct {
	PaymentType        string                `json:"payment_type"`
	TransactionDetails TrxTransactionDetails `json:"transaction_details"`
	CustomerDetails    TrxCustomerDetails    `json:"customer_details"`
	ItemDetails        []TrxItemDetails      `json:"item_details"`
	BankTransfer       TrxBankTransfer       `json:"bank_transfer"`
}

type RespVA struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}
type RespPayment struct {
	StatusCode        string   `json:"status_code"`
	StatusMessage     string   `json:"status_message"`
	TransactionID     string   `json:"transaction_id"`
	IrderID           string   `json:"order_id"`
	MerchantID        string   `json:"merchant_id"`
	GrossAmount       string   `json:"gross_amount"`
	Currency          string   `json:"currency"`
	PaymentType       string   `json:"payment_type"`
	TransactionTime   string   `json:"transaction_time"`
	TransactionStatus string   `json:"transaction_status"`
	VANumbers         []RespVA `json:"va_numbers"`
	FraudStatus       string   `json:"fraud_status"`
}

func (server *Server) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments := models.Payment{}
	data, err := payments.GetPayments(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, data)
}

func (s *Server) paymentCreate(reqCheckout ReqCheckout, resCheckout models.ResCheckout, orderID int64) (*RespPayment, error) {
	fmt.Println("3. Proses Payment ===========>>")

	url := os.Getenv("MIDTRANS_HOST") + "/charge"
	method := "POST"

	payment := models.Payment{}
	_payment, err := payment.GetPaymentByID(s.DB, reqCheckout.Payment)
	if err != nil {
		return nil, err
	}

	trxTransactionDetails := TrxTransactionDetails{
		GrossAmount: resCheckout.Detail.TotalPrice,
		OrderID:     fmt.Sprint(orderID) + `-{{timestamp}}`,
	}
	trxCustomerDetails := TrxCustomerDetails{
		Email:     reqCheckout.Email,
		FirstName: reqCheckout.Name,
		LastName:  reqCheckout.Name,
		Phone:     reqCheckout.Phone,
	}

	trxItemDetails := []TrxItemDetails{}

	trxBankTransfer := TrxBankTransfer{
		Bank:     _payment.Code,
		VANumber: _payment.RekNumber,
		FreeText: map[string]interface{}{},
	}
	trxCheckout := TrxCheckout{
		PaymentType:        _payment.PaymentType,
		TransactionDetails: trxTransactionDetails,
		CustomerDetails:    trxCustomerDetails,
		BankTransfer:       trxBankTransfer,
		ItemDetails:        trxItemDetails,
	}

	for _, row := range resCheckout.Product {
		if row.Checked {
			price := row.Price
			if row.PriceDisc > 0 {
				price = row.PriceDisc
			}
			trxCheckout.ItemDetails = append(trxCheckout.ItemDetails, TrxItemDetails{
				ID:       fmt.Sprint(row.ID),
				Name:     row.Name,
				Price:    price,
				Quantity: row.Qty,
			})
		}
	}
	trxCheckout.ItemDetails = append(trxCheckout.ItemDetails, TrxItemDetails{
		ID:       "shipping.fee",
		Name:     "Shipping Fee",
		Price:    resCheckout.Detail.ShippingFee,
		Quantity: 1,
	})

	body, err := json.Marshal(trxCheckout)
	if err != nil {
		return nil, err
	}

	// fmt.Println("string(body)", string(body))
	payload := strings.NewReader(string(body))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic U0ItTWlkLXNlcnZlci01N1pjUnU2MDNQTkJGX1FvdFlpbW9UVnA6U2MjIzA4NTczMDA2NjE2NQ==")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyResp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var respPayment RespPayment
	json.Unmarshal([]byte(bodyResp), &respPayment)

	if respPayment.StatusCode != "201" {
		fmt.Println("Error, Payment: ", respPayment.StatusMessage)
		return nil, errors.New(respPayment.StatusMessage)
	}

	return &respPayment, nil
}
