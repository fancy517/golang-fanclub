package user

import (
	"fanclub/internal/models"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) SubmitApplication(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.SendSuccess(c, "failed")
	}
	files := form.File["files"]
	var application models.Application
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
		if idx == 0 {
			application.FileFront = newFilename
		} else if idx == 1 {
			application.FileBack = newFilename
		} else if idx == 2 {
			application.FileHandwritten = newFilename
		} else if idx == 3 {
			application.FileVideo = newFilename
		}
	}
	application.Username = form.Value["username"][0]
	application.Firstname = form.Value["firstname"][0]
	application.Lastname = form.Value["lastname"][0]
	application.Birthday = form.Value["birthday"][0]
	application.Address = form.Value["address"][0]
	application.City = form.Value["city"][0]
	application.Country = form.Value["country"][0]
	application.State = form.Value["state"][0]
	application.Zipcode = form.Value["zipcode"][0]
	application.Twitter = form.Value["twitter"][0]
	application.Instagram = form.Value["instagram"][0]
	application.Website = form.Value["website"][0]
	application.Documenttype = form.Value["document_type"][0]
	application.ExplicitContent = form.Value["explicit_content"][0]
	status, err := h.dal.User.SubmitApplication(application)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	if status == "exist" {
		h.SendSuccess(c, "exist")
		return
	} else if status == "success" {
		h.SendSuccess(c, "success")
		return
	} else {
		h.SendSuccess(c, "failed")
		return
	}
}
