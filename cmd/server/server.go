package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/charlieegan3/dotfiled"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db gorm.DB

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		ID       uint
		Contents string
		FileType string
		Count    uint
	}
	var results []Result
	db.Table("chunks").
		Select("chunks.id, chunks.contents, chunks.file_type, count(file_chunks.id)").
		Joins("inner join file_chunks on file_chunks.chunk_id = chunks.id").
		Group("chunks.id").
		Having("count(file_chunks.id) > 3").
		Order("count(file_chunks.id) desc").
		Scan(&results)

	jsonString, _ := json.Marshal(results)
	io.WriteString(w, string(jsonString))
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Missing PORT parameter")
	} else if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("Missing DATABASE_URL environment variable")
	}
	port := os.Args[1]
	fmt.Printf("Listening on port %v\n", port)

	db, _ = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db.AutoMigrate(&models.File{}, &models.Chunk{}, &models.FileChunk{})

	http.HandleFunc("/", index)
	http.ListenAndServe(":"+port, nil)
}
