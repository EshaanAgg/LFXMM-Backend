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
	router.GET("/api/projects", getAllOrgs)
	router.GET("/api/projects/:id", getOrg)

	// creating routes for orgs
	router.GET("/api/orgs" , getAllOrgs)
	router.GET("/api/orgs/:id", getOrg)
	

	router.Run("0.0.0.0:8080")
}
