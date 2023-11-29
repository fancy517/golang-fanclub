package models

import "time"

type Post struct {
	ID            int       `json:"id"`
	UserID        int       `json:"userid"`
	Name          string    `json:"name"`
	ParentID      int       `json:"parentid"`
	Message       string    `json:"message"`
	IsPinned      int       `json:"is_pinned"`
	ReplyRole     string    `json:"reply_role"`
	Files         string    `json:"files"`
	PublishDate   string    `json:"publish_date"`
	DisappearDate string    `json:"disappear_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

type PostCreator struct {
	DisplayName  *string `json:"displayName"`
	UserName     *string `json:"userName"`
	Availability *string `json:"availability"`
	AvatarUrl    *string `json:"avatarUrl"`
	BannerUrl    *string `json:"bannerUrl"`
	Verified     *string `json:"verified"`
}

type Medias struct {
	MediaId   string  `json:"id"`
	Type      string  `json:"type"`
	Sensetive string  `json:"sensetive"`
	Loked     string  `json:"loked"`
	SourceId  string  `json:"sourceid"`
	Timestamp string  `json:"timestamp"`
	Tags      *string `json:"tags"`
}

type Attachment struct {
	Medias []Medias
}

type Postdata struct {
	Creator       PostCreator
	Description   *string `json:"description"`
	Tags          *string `json:"tags"`
	Attachment    Attachment
	LikesCnt      *string `json:"likes_cnt"`
	ReplyRole     *string `json:"reply_role"`
	PublishDate   *string `json:"publish_date"`
	DisappearDate *string `json:"disappear_date"`
	Files         string  `json:"files"`
	CommentsCnt   *string `json:"comments_cnt"`
	IsLiked       *string `json:"is_liked"`
	ID            string  `json:"id"`
	IsPinned      *string `json:"is_pinned"`
}

type Profile struct {
	UserID        int    `json:"userid"`
	DisplayName   string `json:"displayname"`
	Banner        string `json:"userbanner"`
	Aboutme       string `json:"aboutme"`
	LinkTwitter   string `json:"link_twitter"`
	LinkInstagram string `json:"link_instagram"`
	LinkTiktok    string `json:"link_tiktok"`
	Location      string `json:"location"`
}
