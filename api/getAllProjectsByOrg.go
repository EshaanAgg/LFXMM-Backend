package api

import (
	"eshaanagg/lfx/database/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAllProjectsByOrg(c *gin.Context) {

	client := handlers.New()
	defer client.Close()

	orgID := c.Param("orgID")
	projects := client.GetProjectsByParentOrgID(orgID)

	if projects == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There are no projects with this organisation",
		},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"orgID":    orgID,
		"projects": projects,
	})
}
