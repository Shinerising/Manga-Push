package main

import (
	"encoding/json"
	"errors"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/h2non/filetype.v1"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func bookManagement(id string) (bool, int) {
	_, err := os.Stat("./bookinfo/")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./bookinfo/", os.ModePerm)
	}

	book, err := downloadBook(id)
	if err != nil {
		if err.Error() == "Not Found" {
			return false, 0
		} else {
			return true, book.NewData.Book
		}
	}
	bookid := book.NewData.Book
	bookinfo := BookInfo{}
	_, err = os.Stat("./bookinfo/" + strconv.Itoa(bookid) + ".json")
	if os.IsNotExist(err) {
		bookinfo.BookName = book.NewData.BookText
		bookinfo.Author = "未知作者"
	} else {
		getJson("./bookinfo/"+strconv.Itoa(bookid)+".json", &bookinfo)
	}
	bookinfo.LastNumber = strconv.Itoa(book.NewData.Number)
	bookinfo.LastID = id
	now := time.Now()
	bookinfo.UpdateTime = now.Format("2006-01-02")

	jsonString, _ := json.Marshal(bookinfo)
	err = ioutil.WriteFile("./bookinfo/"+strconv.Itoa(bookid)+".json", jsonString, 0644)
	if err != nil {
		log.Panic(err)
	}
	return true, book.NewData.Book
}

func downloadBook(id string) (Book, error) {
	_, err := os.Stat("./books/")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./books/", os.ModePerm)
	}

	book := Book{}
	_, err = os.Stat("./books/" + id + ".json")
	if os.IsNotExist(err) {
		url := "http://hhzapi.ishuhui.com/cartoon/post/ver/76906890/id/" + id + ".json"
		downloadJson(url, "./books/"+id+".json")
	}
	getJson("./books/"+id+".json", &book)

	if book.ErrNo != 0 {
		os.Remove("./books/" + id + ".json")
		return book, errors.New("Not Found")
	}

	if book.NewData.ContentImg == "null" || book.NewData.ContentImg == "{}" {
		os.Remove("./books/" + id + ".json")
		return book, errors.New("No Images")
	}

	_, err = os.Stat("./books/" + id)
	if os.IsNotExist(err) {
		_ = os.Mkdir("./books/"+id, os.ModePerm)
	}

	_, err = os.Stat("./books/" + id + "/.downloaded")
	if os.IsNotExist(err) {

		var objmap map[string]interface{}

		err := json.Unmarshal([]byte(book.NewData.ContentImg), &objmap)
		if err != nil {
			log.Println(err)
			return book, errors.New("Invalid JSON")
		}
		for k := range objmap {
			downloadImage(id, k, objmap[k].(string))
		}

		_ = os.Mkdir("./books/"+id+"/.downloaded", os.ModePerm)

	}

	_, err = os.Stat("./books/" + id + ".pdf")
	if os.IsNotExist(err) {

		files, err := ioutil.ReadDir("./books/" + id)
		if err != nil {
			log.Panic(err)
			return book, errors.New("Folder Error")
		}

		bookid := book.NewData.Book
		bookinfo := BookInfo{}
		_, err = os.Stat("./bookinfo/" + strconv.Itoa(bookid) + ".json")
		if os.IsNotExist(err) {
			bookinfo.BookName = book.NewData.BookText
			bookinfo.Author = "未知作者"
		} else {
			getJson("./bookinfo/"+strconv.Itoa(bookid)+".json", &bookinfo)
		}

		pdf := gofpdf.New("P", "mm", "A5", "")
		pdf.SetTitle("第"+strconv.Itoa(book.NewData.Number)+"话 "+book.NewData.Title, true)
		pdf.SetSubject("第"+strconv.Itoa(book.NewData.Number)+"话 "+book.NewData.Title, true)
		pdf.SetAuthor(bookinfo.Author, true)
		pageSize := gofpdf.SizeType{}
		var imgOptions gofpdf.ImageOptions
		imgOptions.ReadDpi = true

		for _, f := range files {
			if f.Name() != ".downloaded" {
				filename := "./books/" + id + "/" + f.Name()
				buf, _ := ioutil.ReadFile(filename)
				kind, unknown := filetype.Match(buf)
				if unknown != nil {
					break
				}
				imgOptions.ImageType = kind.Extension
				infoPtr := pdf.RegisterImageOptions(filename, imgOptions)
				imgWd := infoPtr.Width()
				imgHt := infoPtr.Height()
				pageSize.Wd = imgWd
				pageSize.Ht = imgHt
				pdf.AddPageFormat("P", pageSize)
				pdf.ImageOptions(filename, 0, 0, imgWd, imgHt, false, imgOptions, 0, "")
			}
		}

		err = pdf.OutputFileAndClose("./books/" + id + ".pdf")
		if err != nil {
			log.Panic(err)
			return book, errors.New("PDF Error")
		}
	}
	return book, nil

}

func downloadImage(id string, name string, url string) {
	name = "./books/" + id + "/" + name
	url = "http://pic01.ishuhui.com" + url[7:len(url)]

	response, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	defer response.Body.Close()

	if filepath.Ext(name) != ".jpg" && filepath.Ext(name) != ".jpeg" && filepath.Ext(name) != ".png" {
		new_name := name[0:len(name)-len(filepath.Ext(name))] + ".jpg"

		file, err := os.Create(new_name)
		if err != nil {
			log.Panic(err)
		}
		img, _, err := image.Decode(response.Body)
		if err != nil {
			log.Panic(err)
		}
		jpeg.Encode(file, img, &jpeg.Options{100})
		file.Close()
	} else {
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
}

func downloadJson(url string, path string) {
	response, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	defer response.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Panic(err)
	}
	file.Close()
}
