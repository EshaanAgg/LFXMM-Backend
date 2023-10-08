package api

import (
	"eshaanagg/lfx/database/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getProjectsByFilter(c *gin.Context) {
	filterText := c.Query("filterText")
	client := handlers.New()
	defer client.Close()

	c.IndentedJSON(http.StatusOK, client.GetProjectsByFilter(filterText))
}
