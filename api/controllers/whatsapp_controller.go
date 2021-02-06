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

	"github.com/victorsteven/fullstack/api/responses"
)

// type RespWhatsapp struct {
// 	Sent        bool   `json:"sent"`
// 	Message     string `json:"message"`
// 	ID          string `json:"id"`
// 	QueueNumber int64  `json:"queueNumber"`
// }
type RespWhatsapp struct {
	Status  bool     `json:"status"`
	Target  []string `json:"target"`
	Message string   `json:"message"`
	ID      []int64  `json:"id"`
	Process string   `json:"process"`
}
type ReqWhatsapp struct {
	Phone   string
	Message string
}

func (s *Server) SendWhatsapp(w http.ResponseWriter, r *http.Request) {
	data, err := s.sendWhatsapp("085730093080", "Halo apa kabar cahyo, alatkita.id")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, data)
}

func (s *Server) sendWhatsapp(phone string, message string) (*RespWhatsapp, error) {

	var err error
	var respWhatsapp RespWhatsapp

	url := os.Getenv("WA_HOST") + "/send_message.php"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("phone", phone)
	_ = writer.WriteField("type", "text")
	_ = writer.WriteField("delay", "2")
	_ = writer.WriteField("delay_req", "2")
	_ = writer.WriteField("text", message)
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return &respWhatsapp, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &respWhatsapp, err
	}
	req.Header.Add("Authorization", os.Getenv("WA_TOKEN"))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &respWhatsapp, err
	}
	defer res.Body.Close()

	bodyResp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &respWhatsapp, err

	}
	json.Unmarshal([]byte(bodyResp), &respWhatsapp)

	if respWhatsapp.Status == false {
		return &respWhatsapp, errors.New(respWhatsapp.Message)
	}

	fmt.Println(string(bodyResp))
	return &respWhatsapp, err

}
