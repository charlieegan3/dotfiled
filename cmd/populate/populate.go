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

	db, _ := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db.DropTable(&dotfiled.File{}, &dotfiled.Chunk{}, &dotfiled.FileChunk{})
	db.AutoMigrate(&dotfiled.File{}, &dotfiled.Chunk{}, &dotfiled.FileChunk{})
	credentials := repofiles.Credentials{
		User:  os.Getenv("GITHUB_USER"),
		Token: os.Getenv("GITHUB_TOKEN"),
	}

	var currentFile dotfiled.File
	for _, url := range dotfileRepos {
		parts := strings.Split(url, "/")
		repo := parts[len(parts)-1]
		user := parts[len(parts)-2]
		fmt.Printf("%v / %v\n", parts[len(parts)-2], parts[len(parts)-1])
		repoData := repofiles.NewRepo(user, repo, "master")
		repoData.List(credentials)
		var files []repofiles.File
		pattern := "bashrc|bash_profile|zshrc|vimrc|emacs\\.el|init\\.el|gitignore|gitconfig"
		files = append(files, repoData.Files(pattern, credentials)...)
		for _, f := range files {
			currentFile = dotfiled.File{
				Name:     f.Name(),
				Contents: f.Contents,
				Repo:     url,
			}
			db.Create(&currentFile)
		}
	}

	var files []dotfiled.File
	db.Find(&files)
	filechunker := filechunker.NewFileChunker(3, "\t")
	var currentChunk dotfiled.Chunk
	var currentFileChunk dotfiled.FileChunk
	for _, f := range files {
		for _, c := range filechunker.Chunk(f.Contents) {
			c = formatChunk(c, f)
			if validChunk(c, f) {
				currentChunk = createOrLinkChunk(c, f, db)
				currentFileChunk = dotfiled.FileChunk{
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

func createOrLinkChunk(chunk string, file dotfiled.File, db gorm.DB) dotfiled.Chunk {
	reducedName := reduceNameToType(file.Name)
	currentChunk := dotfiled.Chunk{}
	chunkHash := hashChunk(chunk)
	db.Where("hash = ? and file_type = ?", chunkHash, reducedName).First(&currentChunk)

	if currentChunk.ID == 0 {
		currentChunk = dotfiled.Chunk{
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
	} else if strings.Contains(name, "zsh") {
		return "zsh"
	} else if strings.Contains(name, "emacs") || strings.Contains(name, ".el") {
		return "emacs"
	} else if strings.Contains(name, "gitignore") {
		return "gitignore"
	} else if strings.Contains(name, "gitconfig") {
		return "gitconfig"
	} else {
		return name
	}
}

func tagsForChunk(chunk string, fileType string) string {
	re := regexp.MustCompile("\\W+")
	cleanChunk := string(re.ReplaceAllLiteralString(chunk, " "))
	cleanChunk = strings.ToLower(strings.TrimSpace(cleanChunk))
	return "{" + strings.Join(append(strings.Split(cleanChunk, " "), fileType), ",") + "}"
}

func validChunk(chunk string, file dotfiled.File) bool {
	re := regexp.MustCompile("^\\W*$")
	if re.MatchString(chunk) {
		return false
	}

	reducedName := reduceNameToType(file.Name)
	if reducedName == "vim" {
		if chunk[0] == '"' {
			return false
		}
		if chunk == "endif" || chunk == "endfunction" {
			return false
		}
	}
	if reducedName == "bash" || reducedName == "zsh" {
		if chunk[0] == '#' || chunk == "}" || chunk == "fi" {
			return false
		} else if chunk[len(chunk)-1] == '{' {
			return false
		} else if len(chunk) > 3 {
			if chunk[0:4] == "elif" || chunk[0:2] == "if" || chunk[0:4] == "main" {
				return false
			}
		}
	}
	if reducedName == "emacs" {
		if len(chunk) > 1 {
			if chunk[0:2] == ";;" {
				return false
			}
		}
	}
	if reducedName == "gitconfig" || reducedName == "gitignore" {
		if chunk[0] == '#' {
			return false
		} else if chunk[len(chunk)-1] == ']' {
			return false
		}
	}
	return true
}

func formatChunk(chunk string, file dotfiled.File) string {
	reducedName := reduceNameToType(file.Name)
	if reducedName == "bash" {
		re := regexp.MustCompile("#(\\w|\\s||[^\"';])*$")
		chunk = re.ReplaceAllLiteralString(chunk, "")
		chunk = strings.TrimSpace(chunk)
		if len(chunk) > 0 && chunk[len(chunk)-1] == ';' {
			chunk = chunk[0 : len(chunk)-1]
		}
	}
	return chunk
}
