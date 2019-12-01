package model

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	u "go-api/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}
type Account struct {
	ID        uint   `schema:"id";gorm:"AUTO INCREMENT;PRIMARY KEY"`
	FirstName string `schema:"first_name"`
	LastName  string `schema:"last_name"`
	Email     string `schema:"email" `
	Password  string `schema:"password"`
	Token     string `schema:"token";sql:"-"`
}

//validating incoming user

func (account *Account) ValidateSignUp() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "invalid Password!"), false
	}

	//Email must be unique
	temp := &Account{}

	err := GetDB().Table("account").Find(&temp).Where("email=?", account.Email).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	//if err!= nil{
	//	return u.Message(false,"Email address already in use by another user."),false
	//}
	return u.Message(false, "Requirement passed"), true

}

func (account *Account) CreateAccount() map[string]interface{} {

	if res, ok := account.ValidateSignUp(); !ok {
		return res
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return u.Message(false, "Failed to create account, connection error.")
	}
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = "Bearer" + tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["data"] = account
	return response
}

func Login(email, password string) map[string]interface{} {
	account := &Account{}

	err := GetDB().Table("account").First(account).Where("email=?", email).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")

	}
	log.Println(account.Password)

	log.Println("\n")

	log.Println(password)

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	//Worked! Logged In
	account.Password = ""

	//Create JWT token

	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = "Bearer" + tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp

}

func GetUser(u uint) *Account {
	acc := &Account{}

	err := GetDB().Table("account").Where("id = ?", u).Find(acc).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}



	acc.Password = " "

	return acc
}
