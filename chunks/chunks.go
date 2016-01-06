package chunks

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/charlieegan3/dotfiled"
)

func HashChunk(chunk string) string {
	h := fnv.New32a()
	h.Write([]byte(chunk))
	return fmt.Sprintf("%v", h.Sum32())
}

func ReduceNameToType(name string) string {
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

func TagsForChunk(chunk string, fileType string) string {
	re := regexp.MustCompile("\\W+")
	cleanChunk := string(re.ReplaceAllLiteralString(chunk, " "))
	cleanChunk = strings.ToLower(strings.TrimSpace(cleanChunk))
	return "{" + strings.Join(append(strings.Split(cleanChunk, " "), fileType), ",") + "}"
}

func ValidChunk(chunk string, file models.File) bool {
	re := regexp.MustCompile("^\\W*$")
	if re.MatchString(chunk) {
		return false
	}

	reducedName := ReduceNameToType(file.Name)
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

func FormatChunk(chunk string, file models.File) string {
	reducedName := ReduceNameToType(file.Name)
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
