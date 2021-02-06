package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
)

func (server *Server) GetPath(w http.ResponseWriter, r *http.Request) {
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

	secret := make(map[string]interface{})
	secret = map[string]interface{}{
		"e": productDet,
	}

	responses.JSON(w, http.StatusOK, secret)

}
