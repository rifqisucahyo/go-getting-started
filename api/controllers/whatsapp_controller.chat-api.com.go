package controllers

// import (
// 	"encoding/json"
// 	"errors"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/victorsteven/fullstack/api/responses"
// )

// type RespWhatsapp struct {
// 	Sent        bool   `json:"sent"`
// 	Message     string `json:"message"`
// 	ID          string `json:"id"`
// 	QueueNumber int64  `json:"queueNumber"`
// }
// type ReqWhatsapp struct {
// 	Phone   string
// 	Message string
// }

// func (s *Server) SendWhatsapp(w http.ResponseWriter, r *http.Request) {
// 	data, err := s.sendWhatsapp("085730093080", "Halo apa kabar cahyo, alatkita.id")
// 	if err != nil {
// 		responses.ERROR(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	responses.JSON(w, http.StatusOK, data)
// }

// func (s *Server) sendWhatsapp(phone string, message string) (*RespWhatsapp, error) {

// 	var err error
// 	var respWhatsapp RespWhatsapp

// 	url := os.Getenv("WA_HOST") + "/sendMessage?token=" + os.Getenv("WA_TOKEN")
// 	method := "POST"

// 	payload := strings.NewReader(`{"phone": "` + phone + `","body": "` + message + `"}`)

// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, payload)

// 	if err != nil {
// 		return &respWhatsapp, err
// 	}
// 	req.Header.Add("Content-Type", "application/json")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return &respWhatsapp, err
// 	}
// 	defer res.Body.Close()

// 	bodyResp, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return &respWhatsapp, err
// 	}
// 	json.Unmarshal([]byte(bodyResp), &respWhatsapp)

// 	if respWhatsapp.Sent == false {
// 		return &respWhatsapp, errors.New(respWhatsapp.Message)
// 	}

// 	return &respWhatsapp, err
// }
