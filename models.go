package dotfiled

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	Name     string `sql:"index:name_repo_index"`
	Repo     string `sql:"index:name_repo_index"`
	Contents string `sql:"type:text"`
}

type Chunk struct {
	gorm.Model
	FileType string
	Hash     string `sql:"index"`
	Contents string `sql:"type:text"`
	Tags     string `sql:"type:text[]"`
	Files    []File `gorm:"many2many:file_chunks;"`
}

type FileChunk struct {
	gorm.Model
	FileID  uint
	ChunkID uint
}
