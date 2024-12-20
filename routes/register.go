package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(ctx *gin.Context) {
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

	err = event.Register(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register for event"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Registered for event"})
}

func unregisterFromEvent(ctx *gin.Context) {
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
	err = event.Unregister(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not unregister from event"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Unregistered from event"})
}
