package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()

	router.GET("/api", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "The API is running and healthy.",
		})
	})

	// Register all the routes and there corresponding handlers
	router.GET("/api/orgs", getAllOrgs)
	router.GET("/api/orgs/:id", getOrg)
	router.GET("/api/projects", getProjectsByFilter)
	router.GET("/api/orgs/:id/projects/:projectId", getProject)
	router.GET("/api/orgs/:id/projects", getProjectsByYear)

	router.Run("0.0.0.0:8080")
}
