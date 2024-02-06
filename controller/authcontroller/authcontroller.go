package authcontroller

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
	"rest_full_api/models"
	"rest_full_api/helper"
	"rest_full_api/config"
)


func Login(w http.ResponseWriter, r *http.Request) {

	var userLogin models.User

	// mengabil data json login
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userLogin); err != nil {
		log.Fatal("Gagal mendecode json")
		message := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, message)
		return
	} 
	defer r.Body.Close()

	// ambil data berdasarkan username
	var user models.User
	if err := models.DB.Where("user_name = ?", userLogin.UserName).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			message := map[string]string{"message": "username atau password salah"}
			helper.ResponseJSON(w, http.StatusUnauthorized, message)
			return
		default:
			message := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, message)
			return
		}
	}

	// memeriksa apakah password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			message := map[string]string{"message": "Password Yang Dimasukkan Salah"}
			helper.ResponseJSON(w, http.StatusUnauthorized, message)
			return
		default:
			message := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, message)
			return
		}
	}

	// proses pembuatan token jwt
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		UserName: userLogin.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:	"go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// mendeklarasi algoritma yang akan digunakan untuk signing
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		message := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, message)
		return
	}

	// set token yang ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:	"token",
		Path:	"/",
		Value:	token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "login berhasil"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// inialisasi validator  
	validate := validator.New(validator.WithRequiredStructEnabled())
	
	// mengambil inputan json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		log.Fatal("Gagal mendecode json")
		message := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, message)
	}
	defer r.Body.Close()

	// validasi data userInput (aunthentication)
	trans, validate := helper.TranslateToIndonesia()
	err := validate.Struct(userInput)
	if err != nil {
		// memanggil funstion translator untuk translate message error
		errs := err.(validator.ValidationErrors)
		messageErrors := errs.Translate(trans)
		response := map[string]interface{}{"message": messageErrors}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}  else {
		var existingUser models.User
		err := models.DB.Where("user_name = ?", userInput.UserName).First(&existingUser).Error 
		switch err {
		case gorm.ErrRecordNotFound:
				// hash password menggunkan bcrypt
				hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
				userInput.Password = string(hashPassword)

				//insert database
				if err := models.DB.Create(&userInput).Error; err != nil {
					log.Fatal("Gagal menyimpan data")
				}

				// response yang di berikan jika berhasil
				response := map[string]string{"message":"success"}
				helper.ResponseJSON(w, http.StatusOK, response)
				return	
		case nil:
			// if username alredy exist in database
			message := map[string]string{"message": "username sudah digunakan"}
			helper.ResponseJSON(w, http.StatusConflict, message)
		default:
			message := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, message)
		}
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	// hapus token yang ada di cookie
	http.SetCookie(w, &http.Cookie{
		Name:	"token",
		Path:	"/",
		Value:	"", 
		HttpOnly: true, 
		MaxAge:		-1,
	})

	response := map[string]string{"message": "logout berhasil"}
	helper.ResponseJSON(w, http.StatusOK, response)
}