package api

import (
	"fmt"
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
	router.GET("/api/projects/:projectId", getProject1)
	router.GET("/api/orgs/:id/projects/:projectId", getProject)
	router.GET("/api/orgs/:id/count", getProjectCount)
	router.GET("/api/orgs/:id/projects", getProjectsByYear)
	router.GET("/api/allSkills", getAllSkills)
	router.GET("/api/projectdesc", getProjectDesc)
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Router was not able to run and map the requests")
		fmt.Println(err)
	}
}
