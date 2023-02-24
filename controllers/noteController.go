package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SunnyRaj84348/do-notes/models"
	"github.com/gin-gonic/gin"
)

func InsertNote(ctx *gin.Context) {
	userid, _ := ctx.Get("userid")
	note := models.Note{}

	err := ctx.BindJSON(&note)
	if err != nil {
		return
	}

	// Insert note into database
	noteID, err := models.InsertNotes(userid.(int), note.NoteTitle, note.NoteBody)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	notes := models.Notes{
		NoteID:    noteID,
		NoteTitle: note.NoteTitle,
		NoteBody:  note.NoteBody,
	}

	// Write inserted note to response body
	ctx.JSON(http.StatusOK, notes)
}

func GetNotes(ctx *gin.Context) {
	userid, _ := ctx.Get("userid")

	// Retrieve notes from database
	rows, err := models.GetNotes(userid.(int))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer rows.Close()

	notes := []models.Notes{}

	for rows.Next() {
		note := models.Notes{}

		err := rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		notes = append(notes, note)
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
	err = models.UpdateNotes(userid.(int), id, note.NoteTitle, note.NoteBody)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	notes := models.Notes{
		NoteID:    id,
		NoteTitle: note.NoteTitle,
		NoteBody:  note.NoteBody,
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
	err = models.DeleteNotes(userid.(int), id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}
}
