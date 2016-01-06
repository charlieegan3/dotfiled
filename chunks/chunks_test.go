package chunks_test

import (
	"testing"

	"github.com/charlieegan3/dotfiled"
	"github.com/charlieegan3/dotfiled/chunks"
	"github.com/stretchr/testify/assert"
)

func TestTagsForChunk(t *testing.T) {
	tags := "{set,nocompatible,vim}"
	assert.Equal(t, tags, chunks.TagsForChunk("set nocompatible", "vim"), "they should be equal")
	tags = "{shopt,s,histappend,bash}"
	assert.Equal(t, tags, chunks.TagsForChunk("shopt -s histappend", "bash"), "they should be equal")
}

func TestValidChunkRejectLackingWordCharacter(t *testing.T) {
	file := models.File{Name: "bashrc"}
	assert.Equal(t, false, chunks.ValidChunk("!@£$%^&*(", file), "they should be equal")
}

func TestValidChunkRejectComments(t *testing.T) {
	file := models.File{Name: "bashrc"}
	assert.Equal(t, false, chunks.ValidChunk("# comment", file), "they should be equal")
	file = models.File{Name: "vimrc"}
	assert.Equal(t, false, chunks.ValidChunk("\" comment", file), "they should be equal")
}

func TestFormatChunkRemovesSemiColon(t *testing.T) {
	file := models.File{Name: "bashrc"}
	chunk := "shopt -s cdspell;"

	assert.Equal(t, "shopt -s cdspell", chunks.FormatChunk(chunk, file), "they should be equal")
}

func TestFormatChunkLeavesValidChunk(t *testing.T) {
	file := models.File{Name: "bashrc"}
	chunk := "shopt -s cdspell"

	assert.Equal(t, "shopt -s cdspell", chunks.FormatChunk(chunk, file), "they should be equal")
}

func TestFormatChunkStripsComment(t *testing.T) {
	file := models.File{Name: "bashrc"}
	chunk := "shopt -s cdspell; # some comment (blah)"

	assert.Equal(t, "shopt -s cdspell", chunks.FormatChunk(chunk, file), "they should be equal")
}
