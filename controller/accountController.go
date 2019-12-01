package controller

import (
	"encoding/json"
	"go-api/model"
	u "go-api/utils"
	"io/ioutil"
	"log"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	//log.Println(reqBody)

	var booking model.Account
	json.Unmarshal(reqBody, &booking)
	//  p = new(model.Account)

	log.Println(booking)

	resp := booking.CreateAccount()
	u.Responds(w, resp)
}

var UserLogin = func(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var booking model.Account
	json.Unmarshal(reqBody, &booking)

	resp := model.Login(booking.Email, booking.Password)
	u.Responds(w, resp)

}

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	deviceGUID := keys.Get("user_id") //Get return empty string if key not found
	data := model.GetUser(deviceGUID)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Responds(w, resp)
}
