package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"regexp"
	"strings"

	"github.com/charlieegan3/dotfiled"
	"github.com/charlieegan3/filechunker"
	"github.com/charlieegan3/repofiles"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

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
	db.DropTable(&models.File{}, &models.Chunk{}, &models.FileChunk{})
	db.AutoMigrate(&models.File{}, &models.Chunk{}, &models.FileChunk{})
	credentials := repofiles.Credentials{
		User:  os.Getenv("GITHUB_USER"),
		Token: os.Getenv("GITHUB_TOKEN"),
	}

	var currentFile models.File
	for _, url := range dotfileRepos {
		parts := strings.Split(url, "/")
		repo := parts[len(parts)-1]
		user := parts[len(parts)-2]
		fmt.Printf("%v / %v\n", parts[len(parts)-2], parts[len(parts)-1])
		repoData := repofiles.NewRepo(user, repo, "master")
		repoData.List(credentials)
		files := repoData.Files("bash", credentials)
		files = append(files, repoData.Files("vimrc", credentials)...)
		for _, f := range files {
			currentFile = models.File{Name: f.Name(), Contents: f.Contents}
			db.Create(&currentFile)
		}
	}

	var files []models.File
	db.Find(&files)
	filechunker := filechunker.NewFileChunker(3, "\t")
	var currentChunk models.Chunk
	var currentFileChunk models.FileChunk
	for _, f := range files {
		for _, c := range filechunker.Chunk(f.Contents) {
			if validChunk(c, f) {
				currentChunk = createOrLinkChunk(c, f, db)
				currentFileChunk = models.FileChunk{
					FileID:  f.ID,
					ChunkID: currentChunk.ID,
				}
				db.Create(&currentFileChunk)
			}
		}
	}
}

func hashChunk(chunk string) string {
	h := fnv.New32a()
	h.Write([]byte(chunk))
	return fmt.Sprintf("%v", h.Sum32())
}

func createOrLinkChunk(chunk string, file models.File, db gorm.DB) models.Chunk {
	reducedName := reduceNameToType(file.Name)
	currentChunk := models.Chunk{}
	chunkHash := hashChunk(chunk)
	db.Where("hash = ? and file_type = ?", chunkHash, reducedName).First(&currentChunk)

	if currentChunk.ID == 0 {
		currentChunk = models.Chunk{
			FileType: reducedName,
			Hash:     chunkHash,
			Contents: chunk,
			Tags:     tagsForChunk(chunk, reducedName),
		}
		db.Create(&currentChunk)
	}
	return currentChunk
}

func reduceNameToType(name string) string {
	if strings.Contains(name, "bash") {
		return "bash"
	} else if strings.Contains(name, "vimrc") {
		return "vim"
	} else {
		//TODO
		return name
	}
}

func tagsForChunk(chunk string, fileType string) string {
	re := regexp.MustCompile("\\W+")
	cleanChunk := string(re.ReplaceAllLiteralString(chunk, " "))
	cleanChunk = strings.TrimSpace(cleanChunk)
	return "{" + strings.Join(append(strings.Split(cleanChunk, " "), fileType), ",") + "}"
}

func validChunk(chunk string, file models.File) bool {
	re := regexp.MustCompile("^\\W*$")
	if re.MatchString(chunk) {
		return false
	}

	reducedName := reduceNameToType(file.Name)
	if reducedName == "vim" {
		if chunk[0] == '"' {
			return false
		}
	}
	if reducedName == "bash" {
		if chunk[0] == '#' || chunk == "}" || chunk == "fi" {
			return false
		} else if len(chunk) > 3 {
			if chunk[0:4] == "elif" || chunk[0:4] == "main" {
				return false
			}
		}
	}
	return true
}
