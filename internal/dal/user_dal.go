package dal

import (
	"database/sql"
	"fanclub/internal/models"
	"fanclub/internal/types"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserDAL interface {
	GetByName(name string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByID(userID int) (models.User, error)
	GetByActivationCode(code string) (models.User, error)
	GetAll(types.UserFilter) ([]models.User, int, error)
	GetDisplayNameByID(id int) (string, error)
	InsertGetID(models.User) error
	Activate(userID int) error
	UpdatePassword(userID int, password string) error
	UpdateProfile(userID int, user models.User) error
	Block(userID int, status types.UserStatusType) error
	SetPasswordResetCode(userID int, code string) error
	GetUserIDFromPasswordResetCode(code string) (int, time.Time, error)
	DeletePasswordResetCode(code string) error
	UpdateActivationCode(userID int, code string) error
	ConfirmActivationCode(userID int, code string) (string, error)
	IsActivated(userID int) (bool, error)
	// GetByToken(token string) (models.CompleteUser, error)
	GetUserDataByName(name string) (models.CompleteUser, error)
	UpdateAvatar(userID int, avatar string) error
	UpdateBanner(userID int, banner string) error
	UpdateProfile1(profile models.Profile) error
	SubmitApplication(models.Application) (string, error)
	GetWaitingApplications() ([]models.ApplicationUser, error)
	GetApplicationByUsername(username string) (models.Application, error)
	AcceptApplication(username string) error
	RejectApplication(username string) error
	SetOnline(userid string) error
	SetAway(userid string) error
	SetHidden(userid string) error
	CreateList(userid string, name string) error
	EditList(userid string, name string, listid string) error
	DeleteList(userid string, listid string) error

	//Settings interface
	ChangeUsername(string, string, string) (string, error)
	ChangeDisplayname(int, string) error
	CheckPassword(string, string) error
	CheckVCode1(string, string, string) (string, error)
	CheckVCode2(string, string, string) (string, error)
	GetEmailByUsername(string) (string, error)
	RequestVCodeEmail(string, string) error
	RequestVCodeEmail2(string, string) error
	SetTwoFactorEnabled(string, string) error
	SetTwoFactorDisabled(string) error
	GetOtpSecret(string, string) (string, error)
	ChangePassword(string, string) error
	SavePermissions([][]interface{}) ([]string, error)
	SaveSafety(username string, content int, location string, permids string, message int) error
	GetSafety(username string) (models.Safety, error)
	GetSessions(username string) ([]models.Session, error)
	CloseSession(username string, sid string) error
	ClearSessions(username string) error
	GetWalletAddr(userid int) (string, error)
	InsertSession(string, string, string) error
	DeleteAccount(userid string) error

	// Tips Interface
	AddTipHistory(userid int, ttype string, value int, amount float64, description string) error
	// Notifications Setting Interface
	AddNotificationSettings(userid int) error
	GetNotificationSettings(userid int) (models.NotificationSettings, error)
	UpdateNotificationSettings(userid int, is_push int, is_message int, is_reply int, is_postlike int, is_follower int) error
}

type userDAL struct {
	DB *sql.DB
}

var _ UserDAL = (*userDAL)(nil)

func NewUserDAL(db *sql.DB) UserDAL {
	return &userDAL{
		DB: db,
	}
}

const (
	BcryptCost = 15
)

// Implementations

func (dal *userDAL) EditList(userid string, name string, listid string) error {
	query := "UPDATE tbl_customlist SET listname = ? WHERE userid = ? AND id =?"
	_, err := dal.DB.Exec(query, name, userid, listid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) DeleteList(userid string, listid string) error {
	query := "DELETE FROM tbl_customlist WHERE userid = ? AND id = ?"
	_, err := dal.DB.Exec(query, userid, listid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) CreateList(userid string, name string) error {
	query := "INSERT INTO tbl_customlist (userid, listname) VALUES(?,?)"
	_, err := dal.DB.Exec(query, userid, name)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) DeleteAccount(userid string) error {
	query := "UPDATE users SET deleted_at = NOW() WHERE id = ?"
	_, err := dal.DB.Exec(query, userid)
	if err != nil {
		fmt.Printf("Error Deleting Account: %v\n", err)
		return err
	}
	return nil
}

// Account Status
func (dal *userDAL) SetOnline(userid string) error {
	query := "UPDATE users SET status = 'online' WHERE id = ?"
	_, err := dal.DB.Exec(query, userid)
	if err != nil {
		fmt.Printf("Error Updating %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) SetAway(userid string) error {
	query := "UPDATE users SET status = 'away' WHERE id = ?"
	_, err := dal.DB.Exec(query, userid)
	if err != nil {
		fmt.Printf("Error Updating %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) SetHidden(userid string) error {
	query := "UPDATE users SET status = 'invisible' WHERE id = ?"
	_, err := dal.DB.Exec(query, userid)
	if err != nil {
		fmt.Printf("Error Updating %v\n", err)
		return err
	}
	return nil
}

// Notifications Interface
func (dal *userDAL) AddNotificationSettings(userid int) error {
	query := "INSERT INTO tbl_notification_setting (userid) VALUES(?)"
	_, err := dal.DB.Exec(query, userid)
	if err != nil {
		fmt.Printf("Error Inserting %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) GetNotificationSettings(userid int) (models.NotificationSettings, error) {
	nsetting := models.NotificationSettings{}
	query := "SELECT userid, is_push, is_message, is_reply, is_postlike, is_follower FROM tbl_notification_setting WHERE userid = ?"
	err := dal.DB.QueryRow(query, userid).Scan(
		&nsetting.UserID,
		&nsetting.IsPush,
		&nsetting.IsMessage,
		&nsetting.IsReply,
		&nsetting.IsPostlike,
		&nsetting.IsFollower,
	)
	if err != nil {
		fmt.Printf("Error Inserting %v\n", err)
		return nsetting, err
	}
	return nsetting, nil
}

func (dal *userDAL) UpdateNotificationSettings(userid int, is_push int, is_message int, is_reply int, is_postlike int, is_follower int) error {
	query := "UPDATE tbl_notification_setting SET is_push= ?, is_message = ?, is_reply = ?, is_postlike = ?, is_follower = ? WHERE userid = ? "
	_, err := dal.DB.Exec(query, is_push, is_message, is_reply, is_postlike, is_follower, userid)
	if err != nil {
		fmt.Printf("Error Inserting %v\n", err)
		return err
	}
	return nil
}

// Tips
func (dal *userDAL) AddTipHistory(userid int, ttype string, value int, amount float64, description string) error {
	query := "INSERT INTO tbl_tips (userid, type, value, amount, description) VALUES(?,?,?,?,?)"
	_, err := dal.DB.Exec(query, userid, ttype, value, amount, description)
	if err != nil {
		fmt.Printf("Error Inserting Tips: %v", err)
		return err
	}
	return nil
}

// Settings

func (dal *userDAL) InsertSession(username string, ip string, location string) error {
	fmt.Printf("%v\n, %v\n, %v\n", username, ip, location)

	query := "SELECT count(*) as cnt FROM tbl_sessions WHERE username = ? AND last_used_ip = ?"
	cnt := 0
	err := dal.DB.QueryRow(query, username, ip).Scan(&cnt)
	if err != nil {
		fmt.Printf("Error counting : %v\n", err)
		return err
	}
	if cnt == 0 {
		query = "INSERT INTO tbl_sessions (username, last_used_time, last_used_ip, last_used_location) VALUES( ?, NOW(), ?, ?)"
		_, err := dal.DB.Exec(query, username, ip, location)
		if err != nil {
			fmt.Printf("Error executing query %v: %v", query, err)
			return err
		}
	} else {
		query = "UPDATE tbl_sessions SET last_used_time = NOW(), last_used_location = ? WHERE username = ? AND last_used_ip = ?"
		_, err := dal.DB.Exec(query, location, username, ip)
		if err != nil {
			fmt.Printf("Error executing query %v: %v", query, err)
			return err
		}
	}
	return nil
}

func (dal *userDAL) GetWalletAddr(userid int) (string, error) {
	query := "SELECT deposit_address FROM wallets WHERE user_id = ?"
	var wallet string
	err := dal.DB.QueryRow(query, userid).Scan(&wallet)
	if err != nil {
		return "", err
	}
	return wallet, nil
}

func (dal *userDAL) GetSessions(username string) ([]models.Session, error) {
	query := "SELECT id, last_used_time, last_used_ip, last_used_location FROM tbl_sessions WHERE username = ?"
	rows, err := dal.DB.Query(query, username)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	defer rows.Close()
	session_list := make([]models.Session, 0)
	for rows.Next() {
		session := models.Session{}
		err := rows.Scan(
			&session.ID,
			&session.LastTime,
			&session.LastIP,
			&session.LastLocation,
		)
		if err != nil {
			return nil, err
		}
		session_list = append(session_list, session)

	}
	return session_list, nil
}

func (dal *userDAL) CloseSession(username string, sid string) error {
	query := "DELETE FROM tbl_sessions WHERE id = ? AND username = ?"
	_, err := dal.DB.Exec(query, sid, username)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) ClearSessions(username string) error {
	query := "DELETE FROM tbl_sessions WHERE username = ?"
	_, err := dal.DB.Exec(query, username)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) GetSafety(username string) (models.Safety, error) {
	res := models.Safety{}
	var perms string
	query := "SELECT content_filter, blocked_location, timeline_permission, message_filter FROM tbl_safety WHERE username = ?"
	dal.DB.QueryRow(query, username).Scan(&res.Content, &res.Locations, &perms, &res.Message)

	query = "SELECT tbl_permissions_list.following, tbl_permissions_list.subscribed, tbl_permissions_list.tipped, tbl_permissions_list.like FROM tbl_permissions_list WHERE FIND_IN_SET(tbl_permissions_list.id, ?)"
	rows, err := dal.DB.Query(query, perms)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	permslist := make([]models.Permission, 0)
	for rows.Next() {
		perm_item := models.Permission{}
		err := rows.Scan(&perm_item.Following, &perm_item.Subscribed, &perm_item.Tipped, &perm_item.Like)
		if err != nil {
			return res, err
		}
		permslist = append(permslist, perm_item)
	}
	res.Permissions = permslist
	return res, nil
}

func (dal *userDAL) SaveSafety(username string, content int, location string, permids string, message int) error {
	query1 := "SELECT count(*) as count FROM tbl_safety WHERE username = ?"
	var count int
	dal.DB.QueryRow(query1, username).Scan(&count)
	if count == 0 {
		query := "INSERT INTO tbl_safety (username, content_filter, blocked_location, timeline_permission, message_filter) VALUES (?, ?, ?, ?, ?)"
		_, err := dal.DB.Exec(query, username, content, location, permids, message)
		if err != nil {
			fmt.Printf("Error Inserting tbl_safety: %v\n", err)
			return err
		}
		return nil
	} else {
		query := "UPDATE tbl_safety SET content_filter=?, blocked_location=?, timeline_permission=?, message_filter=? WHERE username =?"
		_, err := dal.DB.Exec(query, content, location, permids, message, username)
		if err != nil {
			fmt.Printf("Error updating tbl_safety: %v\n", err)
			return err
		}
		return nil
	}

}

func (dal *userDAL) SavePermissions(perms [][]interface{}) ([]string, error) {
	ids := make([]string, 0)
	for _, perm := range perms {
		permission_item := models.Permission{}
		permission_item.Following = 0
		permission_item.Like = 0
		permission_item.Subscribed = "0"
		permission_item.Tipped = 0
		for _, item := range perm {
			dataType := reflect.TypeOf(item).String()
			if dataType == "float64" {
				permission_item.Tipped = item.(float64)
			} else if item == "subscribed" {
				permission_item.Subscribed = "1"
			} else if item == "following" {
				permission_item.Following = 1
			} else if item == "followed" {
				permission_item.Like = 1
			}
		}
		query := "INSERT INTO tbl_permissions_list (following, subscribed, tipped, tbl_permissions_list.like) VALUES(?, ?, ?, ?)"
		res, err := dal.DB.Exec(query, permission_item.Following, permission_item.Subscribed, permission_item.Tipped, permission_item.Like)
		if err != nil {
			fmt.Printf("error1 %v\n", err)
			return nil, err
		}
		id, _ := res.LastInsertId()
		myid := strconv.FormatInt(id, 10)
		ids = append(ids, myid)
	}
	return ids, nil
}

func (dal *userDAL) ChangePassword(username string, newpassword string) error {
	_password, err := bcrypt.GenerateFromPassword([]byte(newpassword), BcryptCost)
	if err != nil {
		return err
	}
	_, err = dal.DB.Exec(`UPDATE users SET password_hash = ? WHERE name = ?;`, _password, username)
	return err
}

func (dal *userDAL) GetOtpSecret(username string, otptoken string) (string, error) {
	query := "SELECT otp_secret FROM users where name = ?"
	var qrcode string
	err := dal.DB.QueryRow(query, username).Scan(&qrcode)
	if err != nil {
		return "", err
	}

	return qrcode, nil
}

func (dal *userDAL) SetTwoFactorEnabled(username string, secret string) error {
	query := "UPDATE users SET otp_enabled = 1, otp_secret = ? WHERE name = ?"
	_, err := dal.DB.Exec(query, secret, username)
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) SetTwoFactorDisabled(username string) error {
	query := "UPDATE users SET otp_enabled = 0 WHERE name = ?"
	_, err := dal.DB.Exec(query, username)
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) RequestVCodeEmail(username string, vcode1 string) error {
	query := "INSERT INTO tbl_vemails (username, vcode1) VALUES (?, ?)"
	_, err := dal.DB.Exec(query, username, vcode1)
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) RequestVCodeEmail2(username string, vcode2 string) error {
	query := "UPDATE tbl_vemails SET vcode2 = ? WHERE username = ?"
	_, err := dal.DB.Exec(query, vcode2, username)
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) GetEmailByUsername(username string) (string, error) {
	query := "SELECT email FROM users WHERE name = ?"
	var email string
	err := dal.DB.QueryRow(query, username).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (dal *userDAL) CheckVCode1(username string, newemail string, vcode1 string) (string, error) {
	query := "SELECT vcode1 FROM tbl_vemails WHERE username=?"
	var _vcode1 string
	err := dal.DB.QueryRow(query, username).Scan(&_vcode1)
	if err != nil {
		return "failed", err
	}
	if vcode1 != _vcode1 {
		return "incorrect", nil
	}
	query = "UPDATE tbl_vemails SET change_email=? WHERE username=?"
	_, err = dal.DB.Exec(query, newemail, username)
	if err != nil {
		return "failed", err
	}
	return "success", nil
}

func (dal *userDAL) CheckVCode2(username string, vcode2 string, newemail string) (string, error) {
	query := "SELECT vcode2 FROM tbl_vemails WHERE username=?"
	var _vcode2 string
	err := dal.DB.QueryRow(query, username).Scan(&_vcode2)
	if err != nil {
		return "failed", err
	}
	if vcode2 != _vcode2 {
		return "incorrect", nil
	}
	query = "DELETE FROM tbl_vemails WHERE username=?"
	_, err = dal.DB.Exec(query, username)
	if err != nil {
		return "failed", err
	}
	query = "UPDATE users SET email = ? WHERE name=?"
	_, err = dal.DB.Exec(query, newemail, username)
	if err != nil {
		return "failed", err
	}
	return "success", nil
}

func (dal *userDAL) CheckPassword(username string, password string) error {
	query := "SELECT password_hash FROM users WHERE name = ?"
	var password_hash string
	_ = dal.DB.QueryRow(query, username).Scan(&password_hash)
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) ChangeDisplayname(userid int, displayname string) error {
	query := "UPDATE tbl_userdata SET displayname = ? WHERE userid = ?"
	_, err := dal.DB.Exec(query, displayname, userid)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	return nil
}
func (dal *userDAL) ChangeUsername(username string, password string, newusername string) (string, error) {
	query := "SELECT password_hash FROM users WHERE name = ?"
	var password_hash string
	_ = dal.DB.QueryRow(query, username).Scan(&password_hash)
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return "password", nil
	}
	query = "SELECT count(*) as count FROM users WHERE name = ?"
	var count string
	_ = dal.DB.QueryRow(query, newusername).Scan(&count)
	if count == "0" {
		query = "UPDATE users SET name = ? WHERE name = ?"
		_, err = dal.DB.Exec(query, newusername, username)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return "failed", nil
		}
		return "success", nil
	} else {
		return "exist", nil
	}
}

// Application

func (dal *userDAL) AcceptApplication(username string) error {
	query := "UPDATE tbl_application SET status = 1 WHERE username = ?"
	_, err := dal.DB.Exec(query, username)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	query = "UPDATE users SET active = 2 WHERE name = ?"
	_, err = dal.DB.Exec(query, username)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) RejectApplication(username string) error {
	query := "UPDATE tbl_application SET status = 2 WHERE username = ?"
	_, err := dal.DB.Exec(query, username)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	return nil
}

func (dal *userDAL) GetApplicationByUsername(username string) (models.Application, error) {
	query := `SELECT
	username,
	firstname,
	lastname,
	birthday,
	address,
	city,
	country,
	state,
	zipcode,
	twitter,
	instagram,
	website,
	document_type,
	explicit_content,
	file_front,
	file_back,
	file_handwritten,
	file_video
  FROM
	tbl_application
  WHERE username = ?`
	application := models.Application{}
	err := dal.DB.QueryRow(query, username).Scan(
		&application.Username,
		&application.Firstname,
		&application.Lastname,
		&application.Birthday,
		&application.Address,
		&application.City,
		&application.Country,
		&application.State,
		&application.Zipcode,
		&application.Twitter,
		&application.Instagram,
		&application.Website,
		&application.Documenttype,
		&application.ExplicitContent,
		&application.FileFront,
		&application.FileBack,
		&application.FileHandwritten,
		&application.FileVideo,
	)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
		return application, err
	}
	return application, nil
}

func (dal *userDAL) GetWaitingApplications() ([]models.ApplicationUser, error) {
	query := `SELECT
	A.displayname,
	tbl_application.username,
	A.avatar
  FROM
	tbl_application
	LEFT JOIN
	  (SELECT
		tbl_userdata.displayname,
		users.name,
		users.avatar
	  FROM
		users
		LEFT JOIN tbl_userdata
		  ON users.id = tbl_userdata.userid) AS A
	  ON A.name = tbl_application.username
	  WHERE tbl_application.status = 0`
	rows, err := dal.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]models.ApplicationUser, 0)
	for rows.Next() {
		user := models.ApplicationUser{}
		err := rows.Scan(&user.DisplayName, &user.Username, &user.Avatar)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (dal *userDAL) SubmitApplication(application models.Application) (string, error) {
	query := "SELECT count(*) as count from tbl_application where username = ?"
	var cnt int
	err := dal.DB.QueryRow(query, application.Username).Scan(&cnt)
	if err != nil {
		return "failed", err
	}
	if cnt == 1 {
		return "exist", nil
	}
	query = "INSERT INTO tbl_application (username, firstname, lastname, birthday, address, city, country, state, zipcode, twitter, instagram, website, document_type, explicit_content, file_front, file_back, file_handwritten, file_video) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = dal.DB.Exec(query, application.Username, application.Firstname, application.Lastname, application.Birthday, application.Address, application.City, application.Country, application.State, application.Zipcode, application.Twitter, application.Instagram, application.Website, application.Documenttype, application.ExplicitContent, application.FileFront, application.FileBack, application.FileHandwritten, application.FileVideo)
	if err != nil {
		return "failed", err
	}
	return "success", nil
}

func (dal *userDAL) UpdateProfile1(profile models.Profile) error {
	query := "UPDATE tbl_userdata SET displayname=?, aboutme=?, link_twitter=?, link_instagram=?, link_tiktok=?, location=? WHERE userid=?"
	_, err := dal.DB.Exec(query, profile.DisplayName, profile.Aboutme, profile.LinkTwitter, profile.LinkInstagram, profile.LinkTiktok, profile.Location, profile.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) UpdateBanner(userID int, banner string) error {
	query := "UPDATE tbl_userdata SET userBanner = ? WHERE userid=?"
	_, err := dal.DB.Exec(query, banner, userID)
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) UpdateAvatar(userID int, avatar string) error {
	query := "UPDATE users SET avatar = ? WHERE id= ?"
	_, err := dal.DB.Exec(query, avatar, userID)
	if err != nil {
		return err
	}
	return nil
}

func (dal *userDAL) IsActivated(userID int) (bool, error) {
	var user models.User
	err := dal.DB.QueryRow("SELECT active FROM users WHERE id =?", userID).Scan(&user.Active)
	if err != nil {
		return false, err
	}
	return user.Active != 0, nil
}

func (dal *userDAL) GetUserDataByName(name string) (models.CompleteUser, error) {
	query := `SELECT
	users.id AS userid,
	tbl_userdata.displayname,
	users.name,
	users.status,
	users.avatar,
	tbl_userdata.userBanner,
	tbl_userdata.location,
	users.active,
	tbl_userdata.likes_cnt,
	tbl_userdata.followers_cnt,
	tbl_userdata.aboutme,
	tbl_userdata.link_twitter,
	tbl_userdata.link_instagram,
	tbl_userdata.link_tiktok
  FROM
	users
	LEFT JOIN tbl_userdata
	  ON users.id = tbl_userdata.userid
  WHERE users.name = ? AND ISNULL(users.deleted_at)`

	user := models.CompleteUser{}
	active := models.ActiveUser{}
	err := dal.DB.QueryRow(query, name).Scan(
		&user.UserID,
		&user.DisplayName,
		&user.UserName,
		&user.Availability,
		&user.AvatarUrl,
		&user.BannerUrl,
		&user.Location,
		&active.Active,
		&user.Likes,
		&user.Followers,
		&user.AboutMe,
		&user.Twitter,
		&user.Instagram,
		&user.Tiktok,
	)
	if err != nil {
		return user, err
	}
	if active.Active == 0 {
		admin := "false"
		verified := "false"
		user.Admin = &admin
		user.Verified = &verified
	} else if active.Active == 1 {
		admin := "false"
		verified := "false"
		user.Admin = &admin
		user.Verified = &verified
	} else if active.Active == 2 {
		admin := "false"
		verified := "true"
		user.Admin = &admin
		user.Verified = &verified
	} else if active.Active == 3 {
		admin := "true"
		verified := "true"
		user.Admin = &admin
		user.Verified = &verified
	}
	return user, nil

}

func (dal *userDAL) GetByName(name string) (models.User, error) {
	return dal.GetOne("name", name)
}
func (dal *userDAL) GetByEmail(email string) (models.User, error) {
	return dal.GetOne("email", email)
}

func (dal *userDAL) GetByID(userID int) (models.User, error) {
	return dal.GetOne("id", userID)
}

func (dal *userDAL) UpdateActivationCode(userID int, code string) error {

	query := "UPDATE users SET activation_code =?, updated_at =? WHERE id =?"

	_, err := dal.DB.Query(query, code, time.Now(), userID)
	return err
}

func (dal *userDAL) ConfirmActivationCode(userID int, code string) (string, error) {
	query := "SELECT id FROM users WHERE id =? AND activation_code =?"

	user := models.User{}
	err := dal.DB.QueryRow(query, userID, code).Scan(
		&user.ID,
	)
	if err != nil {
		return "failed", err
	} else {
		query := "UPDATE users SET active = 1, updated_at =? WHERE id =?"
		_, err := dal.DB.Query(query, time.Now(), userID)
		if err != nil {
			return "failed", err
		}
		return "success", nil
	}

}

func (dal *userDAL) GetByActivationCode(code string) (models.User, error) {
	return dal.GetOne("activation_code", code)
}

func (dal *userDAL) GetOne(whereCondition string, args any) (models.User, error) {
	user := models.User{}
	query := `SELECT
	id, email, COALESCE(name, '') as name, avatar, password_hash, active, created_at, otp_enabled, otp_secret, status
	FROM users
	WHERE %s = ? AND ISNULL(users.deleted_at);`

	err := dal.DB.QueryRow(fmt.Sprintf(query, whereCondition), args).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Avatar,
		&user.PasswordHash,
		&user.Active,
		&user.CreatedAt,
		&user.OtpEnabled,
		&user.OtpSecret,
		&user.Status,
	)

	return user, err
}

func (dal *userDAL) InsertGetID(user models.User) error {
	insertQuery := `
	INSERT INTO users (name, email, password_hash, activation_code, otp_secret)
	VALUES (?, ?, ?, ?, '') RETURNING id;
	`
	var result int
	err := dal.DB.QueryRow(insertQuery, user.Name, user.Email, user.PasswordHash, user.ActivationCode).Scan(&result)

	if err != nil {
		return err
	}
	fmt.Printf("%v\n", result)
	insertID := result
	// insertID, err := result.LastInsertId()
	// if err != nil {
	// 	return err
	// }
	insertQuery = "INSERT INTO tbl_userdata (userid, displayname) VALUES (?, ?)"
	_, err = dal.DB.Exec(insertQuery, insertID, user.Name)
	if err != nil {
		return err
	}
	return nil
}

// Return value: User ID, RowsAffect, Error
func (dal *userDAL) Activate(userID int) error {
	_, err := dal.DB.Exec(`UPDATE users SET active = 1 WHERE id = ?;`, userID)
	return err
}

func (dal *userDAL) GetDisplayNameByID(userID int) (string, error) {
	var displayname string
	query := "SELECT displayname FROM tbl_userdata WHERE userid = ?"
	err := dal.DB.QueryRow(query, userID).Scan(&displayname)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	return displayname, nil
}

func (dal *userDAL) GetAll(filter types.UserFilter) ([]models.User, int, error) {
	queryTotal := `
	SELECT COUNT(*)
	FROM users
	WHERE name LIKE "%` + filter.Query + `%"
		OR email LIKE "%` + filter.Query + `%"
	;
	`

	total := 0
	if err := dal.DB.QueryRow(queryTotal).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
	SELECT id, name, email, avatar, active, created_at
	FROM users
	WHERE name LIKE "%` + filter.Query + `%"
		OR email LIKE "%` + filter.Query + `%"
	ORDER BY name
	LIMIT ? OFFSET ?;
	`

	rows, err := dal.DB.Query(query, filter.PageSize, filter.PageSize*(filter.Page-1))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Avatar,
			&user.ActivationCode,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Status,
		)
		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}
	return users, total, nil
}

func (dal *userDAL) UpdatePassword(userID int, password string) error {
	_, err := dal.DB.Exec(`UPDATE users SET password_hash = ? WHERE id = ?;`, password, userID)
	return err
}

func (dal *userDAL) Block(userID int, status types.UserStatusType) error {
	_, err := dal.DB.Exec(`UPDATE users SET active = ? WHERE id = ?;`, status, userID)
	return err
}

func (dal *userDAL) UpdateProfile(userID int, user models.User) error {
	_, err := dal.DB.Exec(`UPDATE users SET name = ?, avatar = ? WHERE id = ?;`, user.Name, user.Avatar, userID)
	return err
}

func (dal *userDAL) SetPasswordResetCode(userID int, code string) error {
	_, err := dal.DB.Exec(
		`REPLACE INTO password_reset (user_id, code, expiry) VALUES (?, ?, ?);`,
		userID,
		code,
		time.Now().Add(time.Hour*24),
	)
	return err
}

func (dal *userDAL) GetUserIDFromPasswordResetCode(code string) (id int, expiry time.Time, err error) {
	err = dal.DB.QueryRow(`SELECT user_id, expiry FROM password_reset WHERE code = ?;`, code).Scan(&id, &expiry)
	return
}

func (dal *userDAL) DeletePasswordResetCode(code string) error {
	_, err := dal.DB.Exec(`DELETE FROM password_reset WHERE code = ?;`, code)
	return err
}
