package dal

import (
	"database/sql"
	"path/filepath"
	"strings"
)

type MediaDAL interface {
	CreateMedia(filename string, tags string, username string) (string, error)
	GetMediaCounts(username string) ([]string, error)
}

type mediaDAL struct {
	DB *sql.DB
}

var _ MediaDAL = (*mediaDAL)(nil)

func NewMediaDAL(db *sql.DB) MediaDAL {
	return &mediaDAL{db}
}

func (dal *mediaDAL) GetMediaCounts(username string) ([]string, error) {
	query := "SELECT count(*) as count, type FROM tbl_medias WHERE username = ? GROUP BY type ORDER BY type"
	rows, err := dal.DB.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cnts []string
	idx := 0
	for rows.Next() {
		var cnt string
		var mtype string
		err = rows.Scan(&cnt, &mtype)
		if err != nil {
			return nil, err
		}
		if( mtype == "image" && idx == 0) || (mtype =="video" && idx ==1){
			cnts = append(cnts, cnt)
		} else{
			cnts = append(cnts, "0")
			cnts = append(cnts, cnt)
		}
		idx++
	}
	if idx == 0 {
		cnts = append(cnts, "0")
		cnts = append(cnts, "0")
	} else if idx == 1 && len(cnts) == 1{
		cnts = append(cnts, "0")
	}
	return cnts, nil
}

func (dal *mediaDAL) CreateMedia(filename string, tags string, username string) (string, error) {
	query := `
    INSERT INTO tbl_medias(filename,type,tags, username)
    VALUES(?,?,?,?)
    RETURNING id;
    `
	extension := strings.ToLower(filepath.Ext(filename))
	file_ext := "image"
	switch extension {
	case ".mp4", ".mov", ".avi", ".wmv", ".mkv", ".flv", ".mpg", ".mpeg":
		file_ext = "video"
	default:
		file_ext = "image"
	}
	var id string
	if err := dal.DB.QueryRow(query, filename, file_ext, tags, username).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
