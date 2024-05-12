package routes

import (
	"context"
	"fmt"
	"io"
	"lillith/autorender/server/database"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func downloadFormFile(request *http.Request, response *http.ResponseWriter, key string, path string) {
	multipart, _, err := request.FormFile(key)
	if err != nil {
		fmt.Println("error receiving file", err)
		http.Error(*response, err.Error(), http.StatusInternalServerError)
		return
	}
	defer multipart.Close()

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("error creating file", err)
		http.Error(*response, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, multipart)
	if err != nil {
		fmt.Println("error writing file", err)
		http.Error(*response, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DownloadStorage(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		file, err := os.Open(fmt.Sprintf("storage/%s", request.PathValue("path")))
		if err != nil {
			http.Error(response, err.Error(), http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, _ := file.Stat()
		if stat.IsDir() {
			http.Error(response, "404 page not found", http.StatusNotFound)
			return
		}

		response.Header().Set("Content-Type", "application/octet-stream")
		response.Header().Set("Content-Length", fmt.Sprint(stat.Size()))
		io.Copy(response, file)
	}
}

func UploadDemo(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		request.ParseMultipartForm(32 << 20)

		id := uuid.New().String()[:8]
		title := request.FormValue("title")

		downloadFormFile(request, &response, "file", fmt.Sprintf("storage/demos/%s.dem", id))

		rows, err := database.Pool.Query(context.Background(), "INSERT INTO videos (id, status, title) values ($1, 0, $2)", id, title)
		if err != nil {
			fmt.Printf("error Query %v\n", err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		rows.Close()

		response.Header().Add("Hx-Redirect", fmt.Sprintf("videos/%s", id))
	}
}

func UploadVideo(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if request.Header.Get("Worker-ID") != "sample" {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}

		request.ParseMultipartForm(32 << 20)

		id := request.FormValue("videoId")
		downloadFormFile(request, &response, "file", fmt.Sprintf("storage/videos/%s.mp4", id))

		rows, err := database.Pool.Query(context.Background(), "UPDATE videos SET status = 1 WHERE id = $1", id)
		if err != nil {
			fmt.Printf("error Query %v\n", err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		rows.Close()
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}

func FetchDemo(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		workerId := request.Header.Get("Worker-ID")

		if workerId != "sample" {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}

		var id string
		err := database.Pool.QueryRow(context.Background(), "SELECT id FROM videos WHERE status = 0").Scan(&id)
		if err != nil {
			if err == pgx.ErrNoRows {
				response.WriteHeader(http.StatusNoContent)
				return
			}
			fmt.Printf("error QueryRow %v\n", err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := database.Pool.Query(context.Background(), "UPDATE videos SET status = 2 WHERE id = $1", id)
		if err != nil {
			fmt.Printf("error Query %v\n", err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		rows.Close()

		response.Write([]byte(id))
	}
}
