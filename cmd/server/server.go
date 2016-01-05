package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/charlieegan3/dotfiled"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db gorm.DB

func ChunksIndexHandler(w http.ResponseWriter, r *http.Request) {
	template.
		Must(template.ParseFiles("templates/index.html", "templates/base.html")).
		ExecuteTemplate(w, "base", nil)
}

func ChunksShowHandler(w http.ResponseWriter, r *http.Request) {
	data := struct{ ID string }{r.URL.Path[len("/chunks/"):]}
	template.
		Must(template.ParseFiles("templates/show.html", "templates/base.html")).
		ExecuteTemplate(w, "base", data)
}

func ApiChunksIndexHandler(w http.ResponseWriter, r *http.Request) {
	results := dotfiled.ChunksForQuery(
		&db, r.URL.Query().Get("tags"), r.URL.Query().Get("file_type"))

	jsonString, _ := json.Marshal(results)
	io.WriteString(w, string(jsonString))
}

func ApiChunkShowHandler(w http.ResponseWriter, r *http.Request) {
	result := dotfiled.ChunkForID(&db, r.URL.Path[len("/api/chunks/"):])
	jsonString, _ := json.Marshal(result)
	io.WriteString(w, string(jsonString))
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Missing PORT parameter")
	} else if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("Missing DATABASE_URL environment variable")
	}

	db, _ = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db.AutoMigrate(&dotfiled.File{}, &dotfiled.Chunk{}, &dotfiled.FileChunk{})
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", ChunksIndexHandler)
	http.HandleFunc("/chunks/", ChunksShowHandler)
	http.HandleFunc("/api/chunks", ApiChunksIndexHandler)
	http.HandleFunc("/api/chunks/", ApiChunkShowHandler)
	http.ListenAndServe(":"+os.Args[1], nil)
}
