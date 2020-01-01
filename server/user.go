package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//RegistrationRequest ...
type RegistrationRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

//RegisterUserRoutes ...
func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/register", registerUser)
	router.GET("/list", listUsers)
}

func registerUser(c *gin.Context) {
	var request RegistrationRequest
	c.BindJSON(&request)
	sub, err := cognito.RegisterUser(request.Username, request.Password)
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

func listUsers(c *gin.Context) {
	users, err := cognito.ListUsers()
	if err == nil {
		c.JSON(http.StatusAccepted, users)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}
