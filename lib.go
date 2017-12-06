package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Book struct {
	NewData Data `json:"data"`
	ErrNo   int  `json:"errNo"`
}

type BookInfo struct {
	BookName   string `json:"book_text"`
	Author     string `json:"author"`
	LastNumber string `json:"number"`
	LastID     string `json:"id"`
	UpdateTime string `json:"time"`
}

type Data struct {
	BookText   string `json:"book_text"`
	Title      string `json:"title"`
	Number     int    `json:"number"`
	ContentImg string `json:"content_img"`
	Book       int    `json:"book"`
}

type Config struct {
	LastUpdate   int    `json:"last_id"`
	MailAddress  string `json:"mail_address"`
	MailPort     int    `json:"mail_port"`
	MailUser     string `json:"mail_user"`
	MailPassword string `json:"mail_password"`
}

type Subscription struct {
	Mail string `json:"mail"`
	List []int  `json:"list"`
}

func getJson(path string, target interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	return json.Unmarshal(file, target)
}
