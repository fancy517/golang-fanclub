package user

import (
	"fanclub/internal/models"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) GetProfile(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")

	fmt.Printf("username: %v\n", username)

	user, err := h.dal.User.GetUserDataByName(username)
	if err != nil {
		h.SendSuccess(c, "failed")
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("user: %v\n", user)
	h.SendSuccess(c, user)
}

func (h *HandlerImpl) GetOnePost(c *gin.Context) {
	postid := c.Request.URL.Query().Get("postid")
	username := c.Request.URL.Query().Get("loginuser")

	postdata, err := h.dal.Postlist.GetOnePost(postid, username)
	if err != nil {
		h.SendSuccess(c, nil)
	} else {
		h.SendSuccess(c, postdata)
	}
}

func (h *HandlerImpl) GetPostData(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	loginuser := c.Request.URL.Query().Get("loginuser")
	parentid := c.Request.URL.Query().Get("parentid")
	selecttag := c.Request.URL.Query().Get("selecttag")
	fmt.Printf("username: %v\n", username)
	fmt.Printf("loginuser: %v\n", loginuser)

	postdata, err := h.dal.Postlist.GetPostData(username, "%"+loginuser+"%", parentid, selecttag)
	if err != nil {
		h.SendSuccess(c, nil)
		fmt.Printf("error: %v\n", err)
		return
	}
	h.SendSuccess(c, postdata)
}

func (h *HandlerImpl) GetTags(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	fmt.Printf("%v\n", username)
	tagslist, err := h.dal.Postlist.GetTags(username)
	if err != nil {
		h.SendSuccess(c, "")
		fmt.Printf("error: %v\n", err)
		return
	}
	h.SendSuccess(c, tagslist)
}

func (h *HandlerImpl) GetMediaCounts(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	mediacnts, err := h.dal.Media.GetMediaCounts(username)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		h.SendSuccess(c, "")
		return
	}
	h.SendSuccess(c, mediacnts)
}

func (h *HandlerImpl) SaveProfile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.SendSuccess(c, "failed")
	}
	files := form.File["files"]
	fmt.Printf("%v\n", files)

	var profile models.Profile
	var user models.User
	user, err = h.dal.User.GetByName(form.Value["username"][0])
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	profile.UserID = user.ID
	profile.DisplayName = form.Value["displayname"][0]
	profile.Aboutme = form.Value["aboutme"][0]
	profile.LinkTwitter = form.Value["twitter"][0]
	profile.LinkInstagram = form.Value["instagram"][0]
	profile.LinkTiktok = form.Value["tiktok"][0]
	profile.Location = form.Value["location"][0]

	for idx, file := range files {
		// Save each file to the uploads directory
		ext := filepath.Ext(file.Filename)
		currentTime := time.Now()
		timestamp := currentTime.Format("20060102_150405")
		sanitizedFilename := sanitizeFilename(file.Filename)
		newFilename := timestamp + "_" + sanitizedFilename + ext

		err = c.SaveUploadedFile(file, "uploads/"+newFilename)
		if err != nil {
			fmt.Printf("file error: %v\n", err)
			h.SendSuccess(c, "failed")
			return
		}
		fmt.Printf("%v\n", newFilename)
		if idx == 0 && len(form.Value["banner"]) == 1 {
			profile.Banner = newFilename
			err := h.dal.User.UpdateBanner(user.ID, newFilename)
			if err != nil {
				fmt.Printf("File to update banner: %v\n", err)
				h.SendSuccess(c, "failed to update banner")
				return
			}

		} else if idx == 1 || (idx == 0 && form.Value["avatar"][0] == "true") {
			err := h.dal.User.UpdateAvatar(user.ID, newFilename)
			if err != nil {
				fmt.Printf("error1: %v\n", err)
				h.SendSuccess(c, "Failed to update avatar")
				return
			}
		}
	}
	fmt.Printf("%v\n", form.Value)
	err = h.dal.User.UpdateProfile1(profile)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		h.SendSuccess(c, "Failed to update profile")
		return
	}
	h.SendSuccess(c, "success")
}
