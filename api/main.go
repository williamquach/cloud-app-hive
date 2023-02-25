package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/williamquach/cloud-app-hive/docs/api"
	"net/http"
)

// @title Swagger - CloudAppHive API
// @version 1.0
// @description This is a sample CloudAppHive server.
// @termsOfService http://swagger.io/terms/

// @contact.name CP0 Support
// @contact.url http://www.cp0.io/support
// @contact.email support.cp0@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	router := gin.Default()
	router.GET("/", HealthCheck)
	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/applications", getApplications)
	router.POST("/applications", createApplication)

	const port = 8080
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error("Server failed to start: ", err)
		return
	} else {
		log.Info("Server started on port: ", port)
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}

// CodeSource - User Code Source, e.g. from a GitHub repository, or from a zip file.
type CodeSource int64

const (
	GITHUB CodeSource = iota
	ZIP
)

type CodeSourceInfo interface {
	CodeSource() CodeSource
}

type GithubSourceInfo struct {
	Repo   string `json:"repo"`   // e.g. "
	Branch string `json:"branch"` // e.g. "main"
}

func (g GithubSourceInfo) CodeSource() CodeSource {
	return GITHUB
}

type ZipSourceInfo struct {
	ZipFile string `json:"zip_file"` // e.g. "https://example.com/my-app.zip"
}

func (z ZipSourceInfo) CodeSource() CodeSource {
	return ZIP
}

// ApplicationConfig - User App Configuration for his Cloud Application
type ApplicationConfig struct {
	// e.g. "my-app"
	Name string `json:"name" validate:"required"`

	// e.g. "My awesome app"
	Description string `json:"description" validate:"required"`

	// e.g. "example.com"
	Domain string `json:"domain" validate:"required"`

	// e.g. 8080, 80 - default 80
	Port int `json:"port" validate:"required"`

	// e.g. "NodeJS", "Go", "Python", etc.
	Platform string `json:"platform" validate:"required"`

	// e.g. NodeJS 16.x, Go 1.x, Python 3.x, etc.
	Version string `json:"version" validate:"required"`

	// e.g. GithubSourceInfo, ZipSourceInfo
	Source CodeSourceInfo `json:"source" validate:"required" example:"GithubSourceInfo{Repo: \"my-app\", Branch: \"main\"}"`
}

// albums slice to seed record album data.
var applications = []ApplicationConfig{
	{
		Name:        "my-back-end-app",
		Description: "My nodejs back-end app",
		Domain:      "example.com",
		Port:        8080,
		Platform:    "NodeJS",
		Version:     "16.x",
		Source: GithubSourceInfo{
			Repo:   "github.com/username/my-back-end-app",
			Branch: "main",
		},
	},
	{
		Name:        "my-front-end-app",
		Description: "My react front-end app",
		Domain:      "example.com",
		Port:        80,
		Platform:    "React",
		Version:     "17.x",
		Source: ZipSourceInfo{
			ZipFile: "https://example.com/my-front-end-app.zip",
		},
	},
}

// getApplications responds with the list of all applications as JSON.
func getApplications(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, applications)
}

// createApplication adds an application from JSON received in the request body.
func createApplication(context *gin.Context) {
	var newApplication ApplicationConfig
	log.Println("newApplication: ", newApplication)

	// Call BindJSON to bind the received JSON to newApplication.
	if err := context.BindJSON(&newApplication); err != nil {
		// If an error occurred, return an HTTP error.
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"name":    "invalid_body",
			"message": "Cannot parse body as JSON",
			"error":   err.Error(),
		})
		return
	}

	// Add the new album to the slice.
	applications = append(applications, newApplication)
	context.IndentedJSON(http.StatusCreated, newApplication)
}
