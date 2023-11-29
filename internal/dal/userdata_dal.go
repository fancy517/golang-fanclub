package dal

import (
	"database/sql"
	"fanclub/internal/models"
	"fmt"
	"strings"
)

type UserdataDAL interface {
	FollowUser(userid string, followid string) (string, error)
	IsFollowing(userid string, followid string) (int, error)
	IsLike(userid string, followid string) (int, error)
	GetPurchases(userid string) ([]models.Purchase, error)
	CreatePurchase(userid string, ptype string, tierid string, amount float64) error
	GetApplicationStatus(username string) (int, error)
	MuteUser(userid string, creator string) (string, error)
	BlockUser(userid string, creator string) (string, error)
	VipUser(userid string, creator string) (string, error)
	CustomUser(userid string, creator string, listid string) (string, error)
	GetMuteList(userid string) ([]models.SimpleUser, error)
	GetBlockList(userid string) ([]models.SimpleUser, error)
	GetVipList(userid string) ([]models.SimpleUser, error)
	GetCustomList(userid string, listid string) (models.CustomUserList, error)
	GetList(userid string) ([]models.CustomList, error)
	GetDefaultEntries(userid string) (models.DefaultEntries, error)
}

type userdataDAL struct {
	DB *sql.DB
}

var _ UserdataDAL = (*userdataDAL)(nil)

func NewUserdataDAL(db *sql.DB) UserdataDAL {
	return &userdataDAL{db}
}

func (dal *userdataDAL) GetDefaultEntries(userid string) (models.DefaultEntries, error) {
	res := models.DefaultEntries{}
	list := models.DefaultEntriesList{}
	query := "SELECT blocklist, mutelist, viplist FROM tbl_userdata WHERE userid =?"
	dal.DB.QueryRow(query, userid).Scan(&list.Blocklist, &list.Mutelist, &list.Viplist)
	if list.Blocklist.Valid {
		values := strings.Split(list.Blocklist.String, ",")
		// Filter out blank strings
		var result []string
		for _, value := range values {
			if value != "" {
				result = append(result, value)
			}
		}
		res.Block = len(result)
	} else {
		res.Block = 0
	}
	if list.Mutelist.Valid {
		values := strings.Split(list.Mutelist.String, ",")
		// Filter out blank strings
		var result []string
		for _, value := range values {
			if value != "" {
				result = append(result, value)
			}
		}
		res.Mute = len(result)
	} else {
		res.Block = 0
	}
	if list.Viplist.Valid {
		values := strings.Split(list.Viplist.String, ",")
		// Filter out blank strings
		var result []string
		for _, value := range values {
			if value != "" {
				result = append(result, value)
			}
		}
		res.Vip = len(result)
	} else {
		res.Vip = 0
	}
	return res, nil
}

func (dal *userdataDAL) GetList(userid string) ([]models.CustomList, error) {
	query := `SELECT
		tbl_customlist.id,
		tbl_customlist.listname,
		tbl_customlist.userlist
	FROM
		tbl_customlist
	WHERE tbl_customlist.userid = ?`
	customlist := make([]models.CustomList, 0)
	rows, err := dal.DB.Query(query, userid)
	if err != nil {
		fmt.Printf("Error querying userdata %v: %v", userid, err)
		return customlist, err
	}
	defer rows.Close()

	for rows.Next() {
		custom := models.CustomList{}
		err := rows.Scan(
			&custom.ID,
			&custom.Listname,
			&custom.UserList,
		)
		if custom.UserList.Valid {
			custom.Entries = len(strings.Split(custom.UserList.String, ","))
		} else {
			custom.Entries = 0
		}
		if err != nil {
			fmt.Printf("error %v\n", err)
			return customlist, err
		}
		customlist = append(customlist, custom)
	}
	return customlist, nil
}

func (dal *userdataDAL) GetCustomList(userid string, listid string) (models.CustomUserList, error) {
	userlist := models.CustomUserList{}
	query := "SELECT listname FROM tbl_customlist WHERE userid = ? AND id = ?"
	err := dal.DB.QueryRow(query, userid, listid).Scan(&userlist.Listname)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return userlist, err
	}

	users := make([]models.SimpleUser, 0)
	query = `SELECT
		users.id,
		users.avatar,
		users.name AS username,
		tbl_userdata.displayname
	FROM
		users
		LEFT JOIN tbl_userdata
		ON users.id = tbl_userdata.userid
	WHERE FIND_IN_SET (
		users.id,
		(SELECT
			tbl_customlist.userlist
		FROM
			tbl_customlist
		WHERE userid = ? AND id = ?)
		)`
	rows, err := dal.DB.Query(query, userid, listid)
	if err != nil {
		fmt.Printf("SQL error: %v\n", err)
		return userlist, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.SimpleUser{}
		rows.Scan(
			&user.ID,
			&user.Avatar,
			&user.Username,
			&user.Displayname,
		)
		users = append(users, user)
	}
	userlist.Userlist = users
	return userlist, nil
}

func (dal *userdataDAL) GetMuteList(userid string) ([]models.SimpleUser, error) {
	users := make([]models.SimpleUser, 0)
	query := `SELECT
		users.id,
		users.avatar,
		users.name AS username,
		tbl_userdata.displayname
	FROM
		users
		LEFT JOIN tbl_userdata
		ON users.id = tbl_userdata.userid
	WHERE FIND_IN_SET (
		users.id,
		(SELECT
			tbl_userdata.mutelist
		FROM
			tbl_userdata
		WHERE userid = ?)
		)`
	rows, err := dal.DB.Query(query, userid)
	if err != nil {
		fmt.Printf("SQL error: %v\n", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.SimpleUser{}
		rows.Scan(
			&user.ID,
			&user.Avatar,
			&user.Username,
			&user.Displayname,
		)
		users = append(users, user)
	}
	return users, nil
}

func (dal *userdataDAL) GetBlockList(userid string) ([]models.SimpleUser, error) {
	users := make([]models.SimpleUser, 0)
	query := `SELECT
		users.id,
		users.avatar,
		users.name AS username,
		tbl_userdata.displayname
	FROM
		users
		LEFT JOIN tbl_userdata
		ON users.id = tbl_userdata.userid
	WHERE FIND_IN_SET (
		users.id,
		(SELECT
			tbl_userdata.blocklist
		FROM
			tbl_userdata
		WHERE userid = ?)
		)`
	rows, err := dal.DB.Query(query, userid)
	if err != nil {
		fmt.Printf("SQL error: %v\n", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.SimpleUser{}
		rows.Scan(
			&user.ID,
			&user.Avatar,
			&user.Username,
			&user.Displayname,
		)
		users = append(users, user)
	}
	return users, nil
}

func (dal *userdataDAL) GetVipList(userid string) ([]models.SimpleUser, error) {
	users := make([]models.SimpleUser, 0)
	query := `SELECT
		users.id,
		users.avatar,
		users.name AS username,
		tbl_userdata.displayname
	FROM
		users
		LEFT JOIN tbl_userdata
		ON users.id = tbl_userdata.userid
	WHERE FIND_IN_SET (
		users.id,
		(SELECT
			tbl_userdata.viplist
		FROM
			tbl_userdata
		WHERE userid = ?)
		)`
	rows, err := dal.DB.Query(query, userid)
	if err != nil {
		fmt.Printf("SQL error: %v\n", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.SimpleUser{}
		rows.Scan(
			&user.ID,
			&user.Avatar,
			&user.Username,
			&user.Displayname,
		)
		users = append(users, user)
	}
	return users, nil
}

func (dal *userdataDAL) MuteUser(userid string, creator string) (string, error) {
	query := "SELECT mutelist FROM tbl_userdata WHERE userid = ?"
	var _mutelist sql.NullString
	err := dal.DB.QueryRow(query, userid).Scan(&_mutelist)
	if err != nil {
		return "failed", err
	}
	mutelist := make([]string, 0)
	var result []string
	return_flag := "add"
	if _mutelist.Valid {
		mutelist = strings.Split(_mutelist.String, ",")
	}
	for _, str := range mutelist {
		if str == creator {
			return_flag = "remove"
			continue
		} else if str != "" {
			result = append(result, str)
		}
	}
	if return_flag == "add" {
		result = append(result, creator)
	}
	query = "UPDATE tbl_userdata SET mutelist = ? WHERE userid = ?"
	_, err = dal.DB.Exec(query, strings.Join(result, ","), userid)
	if err != nil {
		fmt.Printf("Error updating mutelist: %v\n", err)
		return "failed", err
	}
	return return_flag, nil
}

func (dal *userdataDAL) BlockUser(userid string, creator string) (string, error) {
	query := "SELECT blocklist FROM tbl_userdata WHERE userid = ?"
	var _mutelist sql.NullString
	err := dal.DB.QueryRow(query, userid).Scan(&_mutelist)
	if err != nil {
		return "failed", err
	}
	mutelist := make([]string, 0)
	var result []string
	return_flag := "add"
	if _mutelist.Valid {
		mutelist = strings.Split(_mutelist.String, ",")
	}
	for _, str := range mutelist {
		if str == creator {
			return_flag = "remove"
			continue
		} else if str != "" {
			result = append(result, str)
		}
	}
	if return_flag == "add" {
		result = append(result, creator)
	}
	query = "UPDATE tbl_userdata SET blocklist = ? WHERE userid = ?"
	_, err = dal.DB.Exec(query, strings.Join(result, ","), userid)
	if err != nil {
		fmt.Printf("Error updating blocklist: %v\n", err)
		return "failed", err
	}
	return return_flag, nil
}

func (dal *userdataDAL) CustomUser(userid string, creator string, listid string) (string, error) {
	query := "SELECT userlist FROM tbl_customlist WHERE userid = ? and id = ?"
	var _mutelist sql.NullString
	err := dal.DB.QueryRow(query, userid, listid).Scan(&_mutelist)
	if err != nil {
		return "failed", err
	}
	mutelist := make([]string, 0)
	var result []string
	return_flag := "add"
	if _mutelist.Valid {
		mutelist = strings.Split(_mutelist.String, ",")
	}
	for _, str := range mutelist {
		if str == creator {
			return_flag = "remove"
			continue
		} else if str != "" {
			result = append(result, str)
		}
	}
	if return_flag == "add" {
		result = append(result, creator)
	}
	query = "UPDATE tbl_customlist SET userlist = ? WHERE id = ?"
	_, err = dal.DB.Exec(query, strings.Join(result, ","), listid)
	if err != nil {
		fmt.Printf("Error updating viplist: %v\n", err)
		return "failed", err
	}
	return return_flag, nil
}

func (dal *userdataDAL) VipUser(userid string, creator string) (string, error) {
	query := "SELECT viplist FROM tbl_userdata WHERE userid = ?"
	var _mutelist sql.NullString
	err := dal.DB.QueryRow(query, userid).Scan(&_mutelist)
	if err != nil {
		return "failed", err
	}
	mutelist := make([]string, 0)
	var result []string
	return_flag := "add"
	if _mutelist.Valid {
		mutelist = strings.Split(_mutelist.String, ",")
	}
	for _, str := range mutelist {
		if str == creator {
			return_flag = "remove"
			continue
		} else if str != "" {
			result = append(result, str)
		}
	}
	if return_flag == "add" {
		result = append(result, creator)
	}
	query = "UPDATE tbl_userdata SET viplist = ? WHERE userid = ?"
	_, err = dal.DB.Exec(query, strings.Join(result, ","), userid)
	if err != nil {
		fmt.Printf("Error updating viplist: %v\n", err)
		return "failed", err
	}
	return return_flag, nil
}

func (dal *userdataDAL) GetApplicationStatus(username string) (int, error) {
	query := "SELECT status FROM tbl_application WHERE username = ?"
	status := 0
	err := dal.DB.QueryRow(query, username).Scan(&status)
	if err != nil {
		fmt.Printf("%v\n", err)
		return -1, nil
	}
	return status, nil
}

func (dal *userdataDAL) CreatePurchase(userid string, ptype string, value string, amount float64) error {
	creatorid := "0"
	if ptype == "subscribe" {
		// Get Creator ID from tierid
		query := "SELECT users.id FROM tbl_subscriptions LEFT JOIN users ON tbl_subscriptions.username = users.name WHERE tbl_subscriptions.id = ?"
		dal.DB.QueryRow(query, value).Scan(&creatorid)
	} else if ptype == "media" {
		query := "SELECT users.id FROM tbl_medias LEFT JOIN users ON tbl_medias.username = users.name WHERE tbl_medias.id = ?"
		dal.DB.QueryRow(query, value).Scan(&creatorid)
	} else if ptype == "tippost" {
		query := "SELECT userid FROM tbl_post WHERE id = ?"
		dal.DB.QueryRow(query, value).Scan(&creatorid)
	} else if ptype == "tipuser" {
		creatorid = value
	}

	fmt.Printf("creatorid : %v\n", creatorid)
	// Create New Purchase
	query := "INSERT INTO tbl_purchase (userid, type, value, creatorid, amount) VALUES(?, ?, ?, ?, ?)"
	dal.DB.Exec(query, userid, ptype, value, creatorid, amount)
	return nil
}

func (dal *userdataDAL) GetPurchases(userid string) ([]models.Purchase, error) {
	query := `
	SELECT
		tbl_purchase.type,
		tbl_purchase.value,
		tbl_purchase.amount,
		tbl_purchase.created_at,
		users.name AS username,
		users.avatar,
		tbl_userdata.displayname
	FROM
		tbl_purchase
		LEFT JOIN users
		ON tbl_purchase.creatorid = users.id
		LEFT JOIN tbl_userdata
		ON tbl_purchase.creatorid = tbl_userdata.userid
	WHERE tbl_purchase.userid = ?`
	rows, err := dal.DB.Query(query, userid)
	if err != nil {
		fmt.Printf("Error querying %v\n", err)
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Purchase, 0)
	for rows.Next() {
		record := models.Purchase{}
		err := rows.Scan(
			&record.Type,
			&record.Value,
			&record.Amount,
			&record.Date,
			&record.Creator.Username,
			&record.Creator.Avatar,
			&record.Creator.Displayname,
		)
		if err != nil {
			fmt.Printf("Error scanning %v\n", err)
			return nil, err
		}
		res = append(res, record)
	}
	return res, nil
}

func (dal *userdataDAL) IsFollowing(userid string, followid string) (int, error) {
	query := "SELECT COUNT(*) AS cnt FROM tbl_userdata WHERE userid = ? AND FIND_IN_SET(?, likes)"
	var cnt = 0
	err := dal.DB.QueryRow(query, userid, followid).Scan(&cnt)
	if err != nil {
		return cnt, err
	}
	return cnt, nil
}

func (dal *userdataDAL) IsLike(userid string, followid string) (int, error) {
	query := "SELECT COUNT(*) AS cnt FROM tbl_userdata WHERE userid = ? AND FIND_IN_SET(?, followers)"
	var cnt = 0
	err := dal.DB.QueryRow(query, userid, followid).Scan(&cnt)
	if err != nil {
		return cnt, err
	}
	return cnt, nil
}

func (dal *userdataDAL) FollowUser(userid string, followid string) (string, error) {
	fmt.Printf("userid: %s, followid: %s\n", userid, followid)
	query := "SELECT likes FROM tbl_userdata WHERE userid =?"
	var str sql.NullString
	err := dal.DB.QueryRow(query, userid).Scan(&str)
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
		if part == followid {
			return_flag = "remove"
			continue
		} else if part == "" {
			continue
		}
		result = append(result, part)
	}
	if return_flag == "add" {
		result = append(result, followid)
	}
	likes_cnt = len(result)
	query1 := "UPDATE tbl_userdata SET likes =?, likes_cnt =? WHERE userid =?"
	_, err = dal.DB.Exec(query1, strings.Join(result, ","), likes_cnt, userid)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return "", err
	}

	// step 2:
	query = "SELECT followers FROM tbl_userdata WHERE userid =?"
	var str2 sql.NullString
	err = dal.DB.QueryRow(query, followid).Scan(&str2)
	if err != nil {
		fmt.Printf("error3: %v\n", err)
		return "", err
	}
	followers := ""
	if str2.Valid {
		followers = str2.String
	}
	parts2 := strings.Split(followers, ",")
	var follow_cnts = 0
	var result2 []string
	for _, part := range parts2 {
		if part == userid {
			continue
		} else if part == "" {
			continue
		}
		result2 = append(result2, part)
	}
	if return_flag == "add" {
		result2 = append(result2, userid)
	}
	follow_cnts = len(result2)
	query1 = "UPDATE tbl_userdata SET followers =?, followers_cnt =? WHERE userid =?"
	_, err = dal.DB.Exec(query1, strings.Join(result2, ","), follow_cnts, followid)
	if err != nil {
		fmt.Printf("error4: %v\n", err)
		return "", err
	}
	return return_flag, nil
}
