package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func pushBook(bookid int, id string) {

	sub := Subscription{}
	_, err := os.Stat("./config/sub.json")
	if err != nil {
		return
	}

	getJson("./config/sub.json", &sub)

	if contains(sub.List, bookid) == false {
		return
	}

	bookinfo := BookInfo{}
	_, err = os.Stat("./bookinfo/" + strconv.Itoa(bookid) + ".json")
	if err != nil {
		return
	}
	getJson("./bookinfo/"+strconv.Itoa(bookid)+".json", &bookinfo)

	_, err = os.Stat("./books/" + id + ".pdf")
	if err == nil {
		sendMail(sub.Mail, id, bookinfo.LastID)
	}

}

func pushBookM(bookid int, id string, des string) {

	bookinfo := BookInfo{}
	_, err := os.Stat("./bookinfo/" + strconv.Itoa(bookid) + ".json")
	if err != nil {
		return
	}
	getJson("./bookinfo/"+strconv.Itoa(bookid)+".json", &bookinfo)

	_, err = os.Stat("./books/" + id + ".pdf")
	if err == nil {
		sendMail(des, id, bookinfo.LastNumber)
	}

}

func sendMail(des string, id string, lastid string) {
	fmt.Println("Start Pushing Book!")
	config := Config{}
	_, err := os.Stat("./config/config.json")
	if err != nil {
		return
	}

	getJson("./config/config.json", &config)

	m := gomail.NewMessage()
	m.SetHeader("From", config.MailUser)
	m.SetHeader("To", des)
	m.SetHeader("Subject", "[MangaPush]")
	m.SetBody("text/html", "Email from Manga Push")
	m.Attach("./books/" + id + ".pdf", gomail.Rename("[" + lastid + "].pdf"))

	d := gomail.NewDialer(config.MailAddress, config.MailPort, config.MailUser, config.MailPassword)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		log.Panic(err)
	}
	fmt.Println("Mail Sending Succeed!")
}
