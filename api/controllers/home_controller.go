package controllers

import (
	"net/http"

	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Welcome to alatkita.id API",
	})
}

func (server *Server) GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	infobox := models.InfoBox{}
	infoboxes, err := infobox.GetInfoBox(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	slider := models.Slider{}
	sliders, err := slider.GetSlider(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	product := models.Product{}
	newproducts, err := product.GetProduct(server.DB, 0, 12, "new")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	bestproducts, err := product.GetProduct(server.DB, 0, 12, "best")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	secret := make(map[string]interface{})
	secret = map[string]interface{}{
		"banner": map[string]interface{}{
			"img":   "https://dev.alatkita.id/static/media/home-banner.e2acca42.webp",
			"title": "Belanja Kebutuhan Alat di Alatkita.id",
			"desc":  "Solusi kebutuhan Anda dalam melengkapi kebutuhan rumah, dapur, dan alat-alat lainnya. Harga dijamin murah dan berkualitas.",
		},
		"infobox": infoboxes,
		"slider":  sliders,
		"products": map[string]interface{}{
			"new":  newproducts,
			"best": bestproducts,
		},
	}

	responses.JSON(w, http.StatusOK, secret)
}
