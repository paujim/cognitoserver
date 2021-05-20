package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paujim/cognitoserver/server/pkg/entities"
)

type user struct {
	service entities.UserHandler
}

func NewUser(service entities.UserHandler) *user {
	return &user{
		service: service,
	}
}

func (u *user) RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/register", u.registerUser)
	router.GET("/list", u.listUsers)
}

func (u *user) registerUser(c *gin.Context) {
	var request entities.RegistrationRequest
	c.BindJSON(&request)
	sub, err := u.service.RegisterUser(request.Username, request.Password)
	if err == nil {
		c.JSON(http.StatusAccepted, gin.H{
			"status": "registered",
			"sub":    sub,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}

func (u *user) listUsers(c *gin.Context) {
	users, err := u.service.ListUsers()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"users": users})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}
