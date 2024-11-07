package main

import (
	"context"
	"log"
	"musiclib/connect"
	"musiclib/controllers"
	docs "musiclib/docs"
	auth "musiclib/jwt-authenticate"
	"musiclib/models"
	implements "musiclib/services/implement"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	trackController *controllers.TrackController
	albumController *controllers.AlbumController
	userController  *controllers.UserController
	ctx             context.Context
	mongoClient     *mongo.Client
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	if err := connect.Connect(); err != nil {
		log.Fatal("err connect db", err)
	}
	ctx := context.TODO()
	trackCollection := connect.Ng.Database.Collection("tracks")
	trackService := implements.NewTrackService(trackCollection, ctx)
	trackController = controllers.NewTrackController(trackService)

	userCollection := connect.Ng.Database.Collection("users")
	userService := implements.NewUserService(userCollection, ctx)
	userController = controllers.NewUserController(userService)

	albumCollection := connect.Ng.Database.Collection("albums")
	albumService := implements.NewAlbumService(albumCollection, trackCollection, ctx)
	albumController = controllers.NewAlbumController(albumService)
}

func returnUser(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	user, _ := c.Get("userId")
	c.JSON(200, gin.H{
		"userId": user.(*models.User).UserId,
	})
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
func main() {
	Init()
	authMiddleware := auth.NewJWTAuthMiddleware(userController)
	defer mongoClient.Disconnect(ctx)
	docs.SwaggerInfo.BasePath = "/v1"
	r := gin.Default()
	group := os.Getenv("SERVER_GROUP")
	basepath := r.Group(group)
	// basepath.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	basepath.POST("/login", authMiddleware.LoginHandler)

	basepath.GET("/refresh_token", authMiddleware.RefreshHandler)

	// Apply middleware only to the /currentUser route
	basepath.GET("/currentUser", authMiddleware.MiddlewareFunc(), returnUser)

	// Apply middleware to all routes under basePath, except for GET requests, and /v1/user/create
	basepath.Use(func(c *gin.Context) {
		if c.Request.Method != "GET" && c.FullPath() != "/v1/user/create" {
			authMiddleware.MiddlewareFunc()(c)
		}
	})
	userController.RegisterUserRoute(basepath)
	trackController.RegisterTrackRouter(basepath)
	albumController.RegisterAlbumRouter(basepath)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
