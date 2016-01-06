package dotfiled

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagsForChunk(t *testing.T) {
	tags := "{set,nocompatible,vim}"
	assert.Equal(t, tags, TagsForChunk("set nocompatible", "vim"), "they should be equal")
	tags = "{shopt,s,histappend,bash}"
	assert.Equal(t, tags, TagsForChunk("shopt -s histappend", "bash"), "they should be equal")
}

func TestValidChunkRejectLackingWordCharacter(t *testing.T) {
	file := File{Name: "bashrc"}
	assert.Equal(t, false, ValidChunk("!@£$%^&*(", file), "they should be equal")
}

func TestValidChunkRejectComments(t *testing.T) {
	file := File{Name: "bashrc"}
	assert.Equal(t, false, ValidChunk("# comment", file), "they should be equal")
	file = File{Name: "vimrc"}
	assert.Equal(t, false, ValidChunk("\" comment", file), "they should be equal")
}

func TestFormatChunkRemovesSemiColon(t *testing.T) {
	file := File{Name: "bashrc"}
	chunk := "shopt -s cdspell;"

	assert.Equal(t, "shopt -s cdspell", FormatChunk(chunk, file), "they should be equal")
}

func TestFormatChunkLeavesValidChunk(t *testing.T) {
	file := File{Name: "bashrc"}
	chunk := "shopt -s cdspell"

	assert.Equal(t, "shopt -s cdspell", FormatChunk(chunk, file), "they should be equal")
}

func TestFormatChunkStripsComment(t *testing.T) {
	file := File{Name: "bashrc"}
	chunk := "shopt -s cdspell; # some comment (blah)"

	assert.Equal(t, "shopt -s cdspell", FormatChunk(chunk, file), "they should be equal")
}
