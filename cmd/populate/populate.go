package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"strings"

	"github.com/charlieegan3/filechunker"
	"github.com/charlieegan3/repofiles"
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

func main() {
	dotfileRepos := []string{
		"https://github.com/mathiasbynens/dotfiles",
		"https://github.com/skwp/dotfiles",
		"https://github.com/holman/dotfiles",
		"https://github.com/thoughtbot/dotfiles",
		"https://github.com/ryanb/dotfiles",
		"https://github.com/paulirish/dotfiles",
		"https://github.com/donnemartin/dev-setup",
		"https://github.com/garybernhardt/dotfiles",
		"https://github.com/cowboy/dotfiles",
		"https://github.com/gf3/dotfiles",
		"https://github.com/windelicato/dotfiles",
		"https://github.com/joedicastro/dotfiles",
		"https://github.com/paulmillr/dotfiles",
		"https://github.com/mislav/dotfiles",
		"https://github.com/sontek/dotfiles",
		"https://github.com/necolas/dotfiles",
		"https://github.com/nicknisi/dotfiles",
		"https://github.com/jfrazelle/dotfiles",
		"https://github.com/xero/dotfiles",
		"https://github.com/rmm5t/dotfiles",
		"https://github.com/nelstrom/dotfiles",
		"https://github.com/alrra/dotfiles",
		"https://github.com/dotphiles/dotphiles",
		"https://github.com/tpope/tpope",
		"https://github.com/jferris/config_files",
		"https://github.com/mitsuhiko/dotfiles",
	}

	db, _ := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db.DropTable(&File{}, &Chunk{}, &FileChunk{})
	db.AutoMigrate(&File{}, &Chunk{}, &FileChunk{})

	var currentFile File
	for _, url := range dotfileRepos {
		parts := strings.Split(url, "/")
		repo := parts[len(parts)-1]
		user := parts[len(parts)-2]
		fmt.Printf("%v / %v\n", parts[len(parts)-1], parts[len(parts)-2])
		repoData := repofiles.NewRepo(user, repo, "master")
		repoData.List(repofiles.Credentials{User: os.Getenv("GITHUB_USER"), Token: os.Getenv("GITHUB_TOKEN")})
		files := repoData.Files("vimrc", repofiles.Credentials{User: os.Getenv("GITHUB_USER"), Token: os.Getenv("GITHUB_TOKEN")})
		for _, f := range files {
			currentFile = File{Name: f.Name(), Contents: f.Contents}
			db.Create(&currentFile)
		}
	}

	var files []File
	db.Find(&files)
	filechunker := filechunker.NewFileChunker(3, "\t")
	var currentChunk Chunk
	var currentFileChunk FileChunk
	for _, f := range files {
		for _, c := range filechunker.Chunk(f.Contents) {
			currentChunk = createOrLinkChunk(c, f, db)
			currentFileChunk = FileChunk{FileID: f.ID, ChunkID: currentChunk.ID}
			db.Create(&currentFileChunk)
		}
	}
}

func hashChunk(chunk string) string {
	h := fnv.New32a()
	h.Write([]byte(chunk))
	return fmt.Sprintf("%v", h.Sum32())
}

func createOrLinkChunk(chunk string, file File, db gorm.DB) Chunk {
	currentChunk := Chunk{}
	chunkHash := hashChunk(chunk)
	db.Where("hash = ? and file_type = ?", chunkHash, file.Name).First(&currentChunk)

	if currentChunk.ID == 0 {
		currentChunk = Chunk{FileType: file.Name, Hash: chunkHash, Contents: chunk}
		db.Create(&currentChunk)
	}
	return currentChunk
}
