package main

import (
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	userPoolID    string
	appClientID   string
	paramStore    *ParamStore
	cognito       *CognitoParam
	logger        *AWSLogger
	region        = "us-west-2"
	logGroupName  = "CognitoServer"
	logStreamName = time.Now().UTC().Format("2006010203-") + randomString(10)
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
}

func main() {

	client := cloudwatchlogs.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion(region))
	logger := NewAWSLogger("INFO", logGroupName, logStreamName, client, log.New(os.Stdout, "", 0))

	paramStore = NewParamStore(region, logger)
	userPoolID, _ = paramStore.Get("/pj/userpool/id")
	appClientID, _ = paramStore.Get("/pj/userpool/appclient/id")
	cognito = NewCognitoParam(region, appClientID, userPoolID, cognitoidentityprovider.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion(region)), logger)

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
