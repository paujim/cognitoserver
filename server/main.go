package main

import (
	"path"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	userPoolID  string
	appClientID string
	paramStore  *ParamStore
	cognito     *CognitoParam
	region      = "us-west-2"
	logger      = NewAPILogger("[INFO] ")
)

func init() {
	paramStore = NewParamStore(region, logger)
	userPoolID, _ = paramStore.Get("/pj/userpool/id")
	appClientID, _ = paramStore.Get("/pj/userpool/appclient/id")
	cognito = NewCognitoParam(region, appClientID, userPoolID, logger)
}

func main() {

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile(path.Join("..", "client", "build"), true)))
	router.Use(CorsMiddleware())

	api := router.Group("/api")

	// No auth
	RegisterPing(api)
	RegisterAuthRoutes(api)

	api.Use(AuthMiddleware(region, userPoolID))
	RegisterUserRoutes(api.Group("/user"))

	// Start and run the server
	router.Run(":5000")
}
