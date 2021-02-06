package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
)

type ReqCourier struct {
	Address  ReqAddress   `json:"address"`
	Products []ReqProduct `json:"products"`
}

func (s *Server) getCourier(CourierID int64) ([]map[string]interface{}, string, error) {
	fmt.Println("Proses Courier ===========>>")

	courier := models.Courier{}
	data, err := courier.GetCourier(s.DB, CourierID)
	if err != nil {
		return nil, "", err
	}

	_couriers := make([]map[string]interface{}, 0, 0)
	_courier_code := make([]string, 0, 0)
	for _, row := range *data {
		_row := map[string]interface{}{
			"id":          row.Id,
			"code":        row.Code,
			"name":        row.Name,
			"service":     row.Service,
			"description": row.Description,
			"etd":         "",
			"value":       0,
		}
		_couriers = append(_couriers, _row)
		if !InArrayV2(_courier_code, row.Code) {
			_courier_code = append(_courier_code, row.Code)
		}
	}
	courier_code := strings.Join(_courier_code, ":")

	return _couriers, courier_code, nil
}

func (server *Server) GetCourier(w http.ResponseWriter, r *http.Request) {
	var err error
	userID, _ := auth.ExtractTokenID(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	var reqCourier ReqCourier
	json.Unmarshal([]byte(body), &reqCourier)

	result, err := getDataCheckout(server, userID, reqCourier.Products)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	_couriers, courier_code, err := server.getCourier(0)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	couriers, _, err := shipmentCheck(fmt.Sprint(result.Detail.Weight), fmt.Sprint(reqCourier.Address.SubdistrictID), courier_code, _couriers, "list")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, couriers)
}

func shipmentCheck(weight string, subdistrictID string, courier string, couriers []map[string]interface{}, typeCheck string) ([]map[string]interface{}, int64, error) {
	type RoStatus struct {
		Code        int64  `json:"code"`
		Description string `json:"description"`
	}
	type RoQuery struct {
		Origin          string `json:"origin"`
		OriginType      string `json:"origintype"`
		Destination     string `json:"destination"`
		DestinationType string `json:"destinationtype"`
		Weight          int64  `json:"weight"`
		Courier         string `json:"courier"`
		Length          int64  `json:"length"`
		Width           int64  `json:"width"`
		Height          int64  `json:"height"`
		Diameter        int64  `json:"diameter"`
	}
	type RoOriginDetails struct {
		Type             string `json:"destinationtype"`
		Province_id      string `json:"province_id"`
		Province         string `json:"province"`
		City_id          string `json:"city_id"`
		City             string `json:"city"`
		Subdistrict_id   string `json:"subdistrict_id"`
		Subdistrict_name string `json:"subdistrict_name"`
	}
	type RoDestinationDetails struct {
		Type             string `json:"destinationtype"`
		Province_id      string `json:"province_id"`
		Province         string `json:"province"`
		City_id          string `json:"city_id"`
		City             string `json:"city"`
		Subdistrict_id   string `json:"subdistrict_id"`
		Subdistrict_name string `json:"subdistrict_name"`
	}
	type RoCost struct {
		Value int64  `json:"value"`
		Etd   string `json:"etd"`
		Note  string `json:"note"`
	}
	type RoCosts struct {
		Service     string `json:"service"`
		Description string `json:"description"`
		Cost        []RoCost
	}
	type RoResult struct {
		Code  string `json:"code"`
		Name  string `json:"name"`
		Costs []RoCosts
	}
	type DataRajaOngkir struct {
		Status             RoStatus
		Query              RoQuery
		OriginDetails      RoOriginDetails
		DestinationDetails RoDestinationDetails
		Results            []RoResult
	}
	type RajaOngkir struct {
		Rajaongkir DataRajaOngkir `json:"rajaongkir"`
	}

	url := os.Getenv("RAJAONGKIR_HOST") + "/cost"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("origin", "309")
	_ = writer.WriteField("originType", "subdistrict")
	_ = writer.WriteField("destination", subdistrictID)
	_ = writer.WriteField("destinationType", "subdistrict")
	_ = writer.WriteField("weight", weight)
	_ = writer.WriteField("courier", courier)
	_ = writer.WriteField("length", "")
	_ = writer.WriteField("width", "")
	_ = writer.WriteField("height", "")
	_ = writer.WriteField("diameter", "")
	err := writer.Close()
	if err != nil {
		return couriers, 0, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return couriers, 0, err
	}
	req.Header.Add("key", os.Getenv("RAJAONGKIR_KEY"))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	r, err := client.Do(req)
	if err != nil {
		return couriers, 0, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return couriers, 0, err
	}
	var rajaOngkir RajaOngkir
	json.Unmarshal([]byte(body), &rajaOngkir)

	if typeCheck == "checkout" {
		var total int64
		if rajaOngkir.Rajaongkir.Status.Code == 200 {
			for _, v := range couriers {
				if len(rajaOngkir.Rajaongkir.Results) > 0 {
					for _, row := range rajaOngkir.Rajaongkir.Results {
						if row.Code == v["code"] {
							if len(row.Costs) > 0 {
								for _, rowService := range row.Costs {
									if rowService.Service == v["service"] {
										if len(rowService.Cost) > 0 {
											for _, rowCost := range rowService.Cost {
												if rowCost.Value > 0 {
													total = rowCost.Value
													break
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return nil, total, nil
		} else {
			return nil, 0, errors.New(rajaOngkir.Rajaongkir.Status.Description)
		}

	} else {

		var newCourier = make([]map[string]interface{}, 0, 0)
		var flag bool
		for i, v := range couriers {
			if rajaOngkir.Rajaongkir.Status.Code == 200 {
				if len(rajaOngkir.Rajaongkir.Results) > 0 {
					for _, row := range rajaOngkir.Rajaongkir.Results {
						if row.Code == v["code"] {
							if len(row.Costs) > 0 {
								for _, rowService := range row.Costs {
									if rowService.Service == v["service"] {
										if len(rowService.Cost) > 0 {
											for _, rowCost := range rowService.Cost {
												if rowCost.Value > 0 {
													if rowCost.Value > 0 {
														flag = true
														couriers[i]["etd"] = rowCost.Etd + " Hari"
														couriers[i]["value"] = rowCost.Value
														newCourier = append(newCourier, couriers[i])
													}
													break
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		if flag == false {
			return newCourier, 0, errors.New("Alamat Anda belum terdata pada pengiriman kami.")
		} else {
			return newCourier, 0, err
		}
	}
}
