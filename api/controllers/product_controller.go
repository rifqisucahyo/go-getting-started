package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
)

func (server *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Start, Limit int64
		Type, Q      string
	}
	_start, _ := strconv.ParseInt(r.FormValue("start")[0:], 10, 64)
	_limit, _ := strconv.ParseInt(r.FormValue("limit")[0:], 10, 64)
	if _limit == 0 {
		_limit = 8
	}

	par := param{
		Start: _start,
		Limit: _limit,
		Type:  r.FormValue("type"),
		Q:     r.FormValue("q"),
	}

	vars := mux.Vars(r)
	slug := vars["slug"]
	_slug := strings.Split(slug, "-")
	productID, err := strconv.ParseInt(_slug[len(_slug)-1], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, formaterror.FormatError(err.Error()))
		return
	}

	product := models.Product{}
	productDet, err := product.ProductByID(server.DB, int32(productID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, formaterror.FormatError(err.Error()))
		return
	}

	fmt.Println("par", par)
	bestproducts, _ := product.GetProduct(server.DB, par.Start, par.Limit, par.Type)
	productGallery, _ := models.GalleryByID(server.DB, int32(productID))

	secret := make(map[string]interface{})
	secret = map[string]interface{}{
		"data":          productDet,
		"recomendation": bestproducts,
		"breadcrumb": []map[string]interface{}{
			{
				"name": "Home",
				"link": "/",
			},
			{
				"name": productDet.Brand,
				"link": fmt.Sprintf("/search?q=%s", productDet.Brand),
			},
			{
				"name": productDet.Name,
				"link": "/",
			},
		},
		"gallery": productGallery,
	}

	responses.JSON(w, http.StatusOK, secret)
}

func (server *Server) GetProducts(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	par := models.ReqProduct{}
	json.Unmarshal([]byte(body), &par)

	if par.Limit == 0 {
		par.Limit = 12
	}

	product := models.Product{}
	products, err := product.GetProducts(server.DB, par)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, formaterror.FormatError(err.Error()))
		return
	}

	secret := make(map[string]interface{})
	secret = map[string]interface{}{
		"data": products,
		"breadcrumb": []map[string]interface{}{
			{
				"name": "Home",
				"link": "/",
			},
			{
				"name": "Jakemy",
				"link": "/",
			},
			{
				"name": "Jakemy 79 in 1 Ratchet Screwdriver Set - JM-6107",
				"link": "/",
			},
		},
	}

	responses.JSON(w, http.StatusOK, secret)
}
