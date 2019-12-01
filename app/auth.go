package app

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-api/model"
	"go-api/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/register", "/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                    //current request path
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})

		token := r.Header.Get("X-AUTH-TOKEN")

		if token == "" { //Token is missing, returns with error code 403 Unauthorized
			response = utils.Message(false, "Missing Auth token!!")

			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

			utils.Responds(w, response)
			return

		}
		log.Println(token)

		split := strings.Split(token, "Bearer") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		log.Println(len(split))

		if len(split) != 2 {
			response = utils.Message(false, "Invalid Auth token !!")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

			utils.Responds(w, response)
			return
		}

		interestedTokenPath := split[1]

		tk := &model.Token{}

		fmt.Println("see", interestedTokenPath)

		toke, err := jwt.ParseWithClaims(interestedTokenPath, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = utils.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			utils.Responds(w, response)
			return
		}

		if !toke.Valid { //Token is invalid, maybe not signed on this server
			response = utils.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			utils.Responds(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token

		fmt.Sprint("user%", tk.UserId)

		ctx := context.WithValue(r.Context(), "user", tk.UserId)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!

	})
}
