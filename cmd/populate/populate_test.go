package main

import (
	"testing"

	"github.com/charlieegan3/dotfiled"
	"github.com/stretchr/testify/assert"
)

func TestTagsForChunk(t *testing.T) {
	tags := "{set,nocompatible,vim}"
	assert.Equal(t, tags, tagsForChunk("set nocompatible", "vim"), "they should be equal")
	tags = "{shopt,s,histappend,bash}"
	assert.Equal(t, tags, tagsForChunk("shopt -s histappend", "bash"), "they should be equal")
}

func TestValidChunkRejectLackingWordCharacter(t *testing.T) {
	file := models.File{Name: "bashrc"}
	assert.Equal(t, false, validChunk("!@£$%^&*(", file), "they should be equal")
}

func TestValidChunkRejectComments(t *testing.T) {
	file := models.File{Name: "bashrc"}
	assert.Equal(t, false, validChunk("# comment", file), "they should be equal")
	file = models.File{Name: "vimrc"}
	assert.Equal(t, false, validChunk("\" comment", file), "they should be equal")
}

func TestFormatChunkRemovesSemiColon(t *testing.T) {
	file := models.File{Name: "bashrc"}
	chunk := "shopt -s cdspell;"

	assert.Equal(t, "shopt -s cdspell", formatChunk(chunk, file), "they should be equal")
}

func TestFormatChunkLeavesValidChunk(t *testing.T) {
	file := models.File{Name: "bashrc"}
	chunk := "shopt -s cdspell"

	assert.Equal(t, "shopt -s cdspell", formatChunk(chunk, file), "they should be equal")
}

func TestFormatChunkStripsComment(t *testing.T) {
	file := models.File{Name: "bashrc"}
	chunk := "shopt -s cdspell; # some comment (blah)"

	assert.Equal(t, "shopt -s cdspell", formatChunk(chunk, file), "they should be equal")
}
