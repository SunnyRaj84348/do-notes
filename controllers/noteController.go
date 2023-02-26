package controllers

import (
	"net/http"
	"strconv"

	"github.com/SunnyRaj84348/do-notes/models"
	"github.com/gin-gonic/gin"
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
	notes, err := models.InsertNotes(userid.(uint32), note.NoteTitle, note.NoteBody)
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
	notes, err := models.GetNotes(userid.(uint32))
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

	val := ctx.Param("id")

	// Convert string to int
	id, err := strconv.Atoi(val)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Update specified note
	notes, err := models.UpdateNotes(userid.(uint32), uint32(id), note.NoteTitle, note.NoteBody)
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
	val := ctx.Param("id")

	// Convert string to int
	id, err := strconv.Atoi(val)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Delete specified note from database
	err = models.DeleteNotes(userid.(uint32), uint32(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
