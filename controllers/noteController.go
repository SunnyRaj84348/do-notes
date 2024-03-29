package controllers

import (
	"net/http"

	"github.com/SunnyRaj84348/do-notes/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InsertNote(ctx *gin.Context) {
	userid, _ := ctx.Get("userid")
	note := models.Note{}

	err := ctx.BindJSON(&note)
	if err != nil {
		return
	}

	// Insert note into database
	notes, err := models.InsertNotes(userid.(string), note.NoteTitle, note.NoteBody)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Write inserted note to response body
	ctx.JSON(http.StatusOK, notes)
}

func GetNotes(ctx *gin.Context) {
	userid, _ := ctx.Get("userid")

	// Retrieve notes from database
	notes, err := models.GetNotes(userid.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Write notes to response body
	ctx.JSON(http.StatusOK, notes)
}

func UpdateNote(ctx *gin.Context) {
	userid, _ := ctx.Get("userid")
	note := models.Note{}

	err := ctx.BindJSON(&note)
	if err != nil {
		return
	}

	id := ctx.Param("id")

	// Check for valid uuid syntax
	_, err = uuid.Parse(id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Update specified note
	notes, err := models.UpdateNotes(userid.(string), id, note.NoteTitle, note.NoteBody)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	// Write updated note to response body
	ctx.JSON(http.StatusOK, notes)
}

func DeleteNote(ctx *gin.Context) {
	userid, _ := ctx.Get("userid")
	id := ctx.Param("id")

	// Check for valid uuid syntax
	_, err := uuid.Parse(id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Delete specified note from database
	err = models.DeleteNotes(userid.(string), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
