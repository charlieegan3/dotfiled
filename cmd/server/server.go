package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type File struct {
	gorm.Model
	Name     string
	Contents string `sql:"type:text"`
}

type Chunk struct {
	gorm.Model
	FileType string
	Hash     string `sql:"index"`
	Contents string `sql:"type:text"`
}

type FileChunk struct {
	gorm.Model
	FileID  uint
	ChunkID uint
}

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
	port := os.Getenv("PORT")
	fmt.Printf("Listening on port %v\n", port)

	db, _ = gorm.Open("postgres", "user="+os.Getenv("DB_USER")+" dbname="+os.Getenv("DB_NAME")+" sslmode="+os.Getenv("DB_SSL_MODE"))

	http.HandleFunc("/", index)
	http.ListenAndServe(":"+port, nil)
}
