package dal

import (
	"database/sql"
	"fanclub/internal/models"
	"fmt"
	"strings"
)

type PostlistDAL interface {
	CreatePost(models.Post) error
	UpdatePost(models.Post, string) error
	RemovePost(postid string, userid string) error
	GetPostData(username string, loginuser string, parentid string, selecttag string) ([]models.Postdata, error)
	GetOnePost(postid string, username string) (models.Postdata, error)
	LikePost(postid string, username string) (string, error)
	GetTags(username string) (string, error)
	IncreaseCommentCount(postid int) error
	GetUserIDFromPostid(postid int) (int, error)
}

type postlistDAL struct {
	DB *sql.DB
}

var _ PostlistDAL = (*postlistDAL)(nil)

func NewPostlistDAL(db *sql.DB) PostlistDAL {
	return &postlistDAL{db}
}

func (dal *postlistDAL) GetUserIDFromPostid(postid int) (int, error) {
	query := "SELECT userid FROM tbl_post WHERE id = ?"
	var userid int
	err := dal.DB.QueryRow(query, postid).Scan(&userid)
	if err != nil {
		fmt.Printf("Error getting userid: %v\n", err)
		return 0, err
	}
	return userid, nil
}

func (dal *postlistDAL) IncreaseCommentCount(postid int) error {
	query := "UPDATE tbl_post SET comments_cnt = comments_cnt + 1 WHERE id = ?"
	_, err := dal.DB.Exec(query, postid)
	if err != nil {
		return err
	}
	return nil
}

func (dal *postlistDAL) GetTags(username string) (string, error) {

	if username != "" {
		query := "SELECT GROUP_CONCAT(tags) AS tags FROM (SELECT  tbl_medias.tags FROM tbl_medias WHERE FIND_IN_SET(tbl_medias.id,(SELECT GROUP_CONCAT(tbl_post.files SEPARATOR ',') AS files FROM tbl_post WHERE tbl_post.name = ? AND tbl_post.files != '' AND tbl_post.parent_id = 0 )) AND tbl_medias.tags != '') AS A"
		var tags = ""
		err := dal.DB.QueryRow(query, username).Scan(&tags)
		if err != nil {
			return "", err
		}
		return tags, nil
	} else {
		query := "SELECT tags FROM tbl_medias"
		rows, err := dal.DB.Query(query)
		if err != nil {
			return "", err
		}
		defer rows.Close()
		var tags []string
		for rows.Next() {
			var tag sql.NullString
			err = rows.Scan(&tag)
			if err != nil {
				return "", err
			}
			if tag.Valid && tag.String != "" {
				tags = append(tags, tag.String)
			}

		}
		return strings.Join(tags, ","), nil
	}
}

func (dal *postlistDAL) LikePost(postid string, username string) (string, error) {
	query := "SELECT likes FROM tbl_post WHERE id =?"
	var str sql.NullString
	err := dal.DB.QueryRow(query, postid).Scan(&str)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
		return "", err
	}
	likes := ""
	if str.Valid {
		likes = str.String
	}
	parts := strings.Split(likes, ",")
	var result []string
	var return_flag = "add"
	var likes_cnt = 0
	for _, part := range parts {
		if part == username {
			return_flag = "remove"
			continue
		} else if part == "" {
			continue
		}
		result = append(result, part)
	}
	if return_flag == "add" {
		result = append(result, username)
	}
	likes_cnt = len(result)
	query1 := "UPDATE tbl_post SET likes =?, likes_cnt =? WHERE id =?"
	_, err = dal.DB.Exec(query1, strings.Join(result, ","), likes_cnt, postid)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return "", err
	}
	return return_flag, nil
}

func (dal *postlistDAL) CreatePost(post models.Post) error {
	insertQuery := `
		INSERT INTO tbl_post (userid, name,parent_id, message, is_pinned, reply_role, files, publish_date, disappear_date) VALUES	(?,?,?,?,?,?,?,?,?)`
	_, err := dal.DB.Exec(insertQuery, post.UserID, post.Name, post.ParentID, post.Message, post.IsPinned, post.ReplyRole, post.Files, post.PublishDate, post.DisappearDate)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return err
	}
	return nil
}

func (dal *postlistDAL) UpdatePost(post models.Post, postid string) error {
	updateQuery := `
        UPDATE tbl_post SET message =?, is_pinned =?, reply_role =?, files =?, publish_date =?, disappear_date =? WHERE id =?`
	_, err := dal.DB.Exec(updateQuery, post.Message, post.IsPinned, post.ReplyRole, post.Files, post.PublishDate, post.DisappearDate, postid)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return err
	}
	return nil
}

func (dal *postlistDAL) RemovePost(postid string, userid string) error {
	fmt.Printf("%v\n%v\n", postid, userid)
	var parentid string
	query := "SELECT parent_id FROM tbl_post WHERE id =?"
	err := dal.DB.QueryRow(query, postid).Scan(&parentid)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
		return err
	}
	if parentid != "0" {
		query := "UPDATE tbl_post SET comments_cnt = comments_cnt - 1 WHERE id = ?"
		_, err = dal.DB.Exec(query, parentid)
		if err != nil {
			fmt.Printf("error2: %v\n", err)
			return err
		}
	}

	fmt.Printf("delete: %v\n", postid)
	deleteQuery := "DELETE FROM tbl_post WHERE id =? and userid =?"
	_, err = dal.DB.Exec(deleteQuery, postid, userid)
	if err != nil {
		fmt.Printf("error in remove post: %v\n", err)
		return err
	}
	fmt.Printf("%v\n", postid)
	return nil
}

func (dal *postlistDAL) GetOnePost(postid string, username string) (models.Postdata, error) {
	if username == "" {
		username = "----iazfe"
	}

	var insertQuery = `SELECT
	tbl_userdata.displayName,
	users.name AS userName,
	users.status AS availability,
	users.avatar AS avatarUrl,
	tbl_userdata.userBanner AS bannerUrl,
	users.active AS verified,
	tbl_post.message AS description,
	tbl_post.reply_role,
	tbl_post.publish_date,
	tbl_post.disappear_date,
	tbl_post.likes_cnt,
	tbl_post.files,
	tbl_post.comments_cnt,
	CASE 
	WHEN tbl_post.likes LIKE ? THEN 1
	ELSE 0
	END AS is_liked,
	tbl_post.id AS id,
	tbl_post.is_pinned
  FROM
	tbl_post
	LEFT JOIN tbl_userdata
	  ON tbl_post.userid = tbl_userdata.userid
	LEFT JOIN users
	  ON tbl_post.userid = users.id 
  WHERE tbl_post.id =?`
	item := models.Postdata{}
	err := dal.DB.QueryRow(insertQuery, "%"+username+"%", postid).Scan(
		&item.Creator.DisplayName,
		&item.Creator.UserName,
		&item.Creator.Availability,
		&item.Creator.AvatarUrl,
		&item.Creator.BannerUrl,
		&item.Creator.Verified,
		&item.Description,
		&item.ReplyRole,
		&item.PublishDate,
		&item.DisappearDate,
		&item.LikesCnt,
		&item.Files,
		&item.CommentsCnt,
		&item.IsLiked,
		&item.ID,
		&item.IsPinned,
	)
	if err != nil {
		return models.Postdata{}, err
	}
	result := strings.Split(item.Files, ",")

	medias := make([]models.Medias, 0)
	if len(result) > 0 && result[0] != "" {
		for index, val := range result {
			newmedia := models.Medias{}
			query := `SELECT filename, type, tags from tbl_medias WHERE id=?`
			err := dal.DB.QueryRow(query, val).Scan(
				&newmedia.SourceId,
				&newmedia.Type,
				&newmedia.Tags,
			)
			if err != nil {
				fmt.Printf("error3: %v\n", err)
				return models.Postdata{}, err
			}
			newmedia.MediaId = val
			newmedia.Sensetive = "False"
			newmedia.Loked = "False"
			newmedia.Timestamp = "100"
			medias = append(medias, newmedia)
			index++
		}
	}
	item.Attachment.Medias = medias
	return item, nil
}

func (dal *postlistDAL) GetPostData(username string, loginuser string, parentid string, selecttag string) ([]models.Postdata, error) {

	query := `SELECT id FROM users WHERE users.name = ?`
	var userID int
	if parentid == "0" {
		err := dal.DB.QueryRow(query, username).Scan(&userID)
		if err != nil {
			fmt.Printf("error1: %v\n", err)
			return nil, err
		}
	}
	var insertQuery = `SELECT
	tbl_userdata.displayName,
	users.name AS userName,
	users.status AS availability,
	users.avatar AS avatarUrl,
	tbl_userdata.userBanner AS bannerUrl,
	users.active AS verified,
	tbl_post.message AS description,
	tbl_post.reply_role,
	tbl_post.publish_date,
	tbl_post.disappear_date,
	tbl_post.likes_cnt,
	tbl_post.files,
	tbl_post.comments_cnt,
	CASE 
	WHEN tbl_post.likes LIKE ? THEN 1
	ELSE 0
	END AS is_liked,
	tbl_post.id AS id,
	tbl_post.is_pinned
  FROM
	tbl_post
	LEFT JOIN tbl_userdata
	  ON tbl_post.userid = tbl_userdata.userid
	LEFT JOIN users
	  ON tbl_post.userid = users.id `
	var rows *sql.Rows
	var err error
	fmt.Printf("parentid: %v\n", parentid)
	if parentid == "0" {
		if selecttag == "" {
			insertQuery = insertQuery + ` WHERE tbl_post.userid = ? AND tbl_post.parent_id = ? ORDER BY tbl_post.is_pinned DESC, tbl_post.created_at DESC`
			rows, err = dal.DB.Query(insertQuery, loginuser, userID, parentid)
		} else {
			insertQuery = insertQuery + ` WHERE tbl_post.userid = ? AND tbl_post.parent_id = ? 
			AND EXISTS
			(SELECT
				1
			FROM
				(SELECT
				tbl_medias.id AS fileid
				FROM
				tbl_medias
				WHERE FIND_IN_SET (?, tbl_medias.tags)) AS A
			WHERE FIND_IN_SET (A.fileid, tbl_post.files) > 0)
			ORDER BY tbl_post.is_pinned DESC, tbl_post.created_at DESC`
			rows, err = dal.DB.Query(insertQuery, loginuser, userID, parentid, selecttag)
		}
	} else {
		insertQuery = insertQuery + `WHERE tbl_post.parent_id = ? ORDER BY tbl_post.is_pinned DESC, tbl_post.created_at DESC`
		rows, err = dal.DB.Query(insertQuery, loginuser, parentid)

	}
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	postdata := make([]models.Postdata, 0)
	for rows.Next() {
		item := models.Postdata{}
		err := rows.Scan(
			&item.Creator.DisplayName,
			&item.Creator.UserName,
			&item.Creator.Availability,
			&item.Creator.AvatarUrl,
			&item.Creator.BannerUrl,
			&item.Creator.Verified,
			&item.Description,
			&item.ReplyRole,
			&item.PublishDate,
			&item.DisappearDate,
			&item.LikesCnt,
			&item.Files,
			&item.CommentsCnt,
			&item.IsLiked,
			&item.ID,
			&item.IsPinned,
		)
		if err != nil {
			return nil, err
		}

		result := strings.Split(item.Files, ",")

		medias := make([]models.Medias, 0)
		if len(result) > 0 && result[0] != "" {
			for index, val := range result {
				newmedia := models.Medias{}
				query := `SELECT filename, type, tags from tbl_medias WHERE id=?`
				err := dal.DB.QueryRow(query, val).Scan(
					&newmedia.SourceId,
					&newmedia.Type,
					&newmedia.Tags,
				)
				if err != nil {
					fmt.Printf("error3: %v\n", err)
					return nil, err
				}
				newmedia.MediaId = val
				newmedia.Sensetive = "False"
				newmedia.Loked = "False"
				newmedia.Timestamp = "100"
				medias = append(medias, newmedia)
				index++
			}
		}
		item.Attachment.Medias = medias
		postdata = append(postdata, item)
	}
	return postdata, nil
}
