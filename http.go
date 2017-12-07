package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func pageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./views/index.html")
	if err != nil {
		log.Panic(err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Panic(err)
		return
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	bookManagement(id)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	mail := r.URL.Query().Get("mail")
	if pushBookM(id, mail) {
		fmt.Fprintf(w, "succeed")
	} else {
		fmt.Fprintf(w, "fail")
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := os.Stat("./books/" + id + ".pdf")
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found!")
	} else {
		http.ServeFile(w, r, "./books/"+id+".pdf")
	}
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	bookid := r.URL.Query().Get("id")
	_, err := os.Stat("./bookinfo/")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./bookinfo/", os.ModePerm)
	}

	bookinfo := BookInfo{}
	_, err = os.Stat("./bookinfo/" + bookid + ".json")
	if os.IsNotExist(err) {
		bookinfo.BookName = ""
	} else {
		getJson("./bookinfo/"+bookid+".json", &bookinfo)
	}
	json.NewEncoder(w).Encode(bookinfo)
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	name := "./tmp/" + path.Base(url)

	_, err := os.Stat("./tmp/")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./tmp/", os.ModePerm)
	}

	_, err = os.Stat(name)

	if os.IsNotExist(err) {
		response, err := http.Get(url)
		if err != nil {
			log.Panic(err)
		}

		defer response.Body.Close()

		file, err := os.Create(name)
		if err != nil {
			log.Panic(err)
		}
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Panic(err)
		}
		file.Close()
	}

	http.ServeFile(w, r, name)
}
