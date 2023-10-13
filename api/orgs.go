package api

import (
	"eshaanagg/lfx/database/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getOrg(c *gin.Context) {
	client := handlers.New()
	defer client.Close()

	id := c.Param("id")
	org := client.GetOrganizationByID(id)

	if org == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There is no organization with this ID.",
		},
		)
		return
	}

	projects := client.GetProjectsByParentOrgID(id)

	c.IndentedJSON(http.StatusOK, gin.H{
		"id":           org.ID,
		"logo":         org.Logo,
		"org":          org.Name,
		"projectCount": len(projects),
	})
}

// The following function sends an API response with the description of a given organization.(id)
func getOrgAbout(c *gin.Context){
	client := handlers.New()
	defer client.Close()

	// Get id from HTTP request
	id := c.Param("id")

	// Check if an org exists with given id
	org := client.GetOrganizationByID(id)

	if org == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "There is no organization with this ID.",
		},
		)
		return
	}

	// Use handler to get data from the database
	about := client.GetOrganizationDescription(id) // This function is defined at ../database/handlers/parentOrg.go

	// Send out HTTP response
	c.IndentedJSON(http.StatusOK, gin.H{
		"about": about,
	})
}