package models

import (
	"database/sql"
	"fanclub/internal/types"
	"time"
)

type User struct {
	ID             int                  `json:"id"`
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	Avatar         string               `json:"avatar"`
	PasswordHash   string               `json:"password_hash"`
	Active         types.UserStatusType `json:"active"`
	ActivationCode string               `json:"activation_code"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Status         string               `json:"status"`
	OtpEnabled     int                  `json:"otp_enabled"`
	OtpSecret      string               `json:"otp_secret"`
}

type CompleteUser struct {
	UserID       int     `json:"userid"`
	DisplayName  *string `json:"displayname"`
	UserName     *string `json:"name"`
	Availability *string `json:"status"`
	AvatarUrl    *string `json:"avatar"`
	BannerUrl    *string `json:"userBanner"`
	Location     *string `json:"location"`
	Verified     *string `json:"verified"`
	Admin        *string `json:"admin"`
	Likes        *int    `json:"likes"`
	Followers    *int    `json:"followers"`
	AboutMe      *string `json:"aboutme"`
	Twitter      *string `json:"link_twitter"`
	Instagram    *string `json:"link_instagram"`
	Tiktok       *string `json:"link_tiktok"`
}

type ApplicationUser struct {
	DisplayName string `json:"displayname"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
}

type ActiveUser struct {
	Active int `json:"active"`
}

type Permission struct {
	Following  int     `json:"following"`
	Subscribed string  `json:"subscribed"`
	Tipped     float64 `json:"tipped"`
	Like       int     `json:"like"`
}

type Safety struct {
	Content     int          `json:"content_filter" `
	Locations   string       `json:"blocked_location"`
	Message     int          `json:"message_filter"`
	Permissions []Permission `json:"permissions"`
}

type Session struct {
	ID           int    `json:"id"`
	LastTime     string `json:"last_used_time"`
	LastIP       string `json:"last_used_ip"`
	LastLocation string `json:"last_used_location"`
}

type Creator struct {
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
	Avatar      string `json:"avatar"`
	Banner      string `json:"banner"`
}

type TierCards struct {
	Creator    Creator `json:"creator"`
	Tiername   string  `json:"tiername"`
	Status     string  `json:"astatus"`
	Expiration string  `json:"expiration"`
	Baseprice  string  `json:"baseprice"`
}

type TiersCount struct {
	Expired string `json:"expired"`
	Active  string `json:"active"`
}

type NotificationSettings struct {
	UserID     int `json:"userid"`
	IsPush     int `json:"is_push"`
	IsMessage  int `json:"is_message"`
	IsReply    int `json:"is_reply"`
	IsPostlike int `json:"is_postlike"`
	IsFollower int `json:"is_follower"`
}

type CustomList struct {
	ID       int            `json:"id"`
	Listname string         `json:"listname"`
	UserList sql.NullString `json:"userlist"`
	Entries  int            `json:"entries"`
}

type CustomUserList struct {
	Listname string       `json:"listname"`
	Userlist []SimpleUser `json:"userlist"`
}

type DefaultEntries struct {
	Block int `json:"block"`
	Mute  int `json:"mute"`
	Vip   int `json:"vip"`
}

type DefaultEntriesList struct {
	Blocklist sql.NullString `json:"blocklist"`
	Mutelist  sql.NullString `json:"mutelist"`
	Viplist   sql.NullString `json:"vipllist"`
}
