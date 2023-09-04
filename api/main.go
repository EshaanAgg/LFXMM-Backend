package api

import "github.com/gin-gonic/gin"

func Start() {
	router := gin.Default()

	// Register all the routes and there corresponding handlers
	router.GET("/api/projects", getAllOrgs)
	router.GET("/api/projects/:id", getOrg)

	router.Run("localhost:8080")
}
