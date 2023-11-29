package user

import (
	"fanclub/internal/models"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Function to sanitize the filename
func sanitizeFilename(filename string) string {
	// Replace spaces with underscores
	sanitized := strings.ReplaceAll(filename, " ", "_")

	// Remove special characters using regex or other sanitization techniques

	return sanitized
}

func (h *HandlerImpl) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	files := form.File["files"]
	for _, file := range files {
		// Save each file to the uploads directory
		ext := filepath.Ext(file.Filename)
		currentTime := time.Now()
		timestamp := currentTime.Format("20060102_150405")
		sanitizedFilename := sanitizeFilename(file.Filename)
		newFilename := timestamp + "_" + sanitizedFilename + ext

		err = c.SaveUploadedFile(file, "uploads/"+newFilename)
		if err != nil {
			h.SendSuccess(c, "failed")
			return
		}
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) RemovePost(c *gin.Context) {

	var req RemovePostBindingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding")
		return
	}

	postID := req.PostID
	userID := req.UserID
	err := h.dal.Postlist.RemovePost(postID, userID)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) CreatePost(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	files := form.File["files"]
	tags := form.Value["tags"]
	fileids := form.Value["fileids"]
	fmt.Printf("fileids-----------------%v\n", fileids)
	// fileids := []string{}
	for idx, file := range files {
		// Save each file to the uploads directory
		ext := filepath.Ext(file.Filename)
		currentTime := time.Now()
		timestamp := currentTime.Format("20060102_150405")
		sanitizedFilename := sanitizeFilename(file.Filename)
		newFilename := timestamp + "_" + sanitizedFilename + ext

		err = c.SaveUploadedFile(file, "uploads/"+newFilename)
		if err != nil {
			h.SendSuccess(c, "failed")
			return
		}
		file_id, err := h.dal.Media.CreateMedia(newFilename, tags[idx], form.Value["username"][0])
		if err != nil {
			h.SendSuccess(c, "failed")
			return
		}
		fileids = append(fileids, file_id)
	}
	var post models.Post
	var user models.User
	user, err = h.dal.User.GetByName(form.Value["username"][0])
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	fmt.Printf("%v\n", form.Value["publish_date"][0])

	post.UserID = user.ID
	post.Name = user.Name
	post.ParentID, err = strconv.Atoi(form.Value["parentid"][0])
	fmt.Printf("parentid: %v\n", post.ParentID)
	if err != nil {
		fmt.Printf("Int converting :%v\n", err)
		h.SendSuccess(c, "failed")
	}
	post.Message = form.Value["message"][0]
	post.IsPinned, _ = strconv.Atoi(form.Value["is_pinned"][0])
	post.ReplyRole = form.Value["reply_role"][0]
	post.PublishDate = form.Value["publish_date"][0]
	post.DisappearDate = form.Value["disappear_date"][0]
	post.Files = strings.Join(fileids, ",")

	if form.Value["postid"][0] == "-1" || form.Value["parentid"][0] != "0" {
		err1 := h.dal.Postlist.CreatePost(post)
		if err1 != nil {
			h.SendSuccess(c, "failed")
			return
		}
		if post.ParentID != 0 {
			err1 = h.dal.Postlist.IncreaseCommentCount(post.ParentID)
			if err1 != nil {
				h.SendSuccess(c, "failed")
				return
			}
		}
		// fmt.Printf("parentid: %v\n", post.ParentID)
		h.SendSuccess(c, "success")
	} else {
		err1 := h.dal.Postlist.UpdatePost(post, form.Value["postid"][0])
		if err1 != nil {
			h.SendSuccess(c, "failed")
			return
		}
		h.SendSuccess(c, "updated")
	}
}
