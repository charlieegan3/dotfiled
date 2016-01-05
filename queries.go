package dotfiled

import "github.com/jinzhu/gorm"

type Result struct {
	ID       uint
	Contents string
	FileType string
	Tags     string
	Count    uint
}

func ChunksForQuery(db *gorm.DB, query string, fileType string) []Result {
	var results []Result
	tempDb := db.Table("chunks")
	tempDb = matchingFileTypeParam(tempDb, fileType)
	tempDb = matchingTags(tempDb, query)
	tempDb.Select("chunks.id, chunks.contents, chunks.file_type, chunks.tags, count(file_chunks.id)").
		Joins("inner join file_chunks on file_chunks.chunk_id = chunks.id").
		Group("chunks.id").
		Having("count(file_chunks.id) > 2").
		Order("count(file_chunks.id) desc").
		Limit(100).
		Scan(&results)

	if len(results) == 0 {
		db.Table("chunks").
			Select("chunks.id, chunks.contents, chunks.file_type, chunks.tags, count(file_chunks.id)").
			Joins("inner join file_chunks on file_chunks.chunk_id = chunks.id").
			Where("chunks.contents LIKE ?", "%"+query+"%").
			Group("chunks.id").
			Having("count(file_chunks.id) > 2").
			Order("count(file_chunks.id) desc").
			Limit(100).
			Scan(&results)
	}

	for i, v := range results {
		results[i].Tags = v.Tags[1 : len(v.Tags)-1]
	}

	return results
}

func ChunkForID(db *gorm.DB, id string) Chunk {
	var chunk Chunk
	db.First(&chunk, id)
	db.Model(&chunk).Association("Files").Find(&chunk.Files)
	return chunk
}

func matchingFileTypeParam(db *gorm.DB, fileType string) *gorm.DB {
	if len(fileType) > 0 {
		return db.Where("file_type = ?", fileType)
	} else {
		return db
	}
}

func matchingTags(db *gorm.DB, tags string) *gorm.DB {
	return db.Where("chunks.tags @> ?", "{"+tags+"}")
}
