package models

import "github.com/jinzhu/gorm"

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
	Tags     string `sql:"type:text[]"`
}

type FileChunk struct {
	gorm.Model
	FileID  uint
	ChunkID uint
}
