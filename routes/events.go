package routes

import (
	"log"
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(ctx *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve events"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func getEvent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := models.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func createEvent(ctx *gin.Context) {

	var event models.Event
	err := ctx.ShouldBindJSON(&event)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.UserID = ctx.GetInt64("userId")
	err = event.Save()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save event"})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

func updateEvent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userId := ctx.GetInt64("userId")
	event, err := models.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.UserID != 0 && event.UserID != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	var updateEvent models.Event
	err = ctx.ShouldBindJSON(&updateEvent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request"})
		return
	}

	updateEvent.ID = id
	updateEvent.UserID = userId
	err = updateEvent.Update()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update event"})
		return
	}

	ctx.JSON(http.StatusOK, updateEvent)
}

func deleteEvent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	userId := ctx.GetInt64("userId")
	event, err := models.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.UserID != 0 && event.UserID != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	err = event.Delete()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete event"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
