package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charlieegan3/dotfiled"
	"github.com/charlieegan3/dotfiled/chunks"
	"github.com/charlieegan3/filechunker"
	"github.com/charlieegan3/repofiles"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var credentials repofiles.Credentials
var db gorm.DB

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
		"https://github.com/durdn/cfg",
		"https://github.com/cypher/dotfiles",
		"https://github.com/henrik/dotfiles",
		"https://github.com/tomasr/dotfiles",
		"https://github.com/zanshin/dotfiles",
		"https://github.com/junegunn/dotfiles",
		"https://github.com/joshuaclayton/dotfiles",
		"https://github.com/matthewmccullough/dotfiles",
		"https://github.com/shawncplus/dotfiles",
		"https://github.com/whiteinge/dotfiles",
		"https://github.com/michaeljsmalley/dotfiles",
		"https://github.com/mscoutermarsh/dotfiles",
		"https://github.com/atomantic/dotfiles",
		"https://github.com/tangledhelix/dotfiles",
		"https://github.com/milkbikis/dotfiles-mac",
		"https://github.com/technomancy/dotfiles",
	}

	db, _ = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db.DropTable(&models.File{}, &models.Chunk{}, &models.FileChunk{})
	db.AutoMigrate(&models.File{}, &models.Chunk{}, &models.FileChunk{})
	credentials = repofiles.Credentials{
		User:  os.Getenv("GITHUB_USER"),
		Token: os.Getenv("GITHUB_TOKEN"),
	}

	pattern := "bashrc|bash_profile|zshrc|vimrc|emacs\\.el|init\\.el|gitignore|gitconfig"
	createMatchingFilesFromRepos(pattern, dotfileRepos)

	var files []models.File
	db.Find(&files)
	createFileChunksForFiles(files)
}

func createOrLinkChunk(chunk string, file models.File, db gorm.DB) models.Chunk {
	reducedName := chunks.ReduceNameToType(file.Name)
	currentChunk := models.Chunk{}
	chunkHash := chunks.HashChunk(chunk)
	db.Where("hash = ? and file_type = ?", chunkHash, reducedName).First(&currentChunk)

	if currentChunk.ID == 0 {
		currentChunk = models.Chunk{
			FileType: reducedName,
			Hash:     chunkHash,
			Contents: chunk,
			Tags:     chunks.TagsForChunk(chunk, reducedName),
		}
		db.Create(&currentChunk)
	}
	return currentChunk
}

func createMatchingFilesFromRepos(pattern string, repos []string) {
	var currentFile models.File
	for _, url := range repos {
		parts := strings.Split(url, "/")
		repo := parts[len(parts)-1]
		user := parts[len(parts)-2]
		fmt.Printf("%v / %v\n", parts[len(parts)-2], parts[len(parts)-1])
		repoData := repofiles.NewRepo(user, repo, "master")
		repoData.List(credentials)
		pattern := "bashrc|bash_profile|zshrc|vimrc|emacs\\.el|init\\.el|gitignore|gitconfig"
		files := repoData.Files(pattern, credentials)
		for i := 0; i < len(files); i++ {
			currentFile = models.File{
				Name:     files[i].Name(),
				Contents: files[i].Contents,
				Repo:     url,
			}
			db.Create(&currentFile)
		}
	}
}

func createFileChunksForFiles(files []models.File) {
	filechunker := filechunker.NewFileChunker(3, "\t")
	var currentChunk models.Chunk
	var currentFileChunk models.FileChunk
	for _, f := range files {
		for _, c := range filechunker.Chunk(f.Contents) {
			c = chunks.FormatChunk(c, f)
			if chunks.ValidChunk(c, f) {
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
