package main

import (
	"math/rand"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/paujim/cognitoserver/server/pkg/controllers"
	"github.com/paujim/cognitoserver/server/pkg/entities"
	"github.com/paujim/cognitoserver/server/pkg/services"
	log "github.com/sirupsen/logrus"
)

const (
	region = "us-west-2"
)

var (
	cognito    entities.UserTokenHandler
	userPoolID string
)

func randomString(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	paramStore := services.NewParameterStore(ssm.New(sess))
	userPoolID, _ := paramStore.Get("/pj/userpool/id")
	appClientID, _ := paramStore.Get("/pj/userpool/appclient/id")
	cognito = services.NewCognitoHandler(appClientID, userPoolID, cognitoidentityprovider.New(sess))

}

func main() {

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile(path.Join("..", "client", "build"), true)))
	router.Use(controllers.CorsMiddleware())

	api := router.Group("/api")

	// No auth
	controllers.RegisterPing(api)

	a := controllers.NewAuth(region, userPoolID, cognito)
	a.RegisterAuthRoutes(api)
	api.Use(a.AuthMiddleware())
	controllers.NewUser(cognito).RegisterUserRoutes(api.Group("/user"))

	// Start and run the server
	router.Run(":5000")

}
