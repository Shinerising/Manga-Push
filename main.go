package main

import (
	"io"
	"os"
    "log"
	"fmt"
	"time"
	"path"
	"image"
	"errors"
	"strconv"
    "net/http"
    "io/ioutil"
	"image/jpeg"
	"path/filepath"
	"html/template"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
)

type Book struct {
	NewData Data `json:"data"`
	ErrNo int `json:"errNo"`
}

type BookInfo struct {
	BookName string `json:"book_text"`
	Author string `json:"author"`
	LastNumber string `json:"number"`
	LastID string `json:"id"`
	UpdateTime string `json:"time"`
}

type Data struct {
    BookText string `json:"book_text"`
    Title string `json:"title"`
    Number int `json:"number"`
    ContentImg string `json:"content_img"`
    Book int `json:"book"`
}

type Config struct {
	LastUpdate int `json:"last_id"`
}

func getJson(path string, target interface{}) error {
    file, err := ioutil.ReadFile(path)
    if err != nil {
        log.Panic(err)
    }

    return json.Unmarshal(file, target)
}

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

func fileHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := os.Stat("./books/" + id + ".pdf")
    if os.IsNotExist(err) {
    	w.WriteHeader(http.StatusNotFound)
    	fmt.Fprint(w, "Not Found!")
    } else {
		http.ServeFile(w, r, "./books/" + id + ".pdf")
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
		getJson("./bookinfo/" + bookid + ".json", &bookinfo)
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

func bookManagement(id string) bool {
	_, err := os.Stat("./bookinfo/")
    if os.IsNotExist(err) {
    	_ = os.Mkdir("./bookinfo/", os.ModePerm)
    }

	book, err := downloadBook(id)
    if err != nil {
        log.Panic(err)
        if err.Error() == "Not Found" {
    		return false
    	} else {
    		return true
    	}
    }
    bookid := book.NewData.Book
    bookinfo := BookInfo{}
    _, err = os.Stat("./bookinfo/" + strconv.Itoa(bookid) + ".json")
    if os.IsNotExist(err) {
    	bookinfo.BookName = book.NewData.BookText
    	bookinfo.Author = "未知作者"
    } else {
		getJson("./bookinfo/" + strconv.Itoa(bookid) + ".json", &bookinfo)
    }
	bookinfo.LastNumber = strconv.Itoa(book.NewData.Number)
	bookinfo.LastID = id
	now := time.Now()
	bookinfo.UpdateTime = now.Format("2006-01-02")

	jsonString, _ := json.Marshal(bookinfo)
    err = ioutil.WriteFile("./bookinfo/" + strconv.Itoa(bookid) + ".json", jsonString, 0644)
    if err != nil {
        log.Panic(err)
    }
    return true
}

func startTask() {
	fmt.Println("Start Task!")
	_, err := os.Stat("./config/")
    if os.IsNotExist(err) {
    	_ = os.Mkdir("./config/", os.ModePerm)
    }

    ticker := time.NewTicker(4 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
				case <- ticker.C:
					taskHandler()
				case <- quit:
					ticker.Stop()
				return
			}
		}
	}()
}

func taskHandler () {
    config := Config{}
    _, err := os.Stat("./config/config.json")
    if os.IsNotExist(err) {
    	config.LastUpdate = 10210
    } else {
		getJson("./config/config.json", &config)
    }

    id := config.LastUpdate
	sum := 0
	for sum < 10 {
		fmt.Println("Checking " + strconv.Itoa(id + sum))
		if bookManagement(strconv.Itoa(id + sum)) {
			sum += 1
		} else {
			break
		}
	}
	config.LastUpdate = id + sum;
	jsonString, _ := json.Marshal(config)
    err = ioutil.WriteFile("./config/config.json", jsonString, 0644)
    if err != nil {
        log.Panic(err)
        return
    }
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
    	downloadJson(url, "./books/" + id + ".json")
    }
	getJson("./books/" + id + ".json", &book)

	if (book.ErrNo != 0) {
		os.Remove("./books/" + id + ".json");
		return book, errors.New("Not Found")
	}

	if (book.NewData.ContentImg == "null" || book.NewData.ContentImg == "{}") {
		os.Remove("./books/" + id + ".json");
		return book, errors.New("No Images")
	}

    _, err = os.Stat("./books/" + id)
    if os.IsNotExist(err) {
    	_ = os.Mkdir("./books/" + id, os.ModePerm)
    }

    _, err = os.Stat("./books/" + id + "/.downloaded")
    if os.IsNotExist(err) {

	    var objmap map[string]interface{}

	    err := json.Unmarshal([]byte(book.NewData.ContentImg), &objmap)
	    if err != nil {
			log.Panic(err)
			return book, errors.New("Invalid JSON")
		}
	    for k := range objmap {
	    	downloadImage(id, k, objmap[k].(string))
	    }

	    _ = os.Mkdir("./books/" + id + "/.downloaded", os.ModePerm)

    }

    _, err = os.Stat("./books/" + id + ".pdf")
    if os.IsNotExist(err) {

    	files, err := ioutil.ReadDir("./books/" + id)
	    if err != nil {
	        log.Panic(err)
			return book, errors.New("Folder Error")
	    }

	    pdf := gofpdf.New("P", "mm", "A5", "")
	    pdf.SetTitle("第" + strconv.Itoa(book.NewData.Number) + "话 " + book.NewData.Title, true)
	    pdf.SetSubject("第" + strconv.Itoa(book.NewData.Number) + "话 " + book.NewData.Title, true)
	    pdf.SetAuthor(book.NewData.BookText, true)
	    pageSize := gofpdf.SizeType{}
	    var imgOptions gofpdf.ImageOptions
	    imgOptions.ReadDpi = true

	    for _, f := range files {
	    	if(f.Name() != ".downloaded") {
		    	filename := "./books/" + id + "/" + f.Name()
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

    if filepath.Ext(name) != ".jpg" && filepath.Ext(name) != ".jpeg"  {
		new_name := name[0:len(name)-len(filepath.Ext(name))] + ".jpg"

	    file, err := os.Create(new_name)
	    if err != nil {
	        log.Panic(err)
	    }
	    img, _, err := image.Decode(response.Body)
		if err != nil {
	        log.Panic(err)
		}
		jpeg.Encode(file, img, &jpeg.Options{ 100 })
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

func main() {
	startTask()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    http.HandleFunc("/", pageHandler)
    http.HandleFunc("/file", fileHandler)
    http.HandleFunc("/fetch", fetchHandler)
    http.HandleFunc("/detail", detailHandler)
    http.HandleFunc("/download", downloadHandler)
    port := os.Getenv("PORT")
    if port == "" {
    	port = "8080"
    }
    http.ListenAndServe(":" + port, nil)
	fmt.Println("Start Listening to " + port)
}