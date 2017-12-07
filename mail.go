package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"mime"
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
		sendMail(sub.Mail, id, bookinfo.LastNumber, bookinfo.BookName)
	}

}

func pushBookM(bookid string, des string) bool {

	bookinfo := BookInfo{}
	_, err := os.Stat("./bookinfo/" + bookid + ".json")
	if err != nil {
		return false
	}
	getJson("./bookinfo/" + bookid + ".json", &bookinfo)
	id := bookinfo.LastID

	_, err = os.Stat("./books/" + id + ".pdf")
	if err == nil {
		b := sendMail(des, id, bookinfo.LastNumber, bookinfo.BookName)
		return b
	}
	return false
}

func sendMail(des string, id string, lastid string, bookname string) bool {
	fmt.Println("Start Pushing Book!")
	config := Config{}
	_, err := os.Stat("./config/config.json")
	if err != nil {
		return false
	}

	getJson("./config/config.json", &config)

	m := gomail.NewMessage()
	m.SetHeader("From", config.MailUser)
	m.SetHeader("To", des)
	m.SetHeader("Subject", "[MangaPush]")
	m.SetBody("text/html", "Email from Manga Push")
	fileName := mime.QEncoding.Encode("utf-8", bookname+" 第"+lastid+"话.pdf")
	attachFileName := "./books/" + id + ".pdf"
	m.Attach(attachFileName, gomail.Rename(fileName))

	d := gomail.NewDialer(config.MailAddress, config.MailPort, config.MailUser, config.MailPassword)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Mail Sending Fail!")
		return false
	}
	fmt.Println("Mail Sending Succeed!")
	return true
}
