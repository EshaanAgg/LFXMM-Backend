package api

import (
	"eshaanagg/lfx/database"
	"eshaanagg/lfx/database/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getProjectsByFilter(c *gin.Context) {

	filterText := c.DefaultQuery("filterText", "")
	client := handlers.New()
	defer client.Close()

	var projects []database.ProjectThumbnail
	var err error

	if filterText != "" {
		projects, err = client.GetProjectsByFilter(filterText)
	} else {
		projects = client.GetProjectsByOrganization("1")
	}

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch projects",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, projects)
}
