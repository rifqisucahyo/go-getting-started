package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	dataUser, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, dataUser)
}

func (server *Server) SignIn(email, password string) (map[string]interface{}, error) {

	var err error
	user := models.User{}

	param := make(map[string]interface{})

	err = server.DB.Debug().Model(models.User{}).Select("user_id as id, password, username, name, phone, email, created_at, updated_at").Where("email = ?", email).Take(&user).Error
	if err != nil {
		return param, err
	}

	// dataUser, err := user.FindUserByEmail(server.DB, email)
	// if err != nil {
	// 	return param, err
	// }
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return param, err
	}

	token, _ := auth.CreateToken(user.ID)

	param["token"] = token
	param["data"] = user

	return param, nil
}

func (server *Server) CheckAuth(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, true)
}

type ReqOtp struct {
	Code  int64  `json:"code"`
	Phone string `json:"phone"`
}

func (s *Server) PhoneOTPRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	var reqOtp ReqOtp
	json.Unmarshal([]byte(body), &reqOtp)

	if reqOtp.Phone == "" {
		responses.ERROR(w, http.StatusBadRequest, errors.New("No. Phone tidak boleh kosong."))
		return
	}

	dataVerify := models.Verification{}
	total, err := models.CheckPhoneOTP(s.DB, reqOtp.Phone)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	code, err := strconv.ParseInt(generateCode(4), 10, 64)
	message := `JANGAN BERI kode ini ke siapa pun, TERMASUK ALATKITA. WASPADA PENIPUAN!, masukkan kode verifikasi (OTP) ` + fmt.Sprint(code)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if total > 0 {
		dataVerify, err = models.PhoneOTP(s.DB, reqOtp.Phone)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		if dataVerify.Timer <= 0 {
			fmt.Println("UPDATE OTP ===>")
			//update
			err = models.UpdatePhoneOTP(s.DB, code, reqOtp.Phone)
			if err != nil {
				responses.ERROR(w, http.StatusInternalServerError, err)
				return
			}
			_, err := s.sendWhatsapp(reqOtp.Phone, message)
			if err != nil {
				responses.ERROR(w, http.StatusBadGateway, err)
				return
			}
		}
	} else {
		fmt.Println("CREATE OTP ===>")
		err = models.CreatePhoneOTP(s.DB, code, reqOtp.Phone)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		_, err := s.sendWhatsapp(reqOtp.Phone, message)
		if err != nil {
			responses.ERROR(w, http.StatusBadGateway, err)
			return
		}
	}

	dataVerify, err = models.PhoneOTP(s.DB, reqOtp.Phone)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	result := map[string]interface{}{
		"timer": dataVerify.Timer,
		"phone": reqOtp.Phone,
	}
	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) LoginPhone(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var reqOtp ReqOtp
	json.Unmarshal([]byte(body), &reqOtp)

	if reqOtp.Phone == "" {
		responses.ERROR(w, http.StatusBadRequest, errors.New("No. Phone tidak boleh kosong."))
		return
	}
	if reqOtp.Code == 0 {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Kode Verifikasi tidak boleh kosong."))
		return
	}
	dataUser, err := server.SignInPhone(reqOtp.Phone, reqOtp.Code)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, dataUser)
}

func (server *Server) SignInPhone(phone string, code int64) (map[string]interface{}, error) {

	var err error
	param := make(map[string]interface{})

	dataVerify, err := models.PhoneOTPVerification(server.DB, code, phone)
	if err != nil {
		return param, errors.New("Data tidak ditemukan / tidak sesuai.")
	}

	if dataVerify.Timer == 0 {
		return param, errors.New("Kode OTP tidak ditemukan, lakukan permintaan ulang.")
	}
	if dataVerify.Timer < 0 {
		return param, errors.New("Kode OTP sudah expired, lakukan permintaan ulang.")
	}

	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Select("user_id as id, password, username, name, phone, email, created_at, updated_at").Where("phone = ?", phone).Take(&user).Error
	if err != nil {
		if fmt.Sprint(err) == "record not found" {
			user = models.User{}
			user.ID = 0
			user.Username = ""
			user.Name = phone
			user.Email = ""
			user.Phone = phone
			user.CreatedAt = time.Now()
			user.UpdatedAt = time.Now()

			_, err = user.SaveUser(server.DB)
			if err != nil {
				return param, err
			}
			user = models.User{}
			err = server.DB.Debug().Model(models.User{}).Select("user_id as id, password, username, name, phone, email, created_at, updated_at").Where("phone = ?", phone).Take(&user).Error
			if err != nil {
				return param, err
			}
		} else {
			return param, err
		}
	}

	token, _ := auth.CreateToken(user.ID)

	param["token"] = token
	param["data"] = user

	return param, nil
}
