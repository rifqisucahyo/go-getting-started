package controllers

import (
	"net/http"
	"strconv"

	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
)

func (server *Server) GetRegion(w http.ResponseWriter, r *http.Request) {
	start, err := strconv.ParseInt(r.FormValue("start")[0:], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	limit, err := strconv.ParseInt(r.FormValue("limit")[0:], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if start == 0 {
		start = 0
	}
	if limit == 0 {
		limit = 10
	}
	q := r.FormValue("q")

	data, err := models.GetRegion(server.DB, q, start, limit)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	result := map[string]interface{}{
		"query": map[string]interface{}{
			"q":     q,
			"start": start,
			"limit": limit,
		},
		"data":   data,
		"status": true,
	}
	responses.JSON(w, http.StatusOK, result)
}
