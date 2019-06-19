package main

import (
	utils "CSE350-1/chaincode/cd1/utils-golang"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type User struct {
	Doctype        string   `json:"DocType"`
	Name           string   `json:"Name"`
	Email          string   `json:"Email"`
	PasswordHash   string   `json:"PasswordHash"`
	Token          string   `json:"Token"`
	Key            string   `json:"key"`
	ownDocs        []string `json:"ownDocs"`
	accessableDocs []string `json:"AcsessableDocList"`
}

func main() {
	name := "nammi"
	email := "dsadasd"
	password := "dsadsad"

	h := sha256.New()
	h.Write([]byte(password))
	passwordHash := fmt.Sprintf("%x", h.Sum(nil))

	token := utils.RandomString()

	key := name + utils.RandomString()
	n := "nullDoc"
	var ownDocs = []string{n}
	var acsDocs = []string{n}
	user := User{"user", name, email, passwordHash, token, key, ownDocs, acsDocs}
	fmt.Println(user)
	jsonUser, err := json.Marshal(user)

	fmt.Println(jsonUser)
	fmt.Println(err)
}
