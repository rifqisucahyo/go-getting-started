package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/responses"
)

func (server *Server) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.ExtractTokenID(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	var reqProduct []ReqProduct
	json.Unmarshal([]byte(body), &reqProduct)

	result, err := getDataCheckout(server, userID, reqProduct)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) AddCart(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	userID, err := auth.ExtractTokenID(r)
	var reqProduct []ReqProduct
	json.Unmarshal([]byte(body), &reqProduct)

	err = checkoutProcess(server, userID, reqProduct)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Proses checkout gagal."))
	}

	result, err := getDataCheckout(server, userID, reqProduct)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	responses.JSON(w, http.StatusOK, result)
}
