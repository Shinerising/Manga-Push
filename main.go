package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	startTask()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/file", fileHandler)
	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/fetch", fetchHandler)
	http.HandleFunc("/detail", detailHandler)
	http.HandleFunc("/download", downloadHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Start Listening to " + port)
	http.ListenAndServe(":"+port, nil)
}
