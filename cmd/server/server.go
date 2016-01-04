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

func ApiChunksIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		ID       uint
		Contents string
		FileType string
		Tags     string
		Count    uint
	}
	var results []Result
	tempDb := db.Table("chunks")
	tempDb = matchingFileTypeParam(tempDb, r.URL.Query().Get("file_type"))
	tempDb = matchingTags(tempDb, r.URL.Query().Get("tags"))
	tempDb.Select("chunks.id, chunks.contents, chunks.file_type, chunks.tags, count(file_chunks.id)").
		Joins("inner join file_chunks on file_chunks.chunk_id = chunks.id").
		Group("chunks.id").
		Having("count(file_chunks.id) > 2").
		Order("count(file_chunks.id) desc").
		Limit(100).
		Scan(&results)

	if len(results) == 0 {
		db.Table("chunks").
			Select("chunks.id, chunks.contents, chunks.file_type, chunks.tags, count(file_chunks.id)").
			Joins("inner join file_chunks on file_chunks.chunk_id = chunks.id").
			Where("chunks.contents LIKE ?", "%"+r.URL.Query().Get("tags")+"%").
			Group("chunks.id").
			Having("count(file_chunks.id) > 2").
			Order("count(file_chunks.id) desc").
			Limit(100).
			Scan(&results)
	}

	for i, v := range results {
		results[i].Tags = v.Tags[1 : len(v.Tags)-1]
	}

	jsonString, _ := json.Marshal(results)
	io.WriteString(w, string(jsonString))
}

func ApiChunkShowHandler(w http.ResponseWriter, r *http.Request) {
	var chunk models.Chunk
	db.First(&chunk, r.URL.Path[len("/api/chunks/"):])
	db.Model(&chunk).Related(&chunk.FileChunks)
	jsonString, _ := json.Marshal(chunk)
	io.WriteString(w, string(jsonString))
}

func ChunksShowHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/chunks/"):]
	t, _ := template.ParseFiles("chunksShow.html")
	type ID struct {
		ID string
	}
	t.Execute(w, ID{ID: id})
}

func matchingFileTypeParam(db *gorm.DB, fileType string) *gorm.DB {
	if len(fileType) > 0 {
		return db.Where("file_type = ?", fileType)
	} else {
		return db
	}
}

func matchingTags(db *gorm.DB, tags string) *gorm.DB {
	return db.Where("chunks.tags @> ?", "{"+tags+"}")
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Missing PORT parameter")
	} else if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("Missing DATABASE_URL environment variable")
	}

	db, _ = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db.AutoMigrate(&models.File{}, &models.Chunk{}, &models.FileChunk{})
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/", fs)
	http.HandleFunc("/chunks/", ChunksShowHandler)
	http.HandleFunc("/api/chunks", ApiChunksIndexHandler)
	http.HandleFunc("/api/chunks/", ApiChunkShowHandler)
	http.ListenAndServe(":"+os.Args[1], nil)
}
