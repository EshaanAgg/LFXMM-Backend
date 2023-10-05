package api

import (
	"eshaanagg/lfx/database/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAllProjects(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	c.IndentedJSON(http.StatusOK, client.GetAllProjects())
}
