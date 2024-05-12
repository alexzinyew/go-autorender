package routes

import (
	"context"
	"fmt"
	"html/template"
	"lillith/autorender/server/database"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func VideoPage(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	var status int
	var title string
	err := database.Pool.QueryRow(context.Background(), "SELECT status, title FROM videos WHERE id = $1", id).Scan(&status, &title) // no rows found????
	if err != nil {
		if err == pgx.ErrNoRows {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Printf("error QueryRow %v\n", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	page, err := template.ParseFiles("templates/videopage.html")
	if err != nil {
		fmt.Printf("error ParseFiles %v\n", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"title":  title,
		"id":     id,
		"status": fmt.Sprintf("%d", status),
	}

	page.Execute(response, data)
}
