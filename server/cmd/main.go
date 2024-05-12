package main

import (
	"fmt"
	"lillith/autorender/server/database"
	"lillith/autorender/server/routes"
	"log"
	"net/http"
)

func main() {
	db := database.Connect()
	defer db.Close()

	fmt.Println("Connected to database")

	mux := http.NewServeMux()

	mux.HandleFunc("/storage/{path...}", routes.DownloadStorage)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	mux.HandleFunc("/videos/{id}", routes.VideoPage)

	mux.HandleFunc("/api/upload/demo", routes.UploadDemo)
	mux.HandleFunc("/api/upload/video", routes.UploadVideo)
	mux.HandleFunc("/api/fetch/queue", routes.FetchDemo)

	fmt.Println("Server active")
	err := http.ListenAndServe("0.0.0.0:80", mux)
	if err != nil {
		log.Fatalln(err)
	}
}
