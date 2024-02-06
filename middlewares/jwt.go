package middlewares

import (
	"fmt"
	"net/http"
	"rest_full_api/helper"
	"rest_full_api/config"

	"github.com/golang-jwt/jwt/v5"
)

type joinedError struct {
	errs []error
}


func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		
		c, err := r.Cookie("token")
		if err != nil {
			response := map[string]interface{}{"message": "Unauthorized"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return 
		}

		// mengambil token value 
		tokenString := c.Value

		claims := &config.JWTClaim{}
		fmt.Println(claims)
		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error){
			return config.JWT_KEY, nil
		})

		if err != nil {
			switch err {
			case jwt.ErrSignatureInvalid:
				// token invalid
				response := map[string]interface{}{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case jwt.ErrTokenExpired:
				// token expired
				response := map[string]interface{}{"message": "Unauthorized, Token Expired"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			default:
				response := map[string]interface{}{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
			}
		}

		if !token.Valid {
			response := map[string]interface{}{"message": "Unathorized"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}